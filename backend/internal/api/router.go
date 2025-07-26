package api

import (
	"database/sql"
	"net/http"

	"github.com/tajjjjr/social-network/backend/internal/api/handlers"
	"github.com/tajjjjr/social-network/backend/internal/api/middleware"
	"github.com/tajjjjr/social-network/backend/internal/service"
	"github.com/tajjjjr/social-network/backend/internal/store"
	ws "github.com/tajjjjr/social-network/backend/internal/websocket"
)

func NewRouter(db *sql.DB) http.Handler {
	// Create router
	mux := http.NewServeMux()

	// create the websocket handler
	wsManager := ws.NewManager(
		ws.NewDBSessionResolver(db),
		ws.NewDBGroupMemberFetcher(db),
		ws.NewDBMessagePersister(db),
	)
	mux.Handle("GET /ws", middleware.AuthMiddleware(db)(http.HandlerFunc(wsManager.HandleConnection)))

	// Chat history handlers (paginated HTTP access to messages)
	notifier := ws.NewDBNotificationSender(wsManager)
	chatHandler := ws.NewChatHandler(
		db,
		ws.NewDBSessionResolver(db),
		ws.NewDBMessagePersister(db),
		notifier,
		wsManager,
	)

	// Create stores
	postStore := store.NewPostStore(db)
	authStore := store.NewAuthStore(db)
	followStore := store.NewFollowStore(db)
	unfollowstore := store.NewUnfollowStore(db)
	followRequestStore := store.NewFollowRequestStore(db)

	// Create services
	postService := service.NewPostService(postStore)
	authService := service.NewAuthService(authStore)
	followService := service.NewFollowService(followStore)
	unfollowService := service.NewUnfollowService(unfollowstore)
	followRequestService := service.NewFollowRequestService(followRequestStore)

	// Create handlers
	postHandler := handlers.NewPostHandler(postService)
	authHandler := handlers.NewAuthHandler(authService)
	followHandler := handlers.NewFollowHandler(followService, notifier)
	unfollowHandler := handlers.NewUnfollowHandler(unfollowService)
	followRequestHandler := handlers.NewFollowRequestHandler(followRequestService, notifier)

	// Authentication Handlers
	mux.HandleFunc("POST /validate/step1", authHandler.ValidateAccountStepOne)
	mux.HandleFunc("POST /register", authHandler.Signup)
	mux.HandleFunc("POST /login", authHandler.Login)
	mux.HandleFunc("POST /logout", func(w http.ResponseWriter, r *http.Request) {
		authHandler.LogoutHandler(w, r)
	})

	// mux.HandleFunc("GET /checksession", func(w http.ResponseWriter, r *http.Request) {
	// 	authentication.CheckSessionHandler(w, r, db)
	// })

	// Mount handlers
	mux.Handle("POST /posts", middleware.AuthMiddleware(db)(http.HandlerFunc(postHandler.CreatePost)))
	mux.Handle("GET /posts/{postId}", middleware.AuthMiddleware(db)(http.HandlerFunc(postHandler.GetPostByID)))
	mux.Handle("GET /posts", middleware.AuthMiddleware(db)(http.HandlerFunc(postHandler.GetPosts)))
	mux.Handle("POST /posts/{postId}/comments", middleware.AuthMiddleware(db)(http.HandlerFunc(postHandler.CreateComment)))
	mux.Handle("GET /posts/{postId}/comments", middleware.AuthMiddleware(db)(http.HandlerFunc(postHandler.GetCommentsByPostID)))
	mux.Handle("DELETE /posts/{postId}", middleware.AuthMiddleware(db)(http.HandlerFunc(postHandler.DeletePost)))

	mux.Handle("POST /follow", middleware.AuthMiddleware(db)(http.HandlerFunc(followHandler.Follow)))
	mux.Handle("DELETE /unfollow", middleware.AuthMiddleware(db)(http.HandlerFunc(unfollowHandler.Unfollow)))
	mux.Handle("POST /follow-request/{requestId}/request", middleware.AuthMiddleware(db)(http.HandlerFunc(followRequestHandler.FollowRequestRespond)))

	mux.Handle("GET /me", middleware.AuthMiddleware(db)(http.HandlerFunc(handlers.NewMeHandler(db))))
	mux.Handle("GET /avatar", http.HandlerFunc(handlers.Avatar))

	/*example websocket routes */
	mux.Handle("GET /api/messages/private", middleware.AuthMiddleware(db)(http.HandlerFunc(chatHandler.GetPrivateMessages)))
	mux.Handle("GET /api/messages/group", middleware.AuthMiddleware(db)(http.HandlerFunc(chatHandler.GetGroupMessages)))
	mux.Handle("POST /api/groups/invite", middleware.AuthMiddleware(db)(http.HandlerFunc(chatHandler.SendGroupInvite)))
	mux.Handle("GET /api/notifications", middleware.AuthMiddleware(db)(http.HandlerFunc(chatHandler.GetNotifications)))
	mux.Handle("POST /api/notifications/read", middleware.AuthMiddleware(db)(http.HandlerFunc(chatHandler.MarkNotificationsRead)))
	mux.Handle("GET /api/users/online", middleware.AuthMiddleware(db)(http.HandlerFunc(chatHandler.GetOnlineUsers)))

	return mux
}
