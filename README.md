# social-network

## Description
A Facebook-like social network with features like profiles, posts, groups, real-time chat, and notifications, built with a Go backend (SQLite, Docker, WebSockets, migrations, authentication) and a JavaScript frontend.

## Features
* **User Management:** Registration, Login, Logout, Sessions & Cookies, Optional profile fields (Avatar, Nickname, About Me).
* **Profiles:** Public/Private profiles, User information display, User activity feed, Followers/Following display, Toggle profile privacy.
* **Following System:** Send/accept/decline follow requests, Automatic following for public profiles.
* **Posts & Comments:** Create posts and comments with image/GIF support.
* **Post Privacy:** Public, Almost Private (followers only), Private (selected followers).
* **Groups:** Create groups with title/description, Invite users, Request to join groups, Browse all groups.
* **Group Activities:** Create posts and comments within groups, Create events (Title, Description, Day/Time, Going/Not Going options).
* **Real-time Chat:** Private messaging between followers/following users, Group chat rooms, Emoji support, Instant delivery via WebSockets.
* **Notifications:** Follow requests, Group invitations, Group join requests (for creators), New group events.

## API Contract

### Authentication

*   **POST /login**
    *   Description: Authenticates a user and creates a session.
    *   Request Body:
        ```json
        {
          "email": "user@example.com",
          "password": "password123"
        }
        ```
    * if successful, Response:
        ```json
        {
          "message": "Login successful"
        }
        ```

*   **POST /register**
    *   Description: Registers a new user.
    *   Request Body:
        ```json
        {
          "email": "user@example.com",
          "password": "password123",
          "firstname": "John",
          "lastname": "Doe",
          "dateofbirth": "1990-01-15",
          "nickname": "Johnny",
          "aboutme": "I'm a new user",
          "isprofilepublic": true,
          "avatar": "url-to-avatar"
        }
        ```
    *   if successful, Response:
        ```json
        { 
          "message": "Registration successful"
        }
        ```
* **POST /logout**
    *   Description: Invalidates the user's current session.
    *   if successful, Response:
        ```json
        { 
          "message": "Logout successful"
        }
        ```   

* **POST /validate/step1**
    *   Description: (Optional) First step of account validation (if required by your flow).
    *   Request Body: (see implementation)
    *   Response: (see implementation)

<!--
* **POST /checksession**
    *   Description: Checks if a session is valid. (Currently not active in backend)
    *   if successful, Response:
        ```json
        { 
          "message": "Valid session"
        }
        ```
-->

### Users

*   **GET /api/users/{id}**
    *   Description: Retrieves a user's profile information.
    *   Response:
        ```json
        {
          "id": 123,
          "username": "newuser",
          "email": "newuser@example.com",
          "profile": {
            "avatar": "url-to-avatar",
            "nickname": "Newbie",
            "about": "Hello, I'm new here!"
          }
        }
        ```

### Posts

*   **GET /api/posts**
    *   Description: Retrieves a feed of posts.
    *   Response:
        ```json
        [
          {
            "id": 1,
            "authorId": 123,
            "content": "This is my first post!",
            "createdAt": "2025-06-29T12:00:00Z"
          }
        ]
        ```

*   **POST /api/posts**
    *   Description: Creates a new post.
    *   Request Body:
        ```json
        {
          "content": "My new post content."
        }
        ```
    *   Response:
        ```json
        {
          "id": 2,
          "authorId": 123,
          "content": "My new post content.",
          "createdAt": "2025-06-29T13:00:00Z"
        }
        ```

## Requirements
To run this project, you will need:
* Next JS
* A Go environment for the backend.
* SQLite database.
* Docker.
* Docker Compose.

## Getting Started

Clone the repository and navigate to the frontend directory:

```bash
git clone https://github.com/oduortoni/social-network.git
cd social-network/frontend
npm install
```

**create a .env file**

This project runs a Go backend and a Next.js frontend at the same time. Each must use a different port to avoid conflicts. Using a .env file within the root of the frontend folder, populate it with the following content:

```env
PORT=9000
NEXT_PORT=3000

# development
NEXT_PUBLIC_API_URL=http://localhost:9000

# Authentication credentials (required for OAuth)
FACEBOOK_CLIENT_ID=your_facebook_client_id
FACEBOOK_CLIENT_SECRET=your_facebook_client_secret
GITHUB_CLIENT_ID=your_github_client_id
GITHUB_CLIENT_SECRET=your_github_client_secret
GOOGLE_CLIENT_ID=your_google_client_id
GOOGLE_CLIENT_SECRET=your_google_client_secret

# Note: When making API calls from React, include credentials
# Add this to your fetch options:
# credentials: 'include'

# production
# NEXT_PUBLIC_API_URL=https://api.example.com
```

the NEXT_PUBLIC_API_URL will be used by the browser automatically to access the backend url

finally, run both servers

```bash
npm run dev
```

This command starts both the backend and frontend concurrently using the correct port settings.

**open in your browser**

Follow the link that falls under the next js project i.e either of the last two urls

e.g  http://localhost:3000

since we defined the frontend to run using NEXT_PORT=3000

## Usage
To get started with this project, clone the repository:
```bash
npm run dev
```

This command starts both the backend and frontend concurrently using the correct port settings.

**open in your browser**

Follow the link that falls under the next js project i.e either of the last two urls

e.g  http://localhost:3000

since we defined the frontend to run using NEXT_PORT=3000

## Docker Integration
TO start the docker containers
```bash
docker-compose up --build
```
TO delete the docker containers
```bash
docker-compose down
```

## Authentication System Setup & Testing

### Setup
- Backend: Ensure Go and SQLite are installed. Run `go run server.go` in the backend directory, or use Docker Compose as described above.
- Frontend: See environment setup above. The frontend communicates with the backend via the API endpoints below.

### Running Authentication Tests
- From the backend directory, run:
  ```bash
  go test ./internal/api/handlers/tests/
  ```
- This will execute all authentication-related tests, including SQL injection and XSS prevention.

## Authentication API Endpoints

| Endpoint         | Method | Description                |
|------------------|--------|----------------------------|
| /api/register    | POST   | Register a new user        |
| /api/login       | POST   | Log in a user              |
| /api/logout      | POST   | Log out the current user   |

### /api/register
- **Request:** JSON or multipart/form-data with fields:
  - email, password, firstname, lastname, dateofbirth, nickname, aboutme, isprofilepublic, avatar
- **Success:** 200 OK `{ "message": "Registration successful" }`
- **Errors:** 400 (invalid input), 409 (email exists), 500 (server error)

### /api/login
- **Request:** JSON or form data with fields:
  - email, password
- **Success:** 200 OK `{ "message": "Login successful" }` and sets `session_id` cookie
- **Errors:** 400 (invalid input), 401 (invalid credentials), 500 (server error)

### /api/logout
- **Request:** No body required (must be authenticated)
- **Success:** 200 OK `{ "message": "Logout successful" }`
- **Errors:** 401 (no valid session), 500 (server error)

## Session Management & Cookies
- On login, a `session_id` cookie is set (HttpOnly, Secure, SameSite=Strict, Path=/, Expires=24h).
- All protected endpoints require a valid `session_id` cookie.
- On logout, the session is deleted and the cookie is cleared.
- Session fixation and CSRF are mitigated by secure cookie settings and session regeneration.

## Security
- All user input is validated and sanitized to prevent SQL injection and XSS.
- Passwords are hashed before storage.
- See `docs/authentication.md` and test files for more details.

## Contributors

* [DavJesse](https://github.com/DavJesse)

* [Murzuqisah](https://github.com/Murzuqisah)

* [karodgers](https://github.com/karodgers)

* [siaka385](https://github.com/siaka385)

* [nyunja](https://github.com/nyunja)

* [oduortoni](https://github.com/oduortoni)

## License

This project is licensed under the MIT License.