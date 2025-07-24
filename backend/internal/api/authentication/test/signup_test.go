package test

import (
	"bytes"
	"database/sql"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	_ "github.com/mattn/go-sqlite3"

	"github.com/tajjjjr/social-network/backend/internal/api/authentication"
	"github.com/tajjjjr/social-network/backend/utils"
)

// createTestDB creates a temporary in-memory SQLite DB for testing
func createTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to create test DB: %v", err)
	}
	// Create Users table schema for tests
	_, err = db.Exec(`CREATE TABLE Users (
		email TEXT PRIMARY KEY,
		password TEXT,
		first_name TEXT,
		last_name TEXT,
		date_of_birth TEXT,
		nickname TEXT,
		about_me TEXT,
		is_profile_public BOOLEAN,
		avatar TEXT,
		created_at DATETIME
	)`)
	if err != nil {
		db.Close()
		t.Fatalf("Failed to create Users table: %v", err)
	}
	return db
}

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

func TestSignupHandler_Success(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

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
	w := httptest.NewRecorder()
	authentication.SignupHandler(w, req, db)
	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}
}

func TestSignupHandler_UserAlreadyExists(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()
	// Insert a user first
	_, err := db.Exec(`INSERT INTO Users (email,password,first_name,last_name,date_of_birth,nickname,about_me,is_profile_public,avatar,created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, datetime('now'))`,
		"existing@example.com", "Password1!", "John", "Doe", "2000-01-01", "johndoe", "Hello", true, "avatar.png")
	if err != nil {
		t.Fatalf("Failed to insert user: %v", err)
	}
	fields := map[string]string{
		"email":             "existing@example.com",
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
	w := httptest.NewRecorder()
	authentication.SignupHandler(w, req, db)
	resp := w.Result()
	if resp.StatusCode != http.StatusConflict {
		t.Errorf("expected status %d, got %d", http.StatusConflict, resp.StatusCode)
	}
}

func TestSignupHandler_InvalidFormData(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()
	form := "invalid-form-data"
	req := httptest.NewRequest("POST", "/register", bytes.NewBufferString(form))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	authentication.SignupHandler(w, req, db)
	resp := w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, resp.StatusCode)
	}
}

func TestUserExists_UserDoesNotExist(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()
	exists := authentication.UserExists("nonexistent@example.com", db)
	if exists {
		t.Error("expected user to not exist")
	}
}

func TestUserExists_UserExists(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()
	_, err := db.Exec(`INSERT INTO Users (email,password,first_name,last_name,date_of_birth,nickname,about_me,is_profile_public,avatar,created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, datetime('now'))`,
		"existing@example.com", "Password1!", "John", "Doe", "2000-01-01", "johndoe", "Hello", true, "avatar.png")
	if err != nil {
		t.Fatalf("Failed to insert user: %v", err)
	}
	exists := authentication.UserExists("existing@example.com", db)
	if !exists {
		t.Error("expected user to exist")
	}
}

// Test registration with missing required fields
func TestSignupHandler_MissingFields(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

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
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			body, contentType, err := createMultipartForm(tc.fields)
			if err != nil {
				t.Fatalf("Failed to create multipart form: %v", err)
			}

			req := httptest.NewRequest("POST", "/register", body)
			req.Header.Set("Content-Type", contentType)
			w := httptest.NewRecorder()
			authentication.SignupHandler(w, req, db)
			resp := w.Result()

			// Should return bad request for missing required fields (email/password)
			if resp.StatusCode != http.StatusBadRequest && resp.StatusCode != http.StatusInternalServerError {
				t.Errorf("expected status %d or %d for missing field, got %d",
					http.StatusBadRequest, http.StatusInternalServerError, resp.StatusCode)
			}
		})
	}
}

// Test registration with invalid email formats
func TestSignupHandler_InvalidEmailFormat(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

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
			w := httptest.NewRecorder()
			authentication.SignupHandler(w, req, db)
			resp := w.Result()

			// Should return bad request for invalid email format
			if resp.StatusCode != http.StatusBadRequest {
				t.Errorf("expected status %d for invalid email '%s', got %d",
					http.StatusBadRequest, email, resp.StatusCode)
			}
		})
	}
}

// Test password hashing and storage in the users table
func TestSignupHandler_PasswordHashing(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	// Test data
	email := "test.password@example.com"
	password := "SecurePassword123!"

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
	w := httptest.NewRecorder()
	authentication.SignupHandler(w, req, db)
	resp := w.Result()

	// Verify successful registration
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}

	// Retrieve the stored password from the database
	var storedPassword string
	err = db.QueryRow("SELECT password FROM Users WHERE email = ?", email).Scan(&storedPassword)
	if err != nil {
		t.Fatalf("Failed to retrieve user from database: %v", err)
	}

	// Verify the password is not stored in plain text
	if storedPassword == password {
		t.Error("Password is stored in plain text, should be hashed")
	}

	// Verify the password is properly hashed (should start with bcrypt identifier $2a$ or similar)
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
