package main

import (
	"fmt"
	"net/http"

	_ "github.com/mattn/go-sqlite3"

	"github.com/tajjjjr/social-network/backend/internal/api/authentication"
	"github.com/tajjjjr/social-network/backend/pkg/db/sqlite"
	"github.com/tajjjjr/social-network/backend/pkg/utils"
	"github.com/tajjjjr/social-network/backend/www/controllers"
)

var (
	Host = "0.0.0.0"
	Port = 9000
)

func main() {
	// Initialize DB and run migrations
	db, err := sqlite.Migration()
	if err != nil {
		panic(fmt.Sprintf("DB migration failed: %v", err))
	}
	defer db.Close()

	Port := utils.Port(Port)
	srvAddr := fmt.Sprintf("%s:%d", Host, Port)
	mux := http.NewServeMux()
	mux.HandleFunc("/", controllers.Index)

	// Authentication Handlers
	mux.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		authentication.SignupHandler(w, r, db)
	})
	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		authentication.SigninHandler(w, r, db)
	})
	mux.HandleFunc("/logout", authentication.LogoutHandler)
	// Google Authentication
	mux.HandleFunc("/auth/google/login", authentication.RedirectToGoogleLogin)
	mux.HandleFunc("/auth/google/callback", func(w http.ResponseWriter, r *http.Request) {
		authentication.HandleGoogleCallback(w, r, db)
	})

	fmt.Printf("\n\n\n\t-----------[ server running on http://%s]-------------\n\n", srvAddr)
	http.ListenAndServe(srvAddr, controllers.CORSMiddleware(mux))
}
