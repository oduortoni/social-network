---

## Group Functionalities Documentation

This document outlines the implemented group functionalities, their purpose, and how to interact with them via the API.

### 1. Group Creation

Allows users to create new groups with a specified privacy setting.

*   **Endpoint:** `POST /groups`
*   **Request Body:**
    ```json
    {
      "title": "My Awesome Group",
      "description": "A group for awesome people.",
      "privacy": "public" // or "private"
    }
    ```
*   **Response (201 Created):**
    ```json
    {
      "id": 1,
      "title": "My Awesome Group",
      "description": "A group for awesome people.",
      "creator_id": 123, // ID of the authenticated user
      "privacy": "public",
      "created_at": "2025-08-10T12:00:00Z"
    }
    ```
*   **Notes:**
    *   `creator_id` is automatically set from the authenticated user's context.
    *   If `privacy` is not provided, it defaults to `public`.

### 2. Group Join Requests

Enables users to request to join public groups, and allows group creators/admins to manage these requests.

#### 2.1 Send Join Request

*   **Endpoint:** `POST /groups/{groupID}/join-request`
*   **Path Parameters:**
    *   `groupID`: The ID of the group to join.
*   **Request Body:** (Empty)
*   **Response (201 Created):**
    ```json
    {
      "id": 1,
      "group_id": 1,
      "user_id": 101, // ID of the requesting user
      "status": "pending",
      "created_at": "2025-08-10T12:05:00Z"
    }
    ```
*   **Notes:**
    *   Only applicable for public groups. Attempts to send requests to private groups will result in an error.

#### 2.2 Approve Join Request

*   **Endpoint:** `PUT /groups/{groupID}/join-request/{requestID}/approve`
*   **Path Parameters:**
    *   `groupID`: The ID of the group.
    *   `requestID`: The ID of the join request to approve.
*   **Request Body:** (Empty)
*   **Response (200 OK):**
    ```json
    {
      "message": "Join request approved"
    }
    ```
*   **Notes:**
    *   Requires authentication as a group creator/admin (authorization logic to be fully implemented).
    *   Upon approval, the requesting user is added as a member to the group.

#### 2.3 Reject Join Request

*   **Endpoint:** `PUT /groups/{groupID}/join-request/{requestID}/reject`
*   **Path Parameters:**
    *   `groupID`: The ID of the group.
    *   `requestID`: The ID of the join request to reject.
*   **Request Body:** (Empty)
*   **Response (200 OK):**
    ```json
    {
      "message": "Join request rejected"
    }
    ```
*   **Notes:**
    *   Requires authentication as a group creator/admin (authorization logic to be fully implemented).

### 3. Private Group Chat

Enables members of a group to send and retrieve messages within a private chat accessible only to them.

#### 3.1 Send Group Chat Message

*   **Endpoint:** `POST /groups/{groupID}/chat`
*   **Path Parameters:**
    *   `groupID`: The ID of the group to send a message to.
*   **Request Body:**
    ```json
    {
      "content": "Hello everyone in the group!"
    }
    ```
*   **Response (201 Created):**
    ```json
    {
      "id": 1,
      "group_id": 1,
      "sender_id": 101, // ID of the authenticated user
      "content": "Hello everyone in the group!",
      "created_at": "2025-08-10T12:10:00Z"
    }
    ```
*   **Notes:**
    *   Only group members can send messages. Attempts by non-members will result in an error.

#### 3.2 Get Group Chat Messages

*   **Endpoint:** `GET /groups/{groupID}/chat`
*   **Path Parameters:**
    *   `groupID`: The ID of the group to retrieve messages from.
*   **Query Parameters:**
    *   `limit` (optional): Maximum number of messages to retrieve (default: 10).
    *   `offset` (optional): Number of messages to skip (for pagination, default: 0).
*   **Response (200 OK):**
    ```json
    [
      {
        "id": 2,
        "group_id": 1,
        "sender_id": 102,
        "content": "Nice to meet you all!",
        "created_at": "2025-08-10T12:11:00Z"
      },
      {
        "id": 1,
        "group_id": 1,
        "sender_id": 101,
        "content": "Hello everyone in the group!",
        "created_at": "2025-08-10T12:10:00Z"
      }
    ]
    ```
*   **Notes:**
    *   Only group members can retrieve messages. Attempts by non-members will result in an error.
    *   Messages are returned in reverse chronological order (newest first).

### 4. Group Member Management

This functionality provides the underlying mechanism for tracking group memberships.

*   **Purpose:** To determine if a user is part of a specific group, which is crucial for authorization in private groups and chat.
*   **Details:** The `group_members` table stores associations between users and groups, including their roles (e.g., `member`, `admin`). This is primarily managed internally by the system (e.g., during join request approval or direct admin additions).
