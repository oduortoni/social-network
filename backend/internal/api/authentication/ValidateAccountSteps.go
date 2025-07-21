package authentication

import (
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"regexp"

	"github.com/tajjjjr/social-network/backend/utils"
)

type StepOneCredintial struct {
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

func ValidateAccountStepOne(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var serverresponse Response
	statusCode := http.StatusOK
	var AccountCrediential StepOneCredintial

	body, err := io.ReadAll(r.Body)
	if err != nil {
		serverresponse.Message = "Failed to read request body"
		statusCode = http.StatusInternalServerError
		respondJSON(w, statusCode, serverresponse)
		return
	}
	if err := json.Unmarshal(body, &AccountCrediential); err != nil {
		serverresponse.Message = "Failed to parse request body"
		statusCode = http.StatusBadRequest
		respondJSON(w, statusCode, serverresponse)
		return
	}

	// Validate required fields
	if AccountCrediential.Email == "" || AccountCrediential.Password == "" || AccountCrediential.ConfirmPassword == "" {
		serverresponse.Message = "Missing required fields"
		statusCode = http.StatusBadRequest
		respondJSON(w, statusCode, serverresponse)
		return
	}

	// validate email
	IsEmailValid, err := ValidateEmail(AccountCrediential.Email)
	if err != nil {
		serverresponse.Message = "Regex error in validating Email"
		statusCode = http.StatusInternalServerError
		respondJSON(w, statusCode, serverresponse)
		return
	}
	if !IsEmailValid {
		serverresponse.Message = "Invalid Email format"
		statusCode = http.StatusBadRequest
		respondJSON(w, statusCode, serverresponse)
		return
	}

	// password is same as confirm passwod
	if AccountCrediential.Password != AccountCrediential.ConfirmPassword {
		serverresponse.Message = "Passwords do not match."
		statusCode = http.StatusBadRequest
		respondJSON(w, statusCode, serverresponse)
		return
	}

	// Check if user already exists
	if UserExists(AccountCrediential.Email, db) {
		serverresponse.Message = "Email already exists"
		statusCode = http.StatusConflict
		respondJSON(w, statusCode, serverresponse)
		return
	}
	// validate password
	passwordManager := utils.NewPasswordManager(utils.PasswordConfig{})
	_, err = passwordManager.HashPassword(AccountCrediential.Password)
	if err != nil {
		serverresponse.Message = err.Error()
		statusCode = http.StatusBadRequest
		respondJSON(w, statusCode, serverresponse)
		return
	}

	serverresponse.Message = "Ok"
	respondJSON(w, statusCode, serverresponse)
}

func ValidateEmail(email string) (bool, error) {
	emailPattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re, err := regexp.Compile(emailPattern)
	if err != nil {
		return false, err 
	}
	return re.MatchString(email), nil
}
