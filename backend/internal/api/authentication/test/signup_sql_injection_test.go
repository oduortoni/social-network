package test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/tajjjjr/social-network/backend/internal/api/authentication"
)

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

func TestSignupHandler_SQLInjectionVariants(t *testing.T) {
	db := createTestDB(t)
	defer db.Close()
	injections := []string{
		"' OR '1'='1",
		"' OR 1=1--",
		"' OR 1=1#",
		"admin'--",
		"' UNION SELECT 1,2,3--",
		"' UNION SELECT email,password,id FROM Users--",
		"test@example.com' AND 1=1--",
		"test@example.com' AND (SELECT COUNT(*) FROM Users)>0--",
		"test@example.com'; DROP TABLE Users;--",
		"test@example.com'; INSERT INTO Users (email,password) VALUES ('hacker','hacked');--",
		"test@example.com'; UPDATE Users SET password='hacked' WHERE email='test@example.com';--",
		"test@example.com' AND sqlite_version()--",
		"test@example.com' OR 1=1 LIMIT 1--",
		"test@example.com' OR 1=1 LIMIT 1;--",
		"test@example.com' OR 1=1;--",
		"test@example.com' OR 1=1#",
		"test@example.com' OR 1=1/*",
		"test@example.com' OR 1=1-- -",
		"test@example.com' OR 1=1--+",
		"test@example.com' OR 1=1--%0A",
		"test@example.com' OR 1=1--%0D%0A",
		"test@example.com' OR 1=1--%23",
		"test@example.com' OR 1=1--%3B",
		"test@example.com' OR 1=1--%2F%2A",
		"test@example.com' OR 1=1--%2D%2D",
		"test@example.com' OR 1=1--%20",
		"test@example.com' OR 1=1--%09",
		"test@example.com' OR 1=1--%0B",
		"test@example.com' OR 1=1--%0C",
		"test@example.com' OR 1=1--%0D",
		"test@example.com' OR 1=1--%0A%0D",
		"test@example.com' OR 1=1--%0D%0A",
		"test@example.com' OR 1=1--%0A%0A",
		"test@example.com' OR 1=1--%0D%0D",
		"test@example.com' OR 1=1--%0A%0A%0A",
		"test@example.com' OR 1=1--%0D%0D%0D",
		"test@example.com' OR 1=1--%0A%0D%0A%0D",
		"test@example.com' OR 1=1--%0D%0A%0D%0A",
		"test@example.com' OR 1=1--%0A%0A%0A%0A",
		"test@example.com' OR 1=1--%0D%0D%0D%0D",
		"test@example.com' OR 1=1--%0A%0D%0A%0D%0A",
		"test@example.com' OR 1=1--%0D%0A%0D%0A%0D%0A",
		"test@example.com' OR 1=1--%0A%0A%0A%0A%0A",
		"test@example.com' OR 1=1--%0D%0D%0D%0D%0D",
		"test@example.com' OR 1=1--%0A%0D%0A%0D%0A%0D%0A",
		"test@example.com' OR 1=1--%0D%0A%0D%0A%0D%0A%0D%0A",
		"test@example.com' OR 1=1--%0A%0A%0A%0A%0A%0A",
		"test@example.com' OR 1=1--%0D%0D%0D%0D%0D%0D",
		"test@example.com' OR 1=1--%0A%0D%0A%0D%0A%0D%0A%0D%0A%0D%0A",
		"test@example.com' OR 1=1--%0D%0A%0D%0A%0D%0A%0D%0A%0D%0A",
	}
	for _, inj := range injections {
		fields := map[string]string{
			"email":             inj,
			"password":          "Password1!",
			"firstName":         "John",
			"lastName":          "Doe",
			"dob":               "2000-01-01",
			"nickname":          "hacker",
			"aboutMe":           "test",
			"profileVisibility": "public",
		}
		body, contentType, _ := createMultipartForm(fields)
		req := httptest.NewRequest("POST", "/register", body)
		req.Header.Set("Content-Type", contentType)
		w := httptest.NewRecorder()
		authentication.SignupHandler(w, req, db)
		resp := w.Result()
		if resp.StatusCode == http.StatusOK {
			t.Errorf("SQL injection variant '%s' should not succeed", inj)
		}
	}
}
