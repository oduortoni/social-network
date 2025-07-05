package authentication

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CompareHashpassword(hashedpassword, userinputpassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedpassword), []byte(userinputpassword))
	return err == nil
}
