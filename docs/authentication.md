# Authentication System Documentation

## Overview
This document provides setup instructions, API endpoint details, and session/cookie management information for the authentication system in this project.

---

## 1. Setup Instructions

### Backend
- Ensure you have Go and SQLite installed.
- The backend uses an SQLite database for user and session storage.
- To run the backend:
  1. Clone the repository.
  2. Navigate to the `backend` directory.
  3. Run `go run server.go` (or use Docker Compose as described in the main README).

### Frontend
- The frontend is a Next.js app and communicates with the backend via HTTP.
- See the main README for environment variable setup and running instructions.

### Testing Authentication
- Unit and integration tests for authentication are located in `backend/internal/api/handlers/tests/`.
- To run tests:
  1. Navigate to the `backend` directory.
  2. Run `go test ./internal/api/handlers/tests/`

---

## 2. API Endpoints

### POST `/api/register`
- **Description:** Register a new user.
- **Request Body:** (JSON or multipart/form-data)
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
- **Success Response:**
  - Status: `200 OK`
  - Body:
    ```json
    { "message": "Registration successful" }
    ```
- **Error Codes:**
  - `400 Bad Request`: Invalid input or missing fields
  - `409 Conflict`: Email already exists
  - `500 Internal Server Error`: Server/database error

### POST `/api/login`
- **Description:** Authenticate a user and create a session.
- **Request Body:** (JSON or form data)
  ```json
  {
    "email": "user@example.com",
    "password": "password123"
  }
  ```
- **Success Response:**
  - Status: `200 OK`
  - Body:
    ```json
    { "message": "Login successful" }
    ```
  - **Session Cookie:** `session_id` (HttpOnly, Secure, SameSite=Strict)
- **Error Codes:**
  - `400 Bad Request`: Invalid input
  - `401 Unauthorized`: Invalid credentials
  - `500 Internal Server Error`: Server/database error

### POST `/api/logout`
- **Description:** Invalidate the user's current session.
- **Success Response:**
  - Status: `200 OK`
  - Body:
    ```json
    { "message": "Logout successful" }
    ```
- **Error Codes:**
  - `401 Unauthorized`: No valid session
  - `500 Internal Server Error`: Server/database error

---

## 3. Session Management & Cookie Configuration

- **Session Creation:**
  - On successful login, a new session is created in the database and a `session_id` cookie is set in the user's browser.
- **Session Cookie Properties:**
  - `HttpOnly`: Prevents JavaScript access to the cookie.
  - `Secure`: Only sent over HTTPS (set in production).
  - `SameSite=Strict`: Prevents CSRF by restricting cross-site requests.
  - `Path=/`
  - `Expires`: Set to session expiration time (e.g., 24 hours).
- **Session Validation:**
  - All protected endpoints require a valid `session_id` cookie.
  - The backend checks the session's validity and expiration before granting access.
- **Session Invalidation:**
  - On logout, the session is deleted from the database and the cookie is cleared.

---

## 4. Security Notes
- All user input is validated and sanitized to prevent SQL injection and XSS.
- Passwords are hashed using a secure algorithm before storage.
- Session fixation and CSRF are mitigated via secure cookie settings and session regeneration on login.

---

## 5. References
- See `backend/internal/api/handlers/tests/auth_handler_sql_injection_test.go` for security-related test cases.
- See the main `README.md` for project setup and Docker instructions.
