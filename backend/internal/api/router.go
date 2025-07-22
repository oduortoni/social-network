package api

import (
	"database/sql"
	"net/http"

	"github.com/tajjjjr/social-network/backend/internal/api/authentication"
	"github.com/tajjjjr/social-network/backend/internal/api/handlers"
	"github.com/tajjjjr/social-network/backend/internal/api/middleware"
	"github.com/tajjjjr/social-network/backend/internal/service"
	"github.com/tajjjjr/social-network/backend/internal/store"
)

func NewRouter(db *sql.DB) http.Handler {
	// Create stores
	postStore := store.NewPostStore(db)
	authStore := store.NewAuthStore(db)

	// Create services
	postService := service.NewPostService(postStore)
	authService := service.NewAuthService(authStore)

	// Create handlers
	postHandler := handlers.NewPostHandler(postService)
	authHandler := handlers.NewAuthHandler(authService)

	// Create router
	mux := http.NewServeMux()

	// Authentication Handlers
	mux.HandleFunc("POST /validate/step1", func(w http.ResponseWriter, r *http.Request) {
		authentication.ValidateAccountStepOne(w, r, db)
	})
	mux.HandleFunc("POST /register", func(w http.ResponseWriter, r *http.Request) {
		authentication.SignupHandler(w, r, db)
	})
	mux.HandleFunc("POST /login", authHandler.Login)

	mux.HandleFunc("POST /logout", func(w http.ResponseWriter, r *http.Request) {
		authentication.LogoutHandler(w, r, db)
	})
	mux.HandleFunc("GET /auth/google/login", authentication.RedirectToGoogleLogin)
	mux.HandleFunc("GET /auth/google/callback", func(w http.ResponseWriter, r *http.Request) {
		authentication.HandleGoogleCallback(w, r, db)
	})
	mux.HandleFunc("GET /auth/facebook/login", authentication.RedirectToFacebookLogin)
	mux.HandleFunc("GET /auth/facebook/callback", func(w http.ResponseWriter, r *http.Request) {
		authentication.HandleFacebookCallback(w, r, db)
	})
	mux.HandleFunc("GET /auth/github/login", authentication.RedirectToGitHubLogin)
	mux.HandleFunc("GET /auth/github/callback", func(w http.ResponseWriter, r *http.Request) {
		authentication.HandleGitHubCallback(w, r, db)
	})
	mux.HandleFunc("GET /checksession", func(w http.ResponseWriter, r *http.Request) {
		authentication.CheckSessionHandler(w, r, db)
	})

	// Mount handlers
	mux.Handle("POST /posts", middleware.AuthMiddleware(db)(http.HandlerFunc(postHandler.CreatePost)))
	mux.Handle("GET /posts/{postId}", middleware.AuthMiddleware(db)(http.HandlerFunc(postHandler.GetPostByID)))
	mux.Handle("GET /feed", middleware.AuthMiddleware(db)(http.HandlerFunc(postHandler.GetFeed)))
	mux.Handle("POST /posts/{postId}/comments", middleware.AuthMiddleware(db)(http.HandlerFunc(postHandler.CreateComment)))

	mux.Handle("GET /me", middleware.AuthMiddleware(db)(http.HandlerFunc(handlers.NewMeHandler(db))))
	mux.Handle("GET /avatar", http.HandlerFunc(handlers.Avatar))

	return mux
}
