package authentication

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type Response struct {
	Message string `json:"message,omitempty"`
}

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

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
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
