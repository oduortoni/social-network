package main

import (
	"fmt"
	"net/http"

	"github.com/tajjjjr/social-network/backend/internal/api/authentication"
	"github.com/tajjjjr/social-network/backend/pkg/utils"
	"github.com/tajjjjr/social-network/backend/www/controllers"
)

var (
	Host = "0.0.0.0"
	Port = 9000
)

func main() {
	Port := utils.Port(Port)
	srvAddr := fmt.Sprintf("%s:%d", Host, Port)
	mux := http.NewServeMux()
	mux.HandleFunc("/", controllers.Index)

	// Authentication Handlers
	mux.HandleFunc("/register", authentication.SignupHandler)
	mux.HandleFunc("/login", authentication.SigninHandler)
	mux.HandleFunc("/logout", authentication.LogoutHandler)
	// Google Authentication
	http.HandleFunc("/auth/google/login", authentication.RedirectToGoogleLogin)
	http.HandleFunc("/auth/google/callback", func(w http.ResponseWriter, r *http.Request) {
		authentication.HandleGoogleCallback(w, r, db)
	})

	fmt.Printf("\n\n\n\t-----------[ server running on http://%s]-------------\n\n", srvAddr)
	http.ListenAndServe(srvAddr, controllers.CORSMiddleware(mux))
}