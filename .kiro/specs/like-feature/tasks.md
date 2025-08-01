# Like/Dislike Feature Tasks

## Backend

- [ ] Create database migrations for `post_reactions` and `comment_reactions` tables.
- [ ] Add `likes_count` and `dislikes_count` to `posts` and `comments` tables.
- [ ] Implement `POST /api/posts/{id}/reaction` endpoint.
- [ ] Implement `DELETE /api/posts/{id}/reaction` endpoint.
- [ ] Implement `POST /api/comments/{id}/reaction` endpoint.
- [ ] Implement `DELETE /api/comments/{id}/reaction` endpoint.
- [ ] Implement WebSocket events for real-time updates.

## Frontend

- [ ] Create `ReactionButtons` component.
- [ ] Add `ReactionButtons` to `Post` component.
- [- ] Add `ReactionButtons` to `Comment` component.
- [ ] Implement real-time updates for like and dislike counts.