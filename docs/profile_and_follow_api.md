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
