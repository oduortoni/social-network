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
	mux := http.NewServeMux()

	permissionChecker := ws.NewDBPermissionChecker(db)
	wsManager := ws.NewManager(
		ws.NewDBSessionResolver(db),
		ws.NewDBGroupMemberFetcher(db),
		ws.NewDBMessagePersister(db),
		permissionChecker,
	)

	mux.Handle("GET /ws", middleware.AuthMiddleware(db)(http.HandlerFunc(handlers.NewWebSocketHandler(wsManager).HandleConnection)))

	notifier := ws.NewDBNotificationSender(wsManager)
	chatHandler := ws.NewChatHandler(
		db,
		ws.NewDBSessionResolver(db),
		ws.NewDBMessagePersister(db),
		notifier,
		wsManager,
		permissionChecker,
	)

	postStore := store.NewPostStore(db)
	authStore := store.NewAuthStore(db)
	followStore := store.NewFollowStore(db)
	unfollowstore := store.NewUnfollowStore(db)
	followRequestStore := store.NewFollowRequestStore(db)
	reactionStore := store.NewReactionStore(db)
	profilestore := store.NewProfileStore(db)
	groupStore := store.NewGroupStore(db)
	groupRequestStore := store.NewGroupRequestStore(db)
	groupChatMessageStore := store.NewGroupChatMessageStore(db)
	groupMemberStore := store.NewGroupMemberStore(db)

	postService := service.NewPostService(postStore)
	authService := service.NewAuthService(authStore)
	followService := service.NewFollowService(followStore)
	unfollowService := service.NewUnfollowService(unfollowstore)
	followRequestService := service.NewFollowRequestService(followRequestStore)
	reactionService := service.NewReactionService(reactionStore)
	profileService := service.NewProfileService(profilestore)
	groupService := service.NewGroupService(groupStore)
	groupRequestService := service.NewGroupRequestService(groupRequestStore, groupService)
	groupChatMessageService := service.NewGroupChatMessageService(groupChatMessageStore, groupService, groupMemberStore)

	postHandler := handlers.NewPostHandler(postService)
	authHandler := handlers.NewAuthHandler(authService)
	followHandler := handlers.NewFollowHandler(followService, notifier)
	unfollowHandler := handlers.NewUnfollowHandler(unfollowService)
	followRequestHandler := handlers.NewFollowRequestHandler(followRequestService, notifier)
	reactionHandler := handlers.NewReactionHandler(reactionService)
	profileHandler := handlers.NewProfileHandler(profileService)
	groupHandler := handlers.NewGroupHandler(groupService, groupRequestService, groupChatMessageService)

	mux.HandleFunc("POST /validate/step1", authHandler.ValidateAccountStepOne)
	mux.HandleFunc("POST /register", authHandler.Signup)
	mux.HandleFunc("POST /login", authHandler.Login)
	mux.HandleFunc("POST /logout", func(w http.ResponseWriter, r *http.Request) {
		authHandler.LogoutHandler(w, r)
	})

	mux.Handle("POST /groups", middleware.AuthMiddleware(db)(http.HandlerFunc(groupHandler.CreateGroup)))
	mux.Handle("POST /groups/{groupID}/join-request", middleware.AuthMiddleware(db)(http.HandlerFunc(groupHandler.SendJoinRequest)))
	mux.Handle("PUT /groups/{groupID}/join-request/{requestID}/approve", middleware.AuthMiddleware(db)(http.HandlerFunc(groupHandler.ApproveJoinRequest)))
	mux.Handle("PUT /groups/{groupID}/join-request/{requestID}/reject", middleware.AuthMiddleware(db)(http.HandlerFunc(groupHandler.RejectJoinRequest)))

	mux.Handle("POST /groups/{groupID}/chat", middleware.AuthMiddleware(db)(http.HandlerFunc(groupHandler.SendGroupChatMessage)))
	mux.Handle("GET /groups/{groupID}/chat", middleware.AuthMiddleware(db)(http.HandlerFunc(groupHandler.GetGroupChatMessages)))
	mux.Handle("POST /posts", middleware.AuthMiddleware(db)(http.HandlerFunc(postHandler.CreatePost)))
	mux.Handle("GET /posts/{postId}", middleware.AuthMiddleware(db)(http.HandlerFunc(postHandler.GetPostByID)))
	mux.Handle("GET /posts", middleware.AuthMiddleware(db)(http.HandlerFunc(postHandler.GetPosts)))
	mux.Handle("PUT /posts/{postId}", middleware.AuthMiddleware(db)(http.HandlerFunc(postHandler.UpdatePost)))
	mux.Handle("POST /posts/{postId}/comments", middleware.AuthMiddleware(db)(http.HandlerFunc(postHandler.CreateComment)))
	mux.Handle("GET /posts/{postId}/comments", middleware.AuthMiddleware(db)(http.HandlerFunc(postHandler.GetCommentsByPostID)))
	mux.Handle("PUT /posts/{postId}/comments/{commentId}", middleware.AuthMiddleware(db)(http.HandlerFunc(postHandler.UpdateComment)))
	mux.Handle("DELETE /posts/{postId}/comments/{commentId}", middleware.AuthMiddleware(db)(http.HandlerFunc(postHandler.DeleteComment)))
	mux.Handle("DELETE /posts/{postId}", middleware.AuthMiddleware(db)(http.HandlerFunc(postHandler.DeletePost)))
	mux.Handle("GET /users/search", middleware.AuthMiddleware(db)(http.HandlerFunc(postHandler.SearchUsers)))

	mux.Handle("POST /posts/{postId}/reaction", middleware.AuthMiddleware(db)(http.HandlerFunc(reactionHandler.ReactToPost)))
	mux.Handle("DELETE /posts/{postId}/reaction", middleware.AuthMiddleware(db)(http.HandlerFunc(reactionHandler.UnreactToPost)))
	mux.Handle("POST /comments/{commentId}/reaction", middleware.AuthMiddleware(db)(http.HandlerFunc(reactionHandler.ReactToComment)))
	mux.Handle("DELETE /comments/{commentId}/reaction", middleware.AuthMiddleware(db)(http.HandlerFunc(reactionHandler.UnreactToComment)))

	mux.Handle("POST /follow", middleware.AuthMiddleware(db)(http.HandlerFunc(followHandler.Follow)))
	mux.Handle("DELETE /unfollow", middleware.AuthMiddleware(db)(http.HandlerFunc(unfollowHandler.Unfollow)))
	mux.Handle("POST /follow-request/{requestId}/request", middleware.AuthMiddleware(db)(http.HandlerFunc(followRequestHandler.FollowRequestRespond)))
	mux.Handle("DELETE /follow-request/{requestId}/cancel", middleware.AuthMiddleware(db)(http.HandlerFunc(followRequestHandler.CancelFollowRequest)))
	mux.Handle("GET /pending-follow-requests", middleware.AuthMiddleware(db)(http.HandlerFunc(followRequestHandler.GetPendingFollowRequest)))

	mux.Handle("GET /profile/{userid}", middleware.AuthMiddleware(db)(http.HandlerFunc(profileHandler.ProfileHandler)))
	mux.Handle("GET /profile/{userid}/followers", middleware.AuthMiddleware(db)(http.HandlerFunc(profileHandler.GetFollowers)))
	mux.Handle("GET /profile/{userid}/followees", middleware.AuthMiddleware(db)(http.HandlerFunc(profileHandler.GetFollowees)))
	mux.Handle("PUT /EditProfile", middleware.AuthMiddleware(db)(http.HandlerFunc(authHandler.EditProfile))) // Edit profile handler

	mux.Handle("GET /me", middleware.AuthMiddleware(db)(http.HandlerFunc(handlers.NewMeHandler(db))))
	mux.Handle("GET /avatar", http.HandlerFunc(handlers.GetImage))

	mux.Handle("GET /api/messages/private", middleware.AuthMiddleware(db)(http.HandlerFunc(chatHandler.GetPrivateMessages)))
	mux.Handle("GET /api/messages/group", middleware.AuthMiddleware(db)(http.HandlerFunc(chatHandler.GetGroupMessages)))
	mux.Handle("POST /api/groups/invite", middleware.AuthMiddleware(db)(http.HandlerFunc(chatHandler.SendGroupInvite)))
	mux.Handle("GET /api/notifications", middleware.AuthMiddleware(db)(http.HandlerFunc(chatHandler.GetNotifications)))
	mux.Handle("POST /api/notifications/read", middleware.AuthMiddleware(db)(http.HandlerFunc(chatHandler.MarkNotificationsRead)))
	mux.Handle("GET /api/users/online", middleware.AuthMiddleware(db)(http.HandlerFunc(chatHandler.GetOnlineUsers)))

	mux.Handle("GET /api/users/messageable", middleware.AuthMiddleware(db)(http.HandlerFunc(chatHandler.GetMessageableUsers)))

	return mux
}
