package tests

import (
	"bytes"
	"encoding/json"
	"errors"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/tajjjjr/social-network/backend/internal/api/handlers"
	"github.com/tajjjjr/social-network/backend/internal/models"
	"github.com/tajjjjr/social-network/backend/pkg/utils"
	"golang.org/x/crypto/bcrypt"
)

// createMultipartForm creates a multipart form request body for testing
func createMultipartForm(fields map[string]string) (*bytes.Buffer, string, error) {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	for key, value := range fields {
		if err := writer.WriteField(key, value); err != nil {
			return nil, "", err
		}
	}

	if err := writer.Close(); err != nil {
		return nil, "", err
	}

	return &body, writer.FormDataContentType(), nil
}

func TestSignup_Success(t *testing.T) {
	mockAuthService := &MockAuthService{
		CreateUserFunc: func(user *models.User) (*models.User, error) {
			user.ID = 1
			return user, nil
		},
		ValidateEmailFunc: func(email string) (bool, error) {
			return true, nil
		},
		UserExistsFunc: func(email string) (bool, error) {
			return false, nil
		},
		UserNewEditEmailExistFunc: func(email string, userid int64) (bool, error) {
			return false, nil
		},
		EditUserProfileFunc: func(user *models.User, userid int64) error {
			return nil
		},
	}
	authHandler := handlers.NewAuthHandler(mockAuthService)

	fields := map[string]string{
		"email":             "test@example.com",
		"password":          "Password1!",
		"firstName":         "John",
		"lastName":          "Doe",
		"dob":               "2000-01-01",
		"nickname":          "johndoe",
		"aboutMe":           "Hello",
		"profileVisibility": "public",
	}
	body, contentType, err := createMultipartForm(fields)
	if err != nil {
		t.Fatalf("Failed to create multipart form: %v", err)
	}

	req := httptest.NewRequest("POST", "/register", body)
	req.Header.Set("Content-Type", contentType)
	rr := httptest.NewRecorder()

	authHandler.Signup(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var resp utils.Response
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}

	if resp.Message != "Registration successful" {
		t.Errorf("expected 'Registration successful', got %q", resp.Message)
	}
}

func TestSignup_UserAlreadyExists(t *testing.T) {
	mockAuthService := &MockAuthService{
		CreateUserFunc: func(user *models.User) (*models.User, error) {
			return nil, errors.New("user with email test@example.com already exists")
		},
		ValidateEmailFunc: func(email string) (bool, error) {
			return true, nil
		},
		UserExistsFunc: func(email string) (bool, error) {
			return true, nil // User already exists
		},
	}
	authHandler := handlers.NewAuthHandler(mockAuthService)

	fields := map[string]string{
		"email":             "test@example.com",
		"password":          "Password1!",
		"firstName":         "John",
		"lastName":          "Doe",
		"dob":               "2000-01-01",
		"nickname":          "johndoe",
		"aboutMe":           "Hello",
		"profileVisibility": "public",
	}
	body, contentType, err := createMultipartForm(fields)
	if err != nil {
		t.Fatalf("Failed to create multipart form: %v", err)
	}

	req := httptest.NewRequest("POST", "/register", body)
	req.Header.Set("Content-Type", contentType)
	rr := httptest.NewRecorder()

	authHandler.Signup(rr, req)

	if status := rr.Code; status != http.StatusConflict {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusConflict)
	}

	var resp utils.Response
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}

	if resp.Message != "Email or nickname already taken" {
		t.Errorf("expected 'Email or nickname already taken', got %q", resp.Message)
	}
}

func TestSignup_InvalidEmail(t *testing.T) {
	mockAuthService := &MockAuthService{
		CreateUserFunc: func(user *models.User) (*models.User, error) {
			return user, nil
		},
		ValidateEmailFunc: func(email string) (bool, error) {
			// Return false for invalid emails to trigger the validation error
			return false, nil
		},
		UserExistsFunc: func(email string) (bool, error) {
			return false, nil
		},
	}
	authHandler := handlers.NewAuthHandler(mockAuthService)

	invalidEmails := []string{
		"invalid-email",
		"@example.com",
		"test@",
		"test.example.com",
		"test@.com",
		"test@example.",
		"",
		"   ",
	}

	for _, email := range invalidEmails {
		t.Run("InvalidEmail_"+email, func(t *testing.T) {
			fields := map[string]string{
				"email":             email,
				"password":          "Password1!",
				"firstName":         "John",
				"lastName":          "Doe",
				"dob":               "2000-01-01",
				"nickname":          "johndoe",
				"aboutMe":           "Hello",
				"profileVisibility": "public",
			}

			body, contentType, err := createMultipartForm(fields)
			if err != nil {
				t.Fatalf("Failed to create multipart form: %v", err)
			}

			req := httptest.NewRequest("POST", "/register", body)
			req.Header.Set("Content-Type", contentType)
			rr := httptest.NewRecorder()

			authHandler.Signup(rr, req)

			if status := rr.Code; status != http.StatusBadRequest {
				t.Errorf("expected status %d for invalid email '%s', got %d",
					http.StatusBadRequest, email, status)
			}

			var resp utils.Response
			if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
				t.Fatal(err)
			}

			if resp.Message != "Invalid email format" {
				t.Errorf("expected 'Invalid email format', got %q", resp.Message)
			}
		})
	}
}

func TestSignup_InvalidFormData(t *testing.T) {
	mockAuthService := &MockAuthService{
		CreateUserFunc: func(user *models.User) (*models.User, error) {
			return user, nil
		},
		ValidateEmailFunc: func(email string) (bool, error) {
			return true, nil
		},
		UserExistsFunc: func(email string) (bool, error) {
			return false, nil
		},
	}
	authHandler := handlers.NewAuthHandler(mockAuthService)

	form := "invalid-form-data"
	req := httptest.NewRequest("POST", "/register", bytes.NewBufferString(form))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()

	authHandler.Signup(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, status)
	}

	var resp utils.Response
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}

	if resp.Message != "Failed to parse form" {
		t.Errorf("expected 'Failed to parse form', got %q", resp.Message)
	}
}

func TestSignup_ServiceError(t *testing.T) {
	mockAuthService := &MockAuthService{
		CreateUserFunc: func(user *models.User) (*models.User, error) {
			return nil, errors.New("database connection failed")
		},
		ValidateEmailFunc: func(email string) (bool, error) {
			return true, nil
		},
		UserExistsFunc: func(email string) (bool, error) {
			return false, nil
		},
	}
	authHandler := handlers.NewAuthHandler(mockAuthService)

	fields := map[string]string{
		"email":             "test@example.com",
		"password":          "Password1!",
		"firstName":         "John",
		"lastName":          "Doe",
		"dob":               "2000-01-01",
		"nickname":          "johndoe",
		"aboutMe":           "Hello",
		"profileVisibility": "public",
	}
	body, contentType, err := createMultipartForm(fields)
	if err != nil {
		t.Fatalf("Failed to create multipart form: %v", err)
	}

	req := httptest.NewRequest("POST", "/register", body)
	req.Header.Set("Content-Type", contentType)
	rr := httptest.NewRecorder()

	authHandler.Signup(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
	}

	var resp utils.Response
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}

	if resp.Message != "Failed to create user" {
		t.Errorf("expected 'Failed to create user', got %q", resp.Message)
	}
}

func TestSignup_XSSPrevention(t *testing.T) {
	mockAuthService := &MockAuthService{
		CreateUserFunc: func(user *models.User) (*models.User, error) {
			// Verify that HTML is escaped
			if user.FirstName != nil && strings.Contains(*user.FirstName, "<script>") {
				t.Error("XSS content not properly escaped in firstName")
			}
			if user.LastName != nil && strings.Contains(*user.LastName, "<script>") {
				t.Error("XSS content not properly escaped in lastName")
			}
			if user.AboutMe != nil && strings.Contains(*user.AboutMe, "<script>") {
				t.Error("XSS content not properly escaped in aboutMe")
			}
			user.ID = 1
			return user, nil
		},
		ValidateEmailFunc: func(email string) (bool, error) {
			return true, nil
		},
		UserExistsFunc: func(email string) (bool, error) {
			return false, nil
		},
	}
	authHandler := handlers.NewAuthHandler(mockAuthService)

	fields := map[string]string{
		"email":             "test@example.com",
		"password":          "Password1!",
		"firstName":         "<script>alert('xss')</script>John",
		"lastName":          "<script>alert('xss')</script>Doe",
		"dob":               "2000-01-01",
		"nickname":          "johndoe",
		"aboutMe":           "<script>alert('xss')</script>Hello",
		"profileVisibility": "public",
	}
	body, contentType, err := createMultipartForm(fields)
	if err != nil {
		t.Fatalf("Failed to create multipart form: %v", err)
	}

	req := httptest.NewRequest("POST", "/register", body)
	req.Header.Set("Content-Type", contentType)
	rr := httptest.NewRecorder()

	authHandler.Signup(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestSignup_ProfileVisibility(t *testing.T) {
	testCases := []struct {
		name              string
		profileVisibility string
		expectedPublic    bool
	}{
		{"Public Profile", "public", true},
		{"Private Profile", "private", false},
		{"Empty Profile", "", false},
		{"Invalid Profile", "invalid", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockAuthService := &MockAuthService{
				CreateUserFunc: func(user *models.User) (*models.User, error) {
					if user.IsProfilePublic != tc.expectedPublic {
						t.Errorf("expected IsProfilePublic to be %v, got %v", tc.expectedPublic, user.IsProfilePublic)
					}
					user.ID = 1
					return user, nil
				},
				ValidateEmailFunc: func(email string) (bool, error) {
					return true, nil
				},
				UserExistsFunc: func(email string) (bool, error) {
					return false, nil
				},
			}
			authHandler := handlers.NewAuthHandler(mockAuthService)

			fields := map[string]string{
				"email":             "test@example.com",
				"password":          "Password1!",
				"firstName":         "John",
				"lastName":          "Doe",
				"dob":               "2000-01-01",
				"nickname":          "johndoe",
				"aboutMe":           "Hello",
				"profileVisibility": tc.profileVisibility,
			}
			body, contentType, err := createMultipartForm(fields)
			if err != nil {
				t.Fatalf("Failed to create multipart form: %v", err)
			}

			req := httptest.NewRequest("POST", "/register", body)
			req.Header.Set("Content-Type", contentType)
			rr := httptest.NewRecorder()

			authHandler.Signup(rr, req)

			if status := rr.Code; status != http.StatusOK {
				t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
			}
		})
	}
}

func TestSignup_MissingFields(t *testing.T) {
	testCases := []struct {
		name   string
		fields map[string]string
	}{
		{
			name: "Missing email",
			fields: map[string]string{
				"password":          "Password1!",
				"firstName":         "John",
				"lastName":          "Doe",
				"dob":               "2000-01-01",
				"nickname":          "johndoe",
				"aboutMe":           "Hello",
				"profileVisibility": "public",
			},
		},
		{
			name: "Missing password",
			fields: map[string]string{
				"email":             "test@example.com",
				"firstName":         "John",
				"lastName":          "Doe",
				"dob":               "2000-01-01",
				"nickname":          "johndoe",
				"aboutMe":           "Hello",
				"profileVisibility": "public",
			},
		},
		{
			name: "Missing firstName",
			fields: map[string]string{
				"email":             "test@example.com",
				"password":          "Password1!",
				"lastName":          "Doe",
				"dob":               "2000-01-01",
				"nickname":          "johndoe",
				"aboutMe":           "Hello",
				"profileVisibility": "public",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockAuthService := &MockAuthService{
				CreateUserFunc: func(user *models.User) (*models.User, error) {
					user.ID = 1
					return user, nil
				},
				ValidateEmailFunc: func(email string) (bool, error) {
					return true, nil
				},
				UserExistsFunc: func(email string) (bool, error) {
					return false, nil
				},
			}
			authHandler := handlers.NewAuthHandler(mockAuthService)

			body, contentType, err := createMultipartForm(tc.fields)
			if err != nil {
				t.Fatalf("Failed to create multipart form: %v", err)
			}

			req := httptest.NewRequest("POST", "/register", body)
			req.Header.Set("Content-Type", contentType)
			rr := httptest.NewRecorder()

			authHandler.Signup(rr, req)

			// Missing fields should still be handled gracefully
			// The service layer will handle validation
			if status := rr.Code; status != http.StatusOK && status != http.StatusBadRequest {
				t.Errorf("expected status %d or %d for missing field, got %d",
					http.StatusOK, http.StatusBadRequest, status)
			}
		})
	}
}

// Test password hashing and storage in the users table via handlers.Signup
func TestSignup_PasswordHashing(t *testing.T) {
	email := "test.password@example.com"
	password := "SecurePassword123!"

	// Mock AuthService that stores the password for inspection
	var storedPassword string
	mockAuthService := &MockAuthService{
		CreateUserFunc: func(user *models.User) (*models.User, error) {
			if user.Password == "" {
				t.Fatal("Password should not be empty")
			}
			// Simulate password hashing like the real service does
			if user.Password == password {
				// Hash the password using bcrypt (simulating the real service)
				hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
				if err != nil {
					return nil, err
				}
				user.Password = string(hashedPassword)
			}
			storedPassword = user.Password
			user.ID = 1
			return user, nil
		},
		ValidateEmailFunc: func(email string) (bool, error) {
			return true, nil
		},
		UserExistsFunc: func(email string) (bool, error) {
			return false, nil
		},
	}
	authHandler := handlers.NewAuthHandler(mockAuthService)

	fields := map[string]string{
		"email":             email,
		"password":          password,
		"firstName":         "John",
		"lastName":          "Doe",
		"dob":               "2000-01-01",
		"nickname":          "johndoe",
		"aboutMe":           "Hello",
		"profileVisibility": "public",
	}

	body, contentType, err := createMultipartForm(fields)
	if err != nil {
		t.Fatalf("Failed to create multipart form: %v", err)
	}

	req := httptest.NewRequest("POST", "/register", body)
	req.Header.Set("Content-Type", contentType)
	rr := httptest.NewRecorder()

	authHandler.Signup(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	// Verify the password is not stored in plain text
	if storedPassword == password {
		t.Error("Password is stored in plain text, should be hashed")
	}

	// Verify the password is properly hashed (should start with bcrypt identifier $2)
	if !strings.HasPrefix(storedPassword, "$2") {
		t.Errorf("Password does not appear to be hashed with bcrypt: %s", storedPassword)
	}

	// Verify we can validate the password against the hash
	passwordManager := utils.NewPasswordManager(utils.PasswordConfig{})
	err = passwordManager.ComparePassword(storedPassword, password)
	if err != nil {
		t.Errorf("Failed to validate password against hash: %v", err)
	}

	// Verify wrong password fails validation
	err = passwordManager.ComparePassword(storedPassword, "WrongPassword123!")
	if err == nil {
		t.Error("Wrong password validated successfully, should have failed")
	}
}
