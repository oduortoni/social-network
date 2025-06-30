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

*   **POST /api/login**
    *   Description: Authenticates a user and returns a session token.
    *   Request Body:
        ```json
        {
          "email": "user@example.com",
          "password": "password123"
        }
        ```
    *   Response:
        ```json
        {
          "token": "your-session-token"
        }
        ```

*   **POST /api/register**
    *   Description: Registers a new user.
    *   Request Body:
        ```json
        {
          "username": "newuser",
          "email": "newuser@example.com",
          "password": "password123"
        }
        ```
    *   Response:
        ```json
        {
          "id": 123,
          "username": "newuser",
          "email": "newuser@example.com"
        }
        ```

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
git clone [https://github.com/oduortoni/social-network.git](https://github.com/oduortoni/social-network.git)
cd social-network
```

## Contributors

* [DavJesse](https://github.com/DavJesse)

* [Murzuqisah](https://github.com/Murzuqisah)

* [karodgers](https://github.com/karodgers)

* [siaka385](https://github.com/siaka385)

* [nyunja](https://github.com/nyunja)

* [oduortoni](https://github.com/oduortoni)

## License

This project is licensed under the MIT License.