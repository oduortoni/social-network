package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"html"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/tajjjjr/social-network/backend/internal/models"
	"github.com/tajjjjr/social-network/backend/internal/service"
	"github.com/tajjjjr/social-network/backend/utils"
)

// AuthHandler handles HTTP requests for authentication.
type AuthHandler struct {
	AuthService service.AuthServiceInterface
}

// NewAuthHandler creates a new AuthHandler.
func NewAuthHandler(as service.AuthServiceInterface) *AuthHandler {
	return &AuthHandler{AuthService: as}
}

// LoginRequest represents the request body for a login request.
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// SignupRequest represents the request body for a signup request.
type SignupRequest struct {
	Email           string `json:"email"`
	Password        string `json:"password"`
	FirstName       string `json:"firstName"`
	LastName        string `json:"lastName"`
	DateOfBirth     string `json:"dob"`
	Nickname        string `json:"nickname"`
	AboutMe         string `json:"aboutMe"`
	IsProfilePublic bool   `json:"profileVisibility"`
	Avatar          string `json:"avatar"`
}

// StepOneCredintial represents the user's credentials on step One of registration.
type StepOneCredintial struct {
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

func (auth *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var creds LoginRequest
	var err error

	// Check Content-Type header to determine how to parse the request
	contentType := r.Header.Get("Content-Type")
	fmt.Println("Content-Type:", contentType)

	if strings.Contains(contentType, "application/json") {
		// Parse JSON body
		err = json.NewDecoder(r.Body).Decode(&creds)
		if err != nil {
			models.RespondJSON(w, http.StatusBadRequest, models.Response{Message: "Invalid JSON request body"})
			return
		}
	} else if strings.Contains(contentType, "application/x-www-form-urlencoded") || strings.Contains(contentType, "multipart/form-data") {
		// Parse form data
		err = r.ParseForm()
		if err != nil {
			models.RespondJSON(w, http.StatusBadRequest, models.Response{Message: "Invalid form data"})
			return
		}

		// Extract form values
		creds.Email = r.FormValue("email")
		creds.Password = r.FormValue("password")
	} else {
		// Default to trying JSON first, then form data
		err = json.NewDecoder(r.Body).Decode(&creds)
		if err != nil {
			// If JSON fails, try to parse as form data
			if parseErr := r.ParseForm(); parseErr != nil {
				models.RespondJSON(w, http.StatusBadRequest, models.Response{Message: "Invalid request body format"})
				return
			}
			creds.Email = r.FormValue("email")
			creds.Password = r.FormValue("password")
		}
	}

	fmt.Println("Login credentials:", creds)

	authUser, sessionID, _ := auth.AuthService.AuthenticateUser(creds.Email, creds.Password)
	if authUser == nil {
		if sessionID == service.EXPIRED_SESSION {
			models.RespondJSON(w, http.StatusInternalServerError, models.Response{Message: "Failed to create session"})
		} else if sessionID == service.INVALID_PASSWORD {
			models.RespondJSON(w, http.StatusUnauthorized, models.Response{Message: "Invalid password"})
		} else if sessionID == service.INVALID_EMAIL {
			models.RespondJSON(w, http.StatusUnauthorized, models.Response{Message: "User not found"})
		}
		return
	}

	fmt.Println("Authenticated user:", authUser)

	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // set true in production with HTTPS
		Expires:  time.Now().Add(7 * 24 * time.Hour),
		SameSite: http.SameSiteLaxMode,
	}

	fmt.Println("Setting cookie:", cookie)
	http.SetCookie(w, cookie)
	models.RespondJSON(w, http.StatusOK, models.Response{Message: "Logged in successfully"})
}

// Signup handles user registration
func (auth *AuthHandler) Signup(w http.ResponseWriter, r *http.Request) {
	// Parse multipart form (limit: 10MB)
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		models.RespondJSON(w, http.StatusBadRequest, models.Response{Message: "Failed to parse form"})
		return
	}

	// Extract form values
	email := r.FormValue("email")
	password := r.FormValue("password")
	firstName := r.FormValue("firstName")
	lastName := r.FormValue("lastName")
	dateOfBirth := r.FormValue("dob")
	nickname := r.FormValue("nickname")
	aboutMe := r.FormValue("aboutMe")
	isProfilePublic := r.FormValue("profileVisibility")

	// Sanitize user input to prevent XSS attacks
	email = html.EscapeString(email)
	firstName = html.EscapeString(firstName)
	lastName = html.EscapeString(lastName)
	dateOfBirth = html.EscapeString(dateOfBirth)
	nickname = html.EscapeString(nickname)
	aboutMe = html.EscapeString(aboutMe)

	// Validate email format
	IsEmailValid, err := auth.AuthService.ValidateEmail(email)
	if err != nil {
		models.RespondJSON(w, http.StatusInternalServerError, models.Response{Message: "Regex error in validating Email"})
		return
	}
	if !IsEmailValid {
		models.RespondJSON(w, http.StatusBadRequest, models.Response{Message: "Invalid email format"})
		return
	}

	// Handle avatar upload
	userAvatar := "no profile photo"
	file, header, err := r.FormFile("avatar")
	if err == nil && file != nil {
		defer file.Close()
		userAvatar, err = UploadAvatarImage(file, header)
		if err != nil {
			models.RespondJSON(w, http.StatusInternalServerError, models.Response{Message: userAvatar})
			return
		}
	}

	// Create user model
	user := &models.User{
		Email:           email,
		Password:        password, // Will be hashed in service layer
		FirstName:       &firstName,
		LastName:        &lastName,
		DateOfBirth:     &dateOfBirth,
		Nickname:        &nickname,
		AboutMe:         &aboutMe,
		IsProfilePublic: isProfilePublic == "public",
		Avatar:          &userAvatar,
	}

	// Create user through service
	createdUser, err := auth.AuthService.CreateUser(user)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			models.RespondJSON(w, http.StatusConflict, models.Response{Message: "Email or nickname already taken"})
		} else {
			models.RespondJSON(w, http.StatusInternalServerError, models.Response{Message: "Failed to create user"})
		}
		return
	}

	fmt.Println("User created successfully:", createdUser.ID)
	models.RespondJSON(w, http.StatusOK, models.Response{Message: "Registration successful"})
}

// LogoutHandler deletes session and clears cookie
func (auth *AuthHandler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err == nil {
		_, _ = auth.AuthService.DeleteSession(cookie.Value)

		// Clear cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "session_id",
			Value:    "",
			Path:     "/",
			MaxAge:   -1,
			HttpOnly: true,
			Secure:   false, // set true in production
			SameSite: http.SameSiteLaxMode,
		})
	}

	models.RespondJSON(w, http.StatusOK, models.Response{Message: "Logged out successfully"})
}

// AuthMiddleware verifies session cookie, loads user ID into context
func (auth *AuthHandler) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_id")
		if err != nil {
			models.RespondJSON(w, http.StatusUnauthorized, models.Response{Message: "Authentication required"})
			return
		}

		userID, err := auth.AuthService.GetUserIDBySession(cookie.Value)
		if err != nil {
			models.RespondJSON(w, http.StatusUnauthorized, models.Response{Message: "Invalid or expired session"})
			return
		}

		// Add userID to request context for downstream handlers
		ctx := context.WithValue(r.Context(), utils.User_id, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// ValidateAccountStepOne validates Account Crediential for Step One
func (auth *AuthHandler) ValidateAccountStepOne(w http.ResponseWriter, r *http.Request) {
	var serverresponse models.Response
	statusCode := http.StatusOK
	var AccountCrediential StepOneCredintial

	body, err := io.ReadAll(r.Body)
	if err != nil {
		serverresponse.Message = "Failed to read request body"
		statusCode = http.StatusInternalServerError
		models.RespondJSON(w, statusCode, serverresponse)
		return
	}
	if err := json.Unmarshal(body, &AccountCrediential); err != nil {
		serverresponse.Message = "Failed to parse request body"
		statusCode = http.StatusBadRequest
		models.RespondJSON(w, statusCode, serverresponse)
		return
	}

	// Validate required fields
	if AccountCrediential.Email == "" || AccountCrediential.Password == "" || AccountCrediential.ConfirmPassword == "" {
		serverresponse.Message = "Missing required fields"
		statusCode = http.StatusBadRequest
		models.RespondJSON(w, statusCode, serverresponse)
		return
	}

	// validate email
	IsEmailValid, err := auth.AuthService.ValidateEmail(AccountCrediential.Email)
	if err != nil {
		serverresponse.Message = "Regex error in validating Email"
		statusCode = http.StatusInternalServerError
		models.RespondJSON(w, statusCode, serverresponse)
		return
	}
	if !IsEmailValid {
		serverresponse.Message = "Invalid Email format"
		statusCode = http.StatusBadRequest
		models.RespondJSON(w, statusCode, serverresponse)
		return
	}

	// password is same as confirm passwod
	if AccountCrediential.Password != AccountCrediential.ConfirmPassword {
		serverresponse.Message = "Passwords do not match."
		statusCode = http.StatusBadRequest
		models.RespondJSON(w, statusCode, serverresponse)
		return
	}

	// Check if user already exists
	if UserExists, err := auth.AuthService.UserExists(AccountCrediential.Email); err != nil || UserExists {
		serverresponse.Message = "Email already exists"
		statusCode = http.StatusConflict
		models.RespondJSON(w, statusCode, serverresponse)
		return
	}
	// validate password
	passwordManager := utils.NewPasswordManager(utils.PasswordConfig{})
	_, err = passwordManager.HashPassword(AccountCrediential.Password)
	if err != nil {
		serverresponse.Message = err.Error()
		statusCode = http.StatusBadRequest
		models.RespondJSON(w, statusCode, serverresponse)
		return
	}

	serverresponse.Message = "Ok"
	models.RespondJSON(w, statusCode, serverresponse)
}
