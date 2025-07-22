Read backend-file-structure.md and endpoint.md and entity_relation_diagram.mmd

Given the current posts implementation. let finish up what is left from this check list:

6. Posts and Comments
Tasks:

Verify Posts and Comments Database Structure:

Ensure the POSTS table migration (e.g., 000002_create_posts_table.up.sql, 000002_create_posts_table.down.sql in backend/pkg/db/migrations/sqlite) includes columns: id (PK, INTEGER), user_id (FK to USERS), content (TEXT), image (TEXT, nullable), privacy (TEXT, enum: public, almost_private, private), created_at (DATETIME), as per entity_relation_diagram.mmd.
Create a migration for the COMMENTS table (e.g., 000006_create_comments_table.up.sql, 000006_create_comments_table.down.sql) with columns: id (PK, INTEGER), post_id (FK to POSTS), user_id (FK to USERS), content (TEXT), image (TEXT, nullable), created_at (DATETIME).
Create a migration for the POST_VISIBILITY table (e.g., 000007_create_post_visibility_table.up.sql, 000007_create_post_visibility_table.down.sql) with columns: post_id (PK, FK to POSTS), viewer_id (PK, FK to USERS) for private post visibility.
Apply migrations using golang-migrate and verify table creation with foreign keys and indexes (e.g., on post_id, user_id) in SQLite.


Implement Post Creation Backend:

Create internal/api/handlers/post_handler.go with a handler for POST /api/posts (as per endpoint.md), requiring authentication via session cookie.
Implement business logic in internal/service/post_service.go to validate inputs: content (required), privacy (public, almost_private, private), image (optional), and viewer_ids (for private posts).
Handle image uploads (JPEG, PNG, GIF) in post_service.go, storing images in backend/uploads/posts/ and saving the file path in the image column.
In internal/store/post_store.go, insert post data into the POSTS table and, for private posts, insert allowed viewer IDs into the POST_VISIBILITY table.
Return the created post data (e.g., { "id": 2, "authorId": 123, "content": "My new post content.", "createdAt": "2025-06-29T13:00:00Z" }) or error (e.g., invalid file format).


Implement Comment Creation Backend:

In internal/api/handlers/post_handler.go, create a handler for POST /api/posts/{postId}/comments (as per endpoint.md), requiring authentication.
In internal/service/post_service.go, validate inputs: content (required), image (optional, JPEG, PNG, GIF).
Store images in backend/uploads/comments/ and save the file path in the image column of the COMMENTS table.
In internal/store/post_store.go, insert comment data into the COMMENTS table, linking to the specified post_id and user_id.
Return the created comment data or error (e.g., invalid post ID).


Implement Post and Comment Retrieval Backend:

Update internal/api/handlers/post_handler.go to handle GET /api/posts (renamed from GET /feed in endpoint.md for consistency with README.md) to fetch posts, filtering by privacy:
Public: visible to all.
Almost private: visible to followers (check FOLLOWERS table where is_accepted=true).
Private: visible to users in POST_VISIBILITY table.


In internal/service/post_service.go, implement logic to query POSTS and join with FOLLOWERS or POST_VISIBILITY for access control.
In internal/store/post_store.go, write SQL queries to fetch posts with user details (e.g., nickname, avatar from USERS).
Implement GET /api/posts/{postId}/comments in post_handler.go to fetch comments, respecting post privacy.
Ensure only authorized users (based on session and privacy settings) access posts/comments via post_service.go.


Implement Post Deletion Backend:

In internal/api/handlers/post_handler.go, create a handler for DELETE /api/posts/{postId} (as per endpoint.md), restricted to the post’s author.
In internal/service/post_service.go, validate that the authenticated user is the post’s user_id before deletion.
In internal/store/post_store.go, delete the post from POSTS and associated entries from POST_VISIBILITY and COMMENTS.
Return a success response (e.g., { "message": "Post deleted" }) or error (e.g., unauthorized).


Develop Post Creation Frontend:

In frontend/src/components/PostForm.js, create a post creation form with fields for content (textarea), privacy (dropdown: public, almost_private, private), image (file input), and viewer_ids (multi-select for private posts, fetched from GET /api/users/{userId}/following).
Implement client-side validation in Next.js (e.g., non-empty content, valid image formats: JPEG, PNG, GIF).
Send form data to POST /api/posts using fetch with credentials: 'include' (as per README.md), handling multipart/form-data for images.
Display success (e.g., “Post created”) or error messages (e.g., “Invalid image format”) based on API responses.
Style the form using Tailwind CSS for responsiveness and accessibility.


Develop Comment Creation Frontend:

In frontend/src/components/CommentForm.js, create a comment form component on post pages with fields for content (textarea) and image (file input).
Implement client-side validation for content and image formats (JPEG, PNG, GIF).
Send form data to POST /api/posts/{postId}/comments using fetch with credentials: 'include', handling responses and updating the UI.
Style the form consistently with the post creation form using Tailwind CSS.


Enhance Post Feed Component:

In frontend/src/components/PostFeed.js, update the post feed to fetch posts from GET /api/posts, displaying content, image, author details (from USERS), and created_at.
Add privacy indicators (e.g., icons for public, almost_private, private) using Tailwind CSS classes.
Include a comments section below each post, fetching comments from GET /api/posts/{postId}/comments and reusing CommentForm.js for adding comments.
Implement pagination or infinite scrolling for the post feed using Next.js dynamic routing or state management.
Ensure responsive design across devices using Tailwind CSS.


Implement Image Display Logic:

Create an API endpoint in internal/api/handlers/upload_handler.go (e.g., GET /api/uploads/:type/:filename) to serve images from backend/uploads/ (posts or comments).
In frontend/src/components/Post.js and Comment.js, display images using <img> tags with URLs from the backend (e.g., /api/uploads/posts/image.jpg).
Handle missing or invalid images with placeholders or hidden image areas.
Optimize image loading with lazy loading in Next.js (e.g., using <Image> component).


Test Posts and Comments Functionality:

Test post creation (POST /api/posts) with different privacy settings, verifying database entries in POSTS and POST_VISIBILITY, and correct visibility (public, almost_private, private).
Test comment creation (POST /api/posts/{postId}/comments), ensuring comments are linked in COMMENTS and respect post privacy.
Verify image uploads for posts/comments, checking that JPEG, PNG, and GIF files are stored in backend/uploads/ and paths are saved in image columns.
Test post deletion (DELETE /api/posts/{postId}), confirming removal from POSTS, POST_VISIBILITY, and COMMENTS.
Test post/comment retrieval (GET /api/posts, GET /api/posts/{postId}/comments), ensuring only authorized users see restricted content.
Validate frontend forms and feed components for functionality, responsiveness, and error handling in a Docker environment (docker-compose up).


Document Posts and Comments Features:

Update README.md with instructions for creating, viewing, and deleting posts/comments, including privacy settings and image support.
Document API endpoints (POST /api/posts, GET /api/posts, DELETE /api/posts/{postId}, POST /api/posts/{postId}/comments, GET /api/posts/{postId}/comments, GET /api/uploads/:type/:filename) in docs/api.md, detailing inputs, outputs, and access controls, as per endpoint.md.
Describe frontend components (PostForm.js, CommentForm.js, PostFeed.js) in README.md, including styling (Tailwind CSS) and responsiveness notes.
Note image handling setup (filesystem storage in backend/uploads/, supported formats: JPEG, PNG, GIF).