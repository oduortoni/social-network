// login with twitter
package authentication

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const (
	Twitter_ClientID    = "YOUR_TWITTER_CLIENT_ID"
	Twitter_RedirectURI = "http://localhost:8080/auth/twitter/callback"
	Twitter_AuthURL     = "https://twitter.com/i/oauth2/authorize"
	Twitter_TokenURL    = "https://api.twitter.com/2/oauth2/token"
	Twitter_UserInfoURL = "https://api.twitter.com/2/users/me"
	Twitter_Scopes      = "users.read offline.access" // Removed tweet.read as it's not needed for auth
)

// generateCodeVerifier creates a cryptographically secure random code verifier
func generateCodeVerifier() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b)
}

// generateCodeChallenge creates the code challenge from verifier using SHA256
func generateCodeChallenge(verifier string) string {
	h := sha256.New()
	h.Write([]byte(verifier))
	return base64.RawURLEncoding.EncodeToString(h.Sum(nil))
}

func RedirectToTwitterLogin(w http.ResponseWriter, r *http.Request) {
	// Generate secure code verifier and challenge
	codeVerifier := generateCodeVerifier()
	codeChallenge := generateCodeChallenge(codeVerifier)

	// Save codeVerifier in session/cookie/store
	http.SetCookie(w, &http.Cookie{
		Name:     "code_verifier",
		Value:    codeVerifier,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
		SameSite: http.SameSiteLaxMode,
	})

	// Build auth URL with proper variable references
	authURL := fmt.Sprintf(
		"%s?response_type=code&client_id=%s&redirect_uri=%s&scope=%s&state=xyz&code_challenge=%s&code_challenge_method=S256",
		Twitter_AuthURL, Twitter_ClientID, url.QueryEscape(Twitter_RedirectURI), url.QueryEscape(Twitter_Scopes), codeChallenge,
	)

	http.Redirect(w, r, authURL, http.StatusFound)
}

func HandleTwitterCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Missing code", http.StatusBadRequest)
		return
	}

	// Get code_verifier from cookie/session
	cookie, err := r.Cookie("code_verifier")
	if err != nil {
		http.Error(w, "Missing code verifier", http.StatusBadRequest)
		return
	}
	codeVerifier := cookie.Value

	// Prepare token request
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("redirect_uri", Twitter_RedirectURI)
	data.Set("client_id", Twitter_ClientID)
	data.Set("code_verifier", codeVerifier)

	// Make token request
	req, err := http.NewRequest("POST", Twitter_TokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		http.Error(w, "Token request error", http.StatusInternalServerError)
		return
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Token response error", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Handle token response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Error reading token response", http.StatusInternalServerError)
		return
	}

	if resp.StatusCode != http.StatusOK {
		http.Error(w, fmt.Sprintf("Token request failed: %s", string(body)), http.StatusInternalServerError)
		return
	}

	var tokenRes map[string]interface{}
	if err := json.Unmarshal(body, &tokenRes); err != nil {
		http.Error(w, "Error parsing token response", http.StatusInternalServerError)
		return
	}

	accessToken, ok := tokenRes["access_token"].(string)
	if !ok {
		http.Error(w, "No access token in response", http.StatusInternalServerError)
		return
	}

	// Get user info
	userReq, err := http.NewRequest("GET", Twitter_UserInfoURL, nil)
	if err != nil {
		http.Error(w, "User info request error", http.StatusInternalServerError)
		return
	}
	userReq.Header.Set("Authorization", "Bearer "+accessToken)

	userResp, err := client.Do(userReq)
	if err != nil {
		http.Error(w, "User info response error", http.StatusInternalServerError)
		return
	}
	defer userResp.Body.Close()

	userBody, err := io.ReadAll(userResp.Body)
	if err != nil {
		http.Error(w, "Error reading user info response", http.StatusInternalServerError)
		return
	}

	if userResp.StatusCode != http.StatusOK {
		http.Error(w, fmt.Sprintf("User info request failed: %s", string(userBody)), http.StatusInternalServerError)
		return
	}

	var twitterUser struct {
		Data struct {
			ID       string `json:"id"`
			Name     string `json:"name"`
			Username string `json:"username"`
		} `json:"data"`
	}

	if err := json.Unmarshal(userBody, &twitterUser); err != nil {
		http.Error(w, "Error parsing user info response", http.StatusInternalServerError)
		return
	}

	// Clear the code verifier cookie
	http.SetCookie(w, &http.Cookie{
		Name:   "code_verifier",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})

	fmt.Printf("Logged in as: @%s (%s)", twitterUser.Data.Username, twitterUser.Data.Name)
}
