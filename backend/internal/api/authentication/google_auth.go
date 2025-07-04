package authentication

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	ClientID     = os.Getenv("CLIENT-ID")
	ClientSecret = os.Getenv("CLIENT-SECRET") // client secret
	RedirectURI  = "http://localhost:9000/auth/google/callback"
	AuthURL      = "https://accounts.google.com/o/oauth2/auth"
	TokenURL     = "https://oauth2.googleapis.com/token"
	UserInfoURL  = "https://www.googleapis.com/oauth2/v2/userinfo"
)

type GoogleUserInfo struct {
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

func RedirectToGoogleLogin(w http.ResponseWriter, r *http.Request) {
	authURL := fmt.Sprintf("%s?client_id=%s&redirect_uri=%s&response_type=code&scope=email profile",
		AuthURL, ClientID, RedirectURI)
	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}

func HandleGoogleCallback(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var serverresponse Response
	statusCode := http.StatusOK

	// Step 1: Get the auth code from Google
	code := r.URL.Query().Get("code")
	if code == "" {
		serverresponse.Message = "Bad request: No code received"
		statusCode = http.StatusBadRequest
		respondJSON(w, statusCode, serverresponse)
		return
	}

	// Step 2: Exchange code for access token
	data := url.Values{}
	data.Set("client_id", ClientID)
	data.Set("client_secret", ClientSecret)
	data.Set("code", code)
	data.Set("redirect_uri", RedirectURI)
	data.Set("grant_type", "authorization_code")

	tokenResp, err := http.Post(TokenURL, "application/x-www-form-urlencoded", bytes.NewBufferString(data.Encode()))
	if err != nil {
		serverresponse.Message = "Failed to load token"
		statusCode = http.StatusInternalServerError
		respondJSON(w, statusCode, serverresponse)
		return
	}
	defer tokenResp.Body.Close()

	tokenBody, err := io.ReadAll(tokenResp.Body)
	if err != nil {
		serverresponse.Message = "Failed to read token response"
		statusCode = http.StatusInternalServerError
		respondJSON(w, statusCode, serverresponse)
		return
	}

	var tokenData map[string]interface{}
	if err := json.Unmarshal(tokenBody, &tokenData); err != nil {
		serverresponse.Message = "Failed to parse token"
		statusCode = http.StatusInternalServerError
		respondJSON(w, statusCode, serverresponse)
		return
	}

	accessToken, ok := tokenData["access_token"].(string)
	if !ok || accessToken == "" {
		serverresponse.Message = "Access token missing"
		statusCode = http.StatusInternalServerError
		respondJSON(w, statusCode, serverresponse)
		return
	}

	// Step 3: Get user info from Google
	req, _ := http.NewRequest("GET", UserInfoURL, nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	userResp, err := client.Do(req)
	if err != nil {
		serverresponse.Message = "Failed to fetch user info"
		statusCode = http.StatusInternalServerError
		respondJSON(w, statusCode, serverresponse)
		return
	}
	defer userResp.Body.Close()

	userBody, err := io.ReadAll(userResp.Body)
	if err != nil {
		serverresponse.Message = "Failed to read user body"
		statusCode = http.StatusInternalServerError
		respondJSON(w, statusCode, serverresponse)
		return
	}

	var googleUser GoogleUserInfo
	if err := json.Unmarshal(userBody, &googleUser); err != nil {
		serverresponse.Message = "Failed to parse user info"
		statusCode = http.StatusInternalServerError
		respondJSON(w, statusCode, serverresponse)
		return
	}

	// Step 4: Save or find user in DB
	userID, err := SaveGoogleUser(googleUser, db)
	if err != nil {
		serverresponse.Message = "Failed to save user"
		statusCode = http.StatusInternalServerError
		respondJSON(w, statusCode, serverresponse)
		return
	}

	// Step 5: Manage session
	DeleteUserSessions(userID, db)

	sessionID := uuid.New().String()
	expiray := time.Now().Add(24 * time.Hour)
	if err := StoreSession(userID, sessionID, expiray, db); err != nil {
		serverresponse.Message = "Failed to create session"
		statusCode = http.StatusInternalServerError
		respondJSON(w, statusCode, serverresponse)
		return
	}
	// Step 6: Set session cookie
	cookie := http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		Expires:  expiray,
		HttpOnly: true,
		Secure:   false, // Change to true in production
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, &cookie)

	serverresponse.Message = "Login successful"
	statusCode = http.StatusOK
	respondJSON(w, statusCode, serverresponse)
}

func SaveGoogleUser(userInfo GoogleUserInfo, db *sql.DB) (int, error) {
	var userID int
	var username string

	// Check if the user already exists in the database
	row := db.QueryRow("SELECT id, nickname FROM Users WHERE email = ?", userInfo.Email)
	err := row.Scan(&userID, &username)
	if err != nil {
		// If no existing user is found, create a new account
		if err == sql.ErrNoRows {
			// Generate a unique username based on the user's Google name
			newUsername, err := CreateUniqueUsername(userInfo.Name, db)
			if err != nil {
				return 0, err
			}

			user := Profile_User{
				Email:           userInfo.Email,
				Password:        "",
				FirstName:       userInfo.Name,
				LastName:        "",
				DateOfBirth:     "",
				Nickname:        newUsername,
				AboutMe:         "",
				IsProfilePublic: false,
				Avatar:          userInfo.Picture,
			}
			user.Avatar, err = DownloadAndSavePicture(userInfo.Picture)
			if err != nil {
				return -1, err
			}

			// Insert the new user into the database (password is empty for Google login)
			if err := InsertUserIntoDB(user, db); err != nil {
				return -1, err
			}

			// Retrieve the newly inserted user ID
			err = db.QueryRow("SELECT id FROM Users WHERE email = ?", userInfo.Email).Scan(&userID)
			if err != nil {
				fmt.Println("Error fetching user ID:", err)
				return -1, err
			}

			return userID, nil
		} else {
			// Return other database errors
			return -1, err
		}
	}

	// If the user already exists, return their user ID
	return userID, nil
}

func CreateUniqueUsername(name string, db *sql.DB) (string, error) {
	// Remove spaces from the name to create a base username
	username := strings.ReplaceAll(name, " ", "")
	var count int
	suffix := 1

	for {
		// Check if the username already exists in the database
		err := db.QueryRow("SELECT COUNT(*) FROM Users WHERE nickname = ?", username).Scan(&count)
		if err != nil {
			return "", err // Return an error if the query fails
		}

		// If no existing username matches, return the generated username
		if count == 0 {
			break
		} else {
			// Append a number to the username to ensure uniqueness
			suffix++
			username += strconv.Itoa(suffix)
		}
	}

	return username, nil // Return the unique username
}
