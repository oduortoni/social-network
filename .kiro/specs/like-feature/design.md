# Like/Dislike Feature Design

## Database Schema

We will add two new tables to the database:

- `post_reactions`
  - `user_id` (foreign key to `users` table)
  - `post_id` (foreign key to `posts` table)
  - `reaction_type` (e.g., 'like' or 'dislike')
  - `created_at` (timestamp)

- `comment_reactions`
  - `user_id` (foreign key to `users` table)
  - `comment_id` (foreign key to `comments` table)
  - `reaction_type` (e.g., 'like' or 'dislike')
  - `created_at` (timestamp)

We will also add `likes_count` and `dislikes_count` columns to the `posts` and `comments` tables.

## Backend Design

### API Endpoints

- `POST /api/posts/{id}/reaction`: Add or update a reaction for a post. The request body will contain the reaction type (e.g., `{"reaction": "like"}` or `{"reaction": "dislike"}`).
- `DELETE /api/posts/{id}/reaction`: Remove a reaction from a post.
- `POST /api/comments/{id}/reaction`: Add or update a reaction for a comment.
- `DELETE /api/comments/{id}/reaction`: Remove a reaction from a comment.

### Real-time Updates

We will use WebSockets to push real-time updates to clients.

- When a user reacts to a post, the server will broadcast a `post_reaction_updated` event with the new `likes_count` and `dislikes_count` to all clients viewing that post.
- When a user reacts to a comment, the server will broadcast a `comment_reaction_updated` event with the new `likes_count` and `dislikes_count` to all clients viewing that post.

## Frontend Design

### Components

- `ReactionButtons`: A reusable component that displays like and dislike buttons and handles the logic for adding, updating, and removing reactions.

### State Management

- The `ReactionButtons` component will manage the state of the user's reaction for a given item.
- The `Post` and `Comment` components will receive the `likes_count` and `dislikes_count` from the API and update them in real-time based on WebSocket events.