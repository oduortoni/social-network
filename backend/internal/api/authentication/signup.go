package authentication

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/tajjjjr/social-network/backend/internal/api/handlers"
	"github.com/tajjjjr/social-network/backend/utils"
)

type Profile_User struct {
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

func SignupHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var serverresponse Response
	statusCode := http.StatusOK

	// Parse multipart form (limit: 10MB)
	if err := r.ParseMultipartForm(1 << 20); err != nil {
		serverresponse.Message = "Failed to parse form"
		statusCode = http.StatusBadRequest
		respondJSON(w, statusCode, serverresponse)
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

	// Check if user already exists
	if UserExists(email, db) {
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
		userAvatar, err = handlers.UploadAvatarImage(file, header)
		if err != nil {
			serverresponse.Message = userAvatar
			statusCode = http.StatusInternalServerError
			respondJSON(w, statusCode, serverresponse)
			return
		}
	}

	passwordManager := utils.NewPasswordManager(utils.PasswordConfig{})
	hashedPassword, err := passwordManager.HashPassword(password)
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
		IsProfilePublic: isProfilePublic == "public",
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
	respondJSON(w, statusCode, serverresponse)
}

func UserExists(email string, db *sql.DB) bool {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM Users WHERE email = ?", email).Scan(&count)
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
