package test

import (
	"bytes"
	"database/sql"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	_ "github.com/mattn/go-sqlite3"

	"github.com/tajjjjr/social-network/backend/internal/api/authentication"
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

func TestSignupHandler_SQLInjectionAttempt(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()
	fields := map[string]string{
		"email":             "' OR 1=1;--",
		"password":          "Password1!",
		"firstName":         "John",
		"lastName":          "Doe",
		"dob":               "2000-01-01",
		"nickname":          "hacker",
		"aboutMe":           "<script>alert(1)</script>",
		"profileVisibility": "public",
	}
	body, contentType, _ := createMultipartForm(fields)
	req := httptest.NewRequest("POST", "/register", body)
	req.Header.Set("Content-Type", contentType)
	w := httptest.NewRecorder()
	authentication.SignupHandler(w, req, db)
	resp := w.Result()
	if resp.StatusCode == http.StatusOK {
		t.Error("SQL injection attempt should not succeed")
	}
}

func TestSignupHandler_XSSPrevention(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()

	// Test with potentially malicious XSS content
	fields := map[string]string{
		"email":             "test@example.com",
		"password":          "Password1!",
		"firstName":         "<script>alert('xss')</script>",
		"lastName":          "<img src=x onerror=alert(1)>",
		"dob":               "2000-01-01",
		"nickname":          "&lt;script&gt;",
		"aboutMe":           "Hello <b>world</b> & <script>alert('xss')</script>",
		"profileVisibility": "public",
	}

	body, contentType, err := createMultipartForm(fields)
	if err != nil {
		t.Fatalf("Failed to create multipart form: %v", err)
	}

	req := httptest.NewRequest("POST", "/register", body)
	req.Header.Set("Content-Type", contentType)
	w := httptest.NewRecorder()

	// Call the handler
	authentication.SignupHandler(w, req, db)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}

	// Verify that the data was HTML escaped in the database
	var firstName, lastName, nickname, aboutMe string
	err = db.QueryRow("SELECT first_name, last_name, nickname, about_me FROM Users WHERE email = ?", "test@example.com").
		Scan(&firstName, &lastName, &nickname, &aboutMe)
	if err != nil {
		t.Fatalf("Failed to query user: %v", err)
	}

	// Check that HTML has been escaped
	expectedFirstName := "&lt;script&gt;alert(&#39;xss&#39;)&lt;/script&gt;"
	expectedLastName := "&lt;img src=x onerror=alert(1)&gt;"
	expectedNickname := "&amp;lt;script&amp;gt;"
	expectedAboutMe := "Hello &lt;b&gt;world&lt;/b&gt; &amp; &lt;script&gt;alert(&#39;xss&#39;)&lt;/script&gt;"

	if firstName != expectedFirstName {
		t.Errorf("firstName not properly escaped. Expected: %s, Got: %s", expectedFirstName, firstName)
	}
	if lastName != expectedLastName {
		t.Errorf("lastName not properly escaped. Expected: %s, Got: %s", expectedLastName, lastName)
	}
	if nickname != expectedNickname {
		t.Errorf("nickname not properly escaped. Expected: %s, Got: %s", expectedNickname, nickname)
	}
	if aboutMe != expectedAboutMe {
		t.Errorf("aboutMe not properly escaped. Expected: %s, Got: %s", expectedAboutMe, aboutMe)
	}
}
