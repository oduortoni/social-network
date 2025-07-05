package authentication

import (
	"bytes"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"

	"github.com/google/uuid"
)

type stateToken struct {
	token      string
	createTime time.Time
}

var (
	// Store state tokens with mutex for thread safety
	stateTokens = make(map[string]stateToken)
	stateMutex  sync.RWMutex

	// Token expiration duration
	stateTokenExpiration = 10 * time.Minute
)

// GitHubConfig holds the OAuth configuration
type GitHubConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURI  string
	AuthURL      string
	TokenURL     string
	UserInfoURL  string
	EmailURL     string
}

// GitHubUserInfo represents the user information returned by GitHub
type GitHubUserInfo struct {
	Login      string `json:"login"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Avatar_url string `json:"avatar_url"`
	ID         int    `json:"id"`
	NodeID     string `json:"node_id"`
}

// Initialize GitHub OAuth configuration
var githubConfig = GitHubConfig{
	ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
	ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
	RedirectURI:  "http://localhost:9000/auth/github/callback",
	AuthURL:      "https://github.com/login/oauth/authorize",
	TokenURL:     "https://github.com/login/oauth/access_token",
	UserInfoURL:  "https://api.github.com/user",
	EmailURL:     "https://api.github.com/user/emails", // Fetch user emails
}

type GitHubEmail struct {
	Email    string `json:"email"`
	Primary  bool   `json:"primary"`
	Verified bool   `json:"verified"`
}

func RedirectToGitHubLogin(w http.ResponseWriter, r *http.Request) {
	authURL := fmt.Sprintf("%s?client_id=%s&redirect_uri=%s&scope=user:email&state=%s",
		githubConfig.AuthURL,
		githubConfig.ClientID,
		githubConfig.RedirectURI,
		generateStateToken()) // implement this is important coz it to prevent CSRF attacks

	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}

func HandleGitHubCallback(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var serverresponse Response
	statusCode := http.StatusOK

	// Step 1: Validate state token
	state := r.URL.Query().Get("state")
	if !validateStateToken(state) {
		serverresponse.Message = "Invalid state parameter"
		statusCode = http.StatusBadRequest
		respondJSON(w, statusCode, serverresponse)
		return
	}

	// Step 2: Get authorization code
	code := r.URL.Query().Get("code")
	if code == "" {
		serverresponse.Message = "Authorization code not provided"
		statusCode = http.StatusBadRequest
		respondJSON(w, statusCode, serverresponse)
		return
	}

	// Step 3: Exchange code for access token
	accessToken, err := getGitHubAccessToken(code)
	if err != nil || accessToken == "" {
		serverresponse.Message = "Failed to exchange code for access token"
		statusCode = http.StatusInternalServerError
		respondJSON(w, statusCode, serverresponse)
		return
	}

	// Step 4: Get GitHub user profile
	githubUser, err := getGitHubUserInfo(accessToken)
	if err != nil {
		serverresponse.Message = "Failed to fetch user profile from GitHub"
		statusCode = http.StatusInternalServerError
		respondJSON(w, statusCode, serverresponse)
		return
	}


	// Step 5: Get GitHub user email
	githubUserEmail, err := getGitHubEmail(accessToken)
	if err != nil {
		serverresponse.Message = "Failed to fetch user email from GitHub"
		statusCode = http.StatusInternalServerError
		respondJSON(w, statusCode, serverresponse)
		return
	}

	// Step 6: Prepare user data
	githubUser.Email = GetOrGenerateEmail(githubUserEmail, githubUser.Name)

	user := GoogleUserInfo{
		Name:    githubUser.Name,
		Email:   githubUser.Email,
		Picture: githubUser.Avatar_url,
	}

	// Step 7: Save or find user in DB
	userID, err := SaveGoogleUser(user, db)
	if err != nil {
		serverresponse.Message = "Failed to save user"
		statusCode = http.StatusInternalServerError
		respondJSON(w, statusCode, serverresponse)
		return
	}

	// Step 8: Clear old sessions
	_ = DeleteUserSessions(userID, db)

	// Step 9: Generate new session
	sessionID := uuid.New().String()
	expiry := time.Now().Add(24 * time.Hour)
	if err := StoreSession(userID, sessionID, expiry, db); err != nil {
		serverresponse.Message = "Failed to store session"
		statusCode = http.StatusInternalServerError
		respondJSON(w, statusCode, serverresponse)
		return
	}

	// Step 10: Set session cookie
	cookie := http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		Expires:  expiry,
		HttpOnly: true,
		Secure:   false, // Set to true in production
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, &cookie)

	// Step 11: Respond success
	serverresponse.Message = "Login successful"
	statusCode = http.StatusOK
	respondJSON(w, statusCode, serverresponse)
}

func getGitHubAccessToken(code string) (string, error) {
	data := url.Values{}
	data.Set("client_id", githubConfig.ClientID)
	data.Set("client_secret", githubConfig.ClientSecret)
	data.Set("code", code)
	data.Set("redirect_uri", githubConfig.RedirectURI)

	req, err := http.NewRequest("POST", githubConfig.TokenURL, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return "", err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var tokenResponse struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		Scope       string `json:"scope"`
	}

	if err := json.Unmarshal(body, &tokenResponse); err != nil {
		return "", err
	}

	return tokenResponse.AccessToken, nil
}

func getGitHubUserInfo(accessToken string) (*GitHubUserInfo, error) {
	req, err := http.NewRequest("GET", githubConfig.UserInfoURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "token "+accessToken)
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var userInfo GitHubUserInfo
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, err
	}

	return &userInfo, nil
}

func getGitHubEmail(accessToken string) ([]GitHubEmail, error) {
	req, err := http.NewRequest("GET", githubConfig.EmailURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "token "+accessToken)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var emails []GitHubEmail
	if err := json.Unmarshal(body, &emails); err != nil {
		return nil, err
	}

	return emails, nil
}

func generateStateToken() string {
	// Generate 32 bytes of random data
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		panic(err) // In production, handle this error appropriately
	}

	// Encode to base64URL (URL-safe version of base64)
	token := base64.URLEncoding.EncodeToString(bytes)

	// Store the token with creation time
	stateMutex.Lock()
	stateTokens[token] = stateToken{
		token:      token,
		createTime: time.Now(),
	}

	// Clean up expired tokens
	for k, v := range stateTokens {
		if time.Since(v.createTime) > stateTokenExpiration {
			delete(stateTokens, k)
		}
	}
	stateMutex.Unlock()

	return token
}

func validateStateToken(state string) bool {
	stateMutex.RLock()
	defer stateMutex.RUnlock()

	// Check if token exists and is not expired
	if token, exists := stateTokens[state]; exists {
		if time.Since(token.createTime) <= stateTokenExpiration {
			// Clean up used token
			go func() {
				stateMutex.Lock()
				delete(stateTokens, state)
				stateMutex.Unlock()
			}()
			return true
		}
	}
	return false
}

// GetOrGenerateEmail checks the list of GitHub emails and returns the primary verified email if available.
// If no primary verified email is found, it returns any verified email.
// If no verified email exists, it generates a fake email using the provided username.
func GetOrGenerateEmail(emails []GitHubEmail, username string) string {
	// First, check for a primary verified email.
	for i := 0; i < len(emails); i++ {
		if emails[i].Primary && emails[i].Verified {
			return emails[i].Email
		}
	}

	// If no primary verified email exists, return any verified email.
	for i := 0; i < len(emails); i++ {
		if emails[i].Verified {
			return emails[i].Email
		}
	}

	// If no verified email exists, generate a fake email.
	return username + "@fakeEmail.com"
}
