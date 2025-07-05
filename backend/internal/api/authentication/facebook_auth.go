package authentication

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
)

var (
	Facebook_ClientID     = os.Getenv("FACEBOOK_CLIENT_ID")
	Facebook_ClientSecret = os.Getenv("FACEBOOK_CLIENT_SECRET")
	Facebook_RedirectURI  = "http://localhost:9000/auth/facebook/callback"
	Facebook_AuthURL      = "https://www.facebook.com/v18.0/dialog/oauth"
	Facebook_TokenURL     = "https://graph.facebook.com/v18.0/oauth/access_token"
	Facebook_UserInfoURL  = "https://graph.facebook.com/me?fields=id,name,email"
)

func RedirectToFacebookLogin(w http.ResponseWriter, r *http.Request) {
	authURL := fmt.Sprintf("%s?client_id=%s&redirect_uri=%s&scope=email&response_type=code",
		Facebook_AuthURL, Facebook_ClientID, Facebook_RedirectURI)
	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}

func HandleFacebookCallback(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var serverresponse Response
	statusCode := http.StatusOK

	code := r.URL.Query().Get("code")
	if code == "" {
		serverresponse.Message = "Bad request: No code received"
		statusCode = http.StatusBadRequest
		respondJSON(w, statusCode, serverresponse)
		return
	}

	// Step 1: Exchange code for access token
	tokenResp, err := http.Get(fmt.Sprintf("%s?client_id=%s&redirect_uri=%s&client_secret=%s&code=%s",
		Facebook_TokenURL, Facebook_ClientID, Facebook_RedirectURI, Facebook_ClientSecret, code))
	if err != nil {
		serverresponse.Message = "Failed to load token"
		statusCode = http.StatusInternalServerError
		respondJSON(w, statusCode, serverresponse)
		return
	}
	defer tokenResp.Body.Close()

	body, err := io.ReadAll(tokenResp.Body)
	if err != nil {
		serverresponse.Message = "Failed to read token response"
		statusCode = http.StatusInternalServerError
		respondJSON(w, statusCode, serverresponse)
		return
	}

	var tokenData map[string]interface{}
	if err := json.Unmarshal(body, &tokenData); err != nil {
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

	// Step 2: Get user info
	userInfoURL := fmt.Sprintf("%s&access_token=%s", Facebook_UserInfoURL, accessToken)
	userResp, err := http.Get(userInfoURL)
	if err != nil {
		serverresponse.Message = "Failed to fetch user info"
		statusCode = http.StatusInternalServerError
		respondJSON(w, statusCode, serverresponse)
		return
	}
	defer userResp.Body.Close()

	userData, err := io.ReadAll(userResp.Body)
	if err != nil {
		serverresponse.Message = "Failed to read user body"
		statusCode = http.StatusInternalServerError
		respondJSON(w, statusCode, serverresponse)
		return
	}

	var facebookUser GoogleUserInfo // reusing struct is fine
	if err := json.Unmarshal(userData, &facebookUser); err != nil {
		serverresponse.Message = "Failed to parse user info"
		statusCode = http.StatusInternalServerError
		respondJSON(w, statusCode, serverresponse)
		return
	}

	// Remaining logic (saving user, generating session) remains unchanged...
	userID, err := SaveGoogleUser(facebookUser, db)
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
