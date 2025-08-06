# Profile and Follow API (Backend)

This document explains the backend endpoints for user profile and follow functionality in the social network application, including the structure of the profile response and button status logic.

## Profile Endpoints

### 1. Get User Profile
- **Endpoint:** `GET /profile/{userid}`
- **Description:** Returns the profile information for a user, including name, bio, avatar, posts, and follow/message button status.
- **Authentication:** Required
- **Response Example:**
```json
{
  "profile_details": {
    "firstname": "John",
    "lastname": "Doe",
    "email": "john@example.com",
    "avatar": "/avatar/1.jpg",
    "id": 1,
    "about": "Web developer. Coffee enthusiast.",
    "nickname": "johnny",
    "followbtnstatus": "follow", // see below for explanation
    "messagebtnstatus": "visible", // see below for explanation
    "dateofbirth": "1990-01-01",
    "profile": true,
    "numberoffollowers": 123,
    "numberoffollowees": 87,
    "numberofposts": 5
  },
  "posts": [
    // Array of user's posts (see Post structure)
  ]
}
```

#### ProfileDetails Fields
- **firstname, lastname, email, avatar, id, about, nickname, dateofbirth, profile, numberoffollowers, numberofposts, numberoffollowees:** Basic user info and stats.
- **followbtnstatus:** Controls the visibility and state of the follow button:
  - `"hide"`: Button is hidden (for the logged-in user's own profile)
  - `"pending"`: Follow request sent but not yet accepted
  - `"follow"`: Not following, follow option available
  - `"following"`: Already following the user
- **messagebtnstatus:** Controls the visibility of the message button:
  - `"hide"`: Button is hidden (for the logged-in user's own profile)
  - `"visible"`: Visible to followers (messaging allowed)

#### Example: Button Status Logic
- If you view your own profile: `followbtnstatus = "hide"`, `messagebtnstatus = "hide"`
- If you view another user you do not follow: `followbtnstatus = "follow"`, `messagebtnstatus = "hide"`
- If you sent a follow request: `followbtnstatus = "pending"`, `messagebtnstatus = "hide"`
- If you follow the user: `followbtnstatus = "following"`, `messagebtnstatus = "visible"`

### 2. Get User Followers
- **Endpoint:** `GET /profile/{userid}/followers`
- **Description:** Returns a list of users who follow the specified user.
- **Authentication:** Required
- **Response Example:**
  ```json
  {
    "user": [
      {
        "firstname": "Alice",
        "lastname": "Smith",
        "avatar": "/avatar/2.jpg",
        "follower_id": 2
      },
      {
        "firstname": "Bob",
        "lastname": "Johnson",
        "avatar": "/avatar/3.jpg",
        "follower_id": 3
      }
    ]
  }
  ```

### 3. Get User Followees (Following)
- **Endpoint:** `GET /profile/{userid}/followees`
- **Description:** Returns a list of users that the specified user is following.
- **Authentication:** Required
- **Response Example:**
  ```json
  {
    "user": [
      {
        "firstname": "Charlie",
        "lastname": "Brown",
        "avatar": "/avatar/4.jpg",
        "follower_id": 4
      },
      {
        "firstname": "Diana",
        "lastname": "Prince",
        "avatar": "/avatar/5.jpg",
        "follower_id": 5
      }
    ]
  }
  ```

### 4. Edit User Profile
- **Endpoint:** `PUT /EditProfile`
- **Description:** Updates the authenticated user's profile information including personal details, avatar, and privacy settings.
- **Authentication:** Required
- **Content-Type:** `multipart/form-data`
- **Request Parameters:**
  - `email` (string, required): User's email address
  - `firstName` (string, required): User's first name
  - `lastName` (string, required): User's last name
  - `dob` (string, required): Date of birth in YYYY-MM-DD format
  - `nickname` (string, required): User's nickname/username
  - `aboutMe` (string, optional): User's bio/description
  - `profileVisibility` (string, required): Either "public" or "private"
  - `avatar` (file, optional): Profile picture image file

- **Success Response (200 OK):**
  ```json
  {
    "message": "Profile updated successfully"
  }
  ```

- **Error Responses:**
  - **400 Bad Request:** Invalid email format or failed to parse form
    ```json
    {
      "message": "Invalid email format"
    }
    ```
  - **401 Unauthorized:** User not authenticated
    ```json
    {
      "message": "User not found in context"
    }
    ```
  - **409 Conflict:** Email already exists for another user
    ```json
    {
      "message": "Email already exists"
    }
    ```
  - **500 Internal Server Error:** Server error during profile update
    ```json
    {
      "message": "Error message details"
    }
    ```

- **Example Request (using curl):**
  ```bash
  curl -X PUT http://localhost:8080/EditProfile \
    -H "Cookie: session_id=your_session_id" \
    -F "email=newemail@example.com" \
    -F "firstName=John" \
    -F "lastName=Doe" \
    -F "dob=1990-01-01" \
    -F "nickname=johndoe" \
    -F "aboutMe=Updated bio" \
    -F "profileVisibility=public" \
    -F "avatar=@/path/to/image.jpg"
  ```

- **Security Features:**
  - Input sanitization to prevent XSS attacks
  - Email validation using regex
  - Duplicate email checking (excluding current user)
  - File upload validation for avatar images
  - Authentication required via session middleware

- **Notes:**
  - All text fields are HTML-escaped to prevent XSS attacks
  - Profile visibility can be set to "public" (visible to all) or "private" (visible to followers only)
  - Avatar upload is optional; if not provided, existing avatar is preserved
  - Users can update their email to the same email they currently have
  - Date of birth should be in YYYY-MM-DD format

## Follow Endpoints

### 1. Follow a User
- **Endpoint:** `POST /follow`
- **Description:** Authenticated user sends a follow request to another user.
- **Request Body:**
  ```json
  {
    "followeeid": 2
  }
  ```

### 2. Unfollow a User
- **Endpoint:** `DELETE /unfollow`
- **Description:** Authenticated user unfollows another user.
- **Request Body:**
  ```json
  {
    "followeeid": 2
  }
  ```

### 3. Respond to Follow Request
- **Endpoint:** `POST /follow-request/{requestId}/request`
- **Description:** Accept or reject a follow request (for private accounts).
- **Response Example:**
  ```json
  {
    "status": "accepted"
  }
  ```

### 4. Cancel Follow Request
- **Endpoint:** `DELETE /follow-request/{requestId}/cancel`
- **Description:** Cancel a pending follow request.
- **Response Example:**
  ```json
  {
    "status": "cancelled"
  }
  ```

---

These endpoints allow clients to fetch user profile data, manage followers/following, and handle follow requests. All endpoints require authentication and return appropriate error messages for unauthorized or invalid actions.
