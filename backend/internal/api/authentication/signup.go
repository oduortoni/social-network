package authentication

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/tajjjjr/social-network/backend/pkg/utils"
)

type Profile_User struct {
	Email           string `json:"email"`
	Password        string `json:"password"`
	FirstName       string `json:"firstname"`
	LastName        string `json:"lastname"`
	DateOfBirth     string `json:"dateofbirth"`
	Nickname        string `json:"nickname"`
	AboutMe         string `json:"aboutme"`
	IsProfilePublic bool   `json:"isprofilepublic"`
	Avatar          string `json:"avatar"`
}

func SignupHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var serverresponse Response
	statusCode := http.StatusOK

	// Parse multipart form (limit: 10MB)
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		serverresponse.Message = "Failed to parse form"
		statusCode = http.StatusBadRequest
		respondJSON(w, statusCode, serverresponse)
		return
	}

	// Extract form values
	email := r.FormValue("email")
	password := r.FormValue("password")
	firstName := r.FormValue("first_name")
	lastName := r.FormValue("last_name")
	dateOfBirth := r.FormValue("date_of_birth")
	nickname := r.FormValue("nickname")
	aboutMe := r.FormValue("about_me")
	isProfilePublic := r.FormValue("is_profile_public")

	// Validate required fields
	if email == "" || password == "" {
		serverresponse.Message = "Missing required fields"
		statusCode = http.StatusBadRequest
		respondJSON(w, statusCode, serverresponse)
		return
	}

	// Check if user already exists
	if UserExists(email, nickname, db) {
		serverresponse.Message = "Email or nickname already taken"
		statusCode = http.StatusConflict
		respondJSON(w, statusCode, serverresponse)
		return
	}

	// Handle avatar upload
	userAvatar := "no profile photo"
	file, header, err := r.FormFile("avatar")
	if err == nil && file != nil {
		defer file.Close()
		userAvatar, err = UploadAvatarImage(file, header)
		if err != nil {
			serverresponse.Message = userAvatar
			statusCode = http.StatusInternalServerError
			respondJSON(w, statusCode, serverresponse)
			return
		}
	}

	// Hash the password
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		serverresponse.Message = "Failed to secure password"
		statusCode = http.StatusInternalServerError
		respondJSON(w, statusCode, serverresponse)
		return
	}

	// Create new user struct
	newUser := Profile_User{
		Email:           email,
		Password:        hashedPassword,
		FirstName:       firstName,
		LastName:        lastName,
		DateOfBirth:     dateOfBirth,
		Nickname:        nickname,
		AboutMe:         aboutMe,
		IsProfilePublic: isProfilePublic == "true",
		Avatar:          userAvatar,
	}

	// Insert user into DB
	if err := InsertUserIntoDB(newUser, db); err != nil {
		serverresponse.Message = "Failed to create user"
		statusCode = http.StatusInternalServerError
		respondJSON(w, statusCode, serverresponse)
		return
	}

	// Success
	serverresponse.Message = "Registration successful"
	statusCode = http.StatusOK
	respondJSON(w, statusCode, serverresponse)
}

func UserExists(email, nickname string, db *sql.DB) bool {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM Users WHERE email = ? OR nickname = ?", email, nickname).Scan(&count)
	if err != nil {
		log.Println(err)
		return false // return error if something goes wrong
	}
	return count > 0 // return true if user exists
}

func InsertUserIntoDB(user Profile_User, db *sql.DB) error {
	_, err := db.Exec("INSERT INTO Users (email,password,first_name,last_name,date_of_birth,nickname,about_me,is_profile_public,avatar,created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?,?)", user.Email, user.Password, user.FirstName, user.LastName, user.DateOfBirth, user.Nickname, user.AboutMe, user.IsProfilePublic, user.Avatar, time.Now())
	return err
}
