Here is a potential list of API endpoints. This list follows RESTful conventions and would form the basis of the API contract.

### **Authentication**
* `POST /register`
    * Creates a new user account.
* `POST /login`
    * Authenticates a user and creates a session.
* `POST /logout`
    * Invalidates the user's current session.

### **Users & Profiles**
* `GET /users/{userId}`
    * Retrieves the profile information and activity for a specific user.
* `GET /profile`
    * Retrieves the profile of the currently authenticated user.
* `PUT /profile`
    * Updates the profile information for the currently authenticated user.
* `GET /users/{userId}/followers`
    * Retrieves a list of users who are following the specified user.
* `GET /users/{userId}/following`
    * Retrieves a list of users that the specified user is following.

### **Followers**
* `POST /users/{userId}/follow`
    * Sends a follow request to a user or instantly follows a public profile.
* `DELETE /users/{userId}/follow`
    * Unfollows a user.
* `GET /follow-requests`
    * Retrieves a list of pending follow requests for the authenticated user.
* `POST /follow-requests/{requestId}/accept`
    * Accepts a pending follow request.
* `POST /follow-requests/{requestId}/decline`
    * Declines a pending follow request.

### **Posts & Comments**
* `POST /posts`
    * Creates a new post (with text, optional image/gif, and privacy settings).
* `GET /feed`
    * Retrieves the post feed for the authenticated user.
* `GET /posts/{postId}`
    * Retrieves a single post.
* `DELETE /posts/{postId}`
    * Deletes a post created by the authenticated user.
* `POST /posts/{postId}/comments`
    * Adds a new comment to a specific post.
* `GET /posts/{postId}/comments`
    * Retrieves all comments for a specific post.

### **Groups**
* `POST /groups`
    * Creates a new group.
* `GET /groups`
    * Retrieves a list of all groups for Browse.
* `GET /groups/{groupId}`
    * Retrieves the details for a specific group.
* `GET /groups/{groupId}/posts`
    * Retrieves all posts made within a specific group.
* `POST /groups/{groupId}/posts`
    * Creates a new post within a specific group.

### **Group Membership**
* `POST /groups/{groupId}/invites`
    * Invites a user to join a group.
* `POST /groups/{groupId}/join`
    * Sends a request to join a group.
* `GET /group-invitations`
    * Retrieves pending group invitations for the authenticated user.
* `POST /group-invitations/{invitationId}/accept`
    * Accepts a pending group invitation.
* `POST /group-invitations/{invitationId}/decline`
    * Declines a pending group invitation.
* `GET /groups/{groupId}/join-requests`
    * (For group creator) Retrieves pending join requests for a group.
* `POST /join-requests/{requestId}/accept`
    * (For group creator) Accepts a user's request to join.
* `POST /join-requests/{requestId}/decline`
    * (For group creator) Declines a user's request to join.

### **Group Events**
* `POST /groups/{groupId}/events`
    * Creates a new event within a group.
* `GET /events/{eventId}`
    * Retrieves the details for a specific event.
* `DELETE /events/{eventId}`
    * Deletes an event (if creator).
* `POST /events/{eventId}/respond`
    * Responds to an event (e.g., with "Going" or "Not Going").

### **Chat & Notifications**
* `GET /chats`
    * Retrieves the list of active chat conversations for the user.
* `GET /chats/{userId}/messages`
    * Retrieves the message history with a specific user.
* `GET /notifications`
    * Retrieves all notifications for the authenticated user.
* `POST /notifications/{notificationId}/read`
    * Marks a specific notification as read.

### **WebSocket Endpoint**
* `ws://your-api-domain.com/chat`
    * A WebSocket connection for real-time sending and receiving of chat messages and notifications.