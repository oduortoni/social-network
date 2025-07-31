# Design Document

## Overview

This design extends the existing comment system to support edit and delete operations, mirroring the functionality already implemented for posts. The solution maintains architectural consistency by following the same patterns used in post management while adding the necessary backend endpoints and frontend UI components.

The design leverages the existing three-tier architecture (handlers → services → stores) and reuses established UI patterns from the post management system. A new database migration will add an `updated_at` column to track comment edits, and new API endpoints will handle comment updates and deletions.

## Architecture

### Backend Architecture

The backend follows the existing layered architecture:

```
API Layer (handlers/post_handler.go)
├── UpdateComment(w http.ResponseWriter, r *http.Request)
├── DeleteComment(w http.ResponseWriter, r *http.Request)

Service Layer (service/post_service.go)
├── UpdateComment(commentID, userID int64, content string, imageData []byte, imageMimeType string) (*models.Comment, error)
├── DeleteComment(commentID, userID int64) error

Store Layer (store/post_store.go)
├── UpdateComment(commentID int64, content, imagePath string) (*models.Comment, error)
├── DeleteComment(commentID int64) error
├── GetCommentByID(commentID int64) (*models.Comment, error)
```

### Frontend Architecture

The frontend extends the existing comment system:

```
CommentList.jsx
├── Comment dropdown menu (three-dot icon)
├── Edit modal (similar to post edit modal)
├── Delete confirmation modal (similar to post delete modal)
├── State management for edit/delete operations

API Integration (lib/auth.js)
├── updateComment(commentId, content, image)
├── deleteComment(commentId)
```

## Components and Interfaces

### Database Schema Changes

A new migration will add an `updated_at` column to the Comments table:

```sql
-- 000014_add_updated_at_to_comments.up.sql
ALTER TABLE Comments ADD COLUMN updated_at DATETIME;

-- Update existing comments to have updated_at = created_at
UPDATE Comments SET updated_at = created_at WHERE updated_at IS NULL;
```

### Backend API Endpoints

**PUT /posts/{postId}/comments/{commentId}**
- Updates a specific comment
- Validates user ownership
- Handles multipart form data (content + optional image)
- Returns updated comment with author information

**DELETE /posts/{postId}/comments/{commentId}**
- Deletes a specific comment
- Validates user ownership
- Returns 204 No Content on success

### Backend Models

The Comment model will be extended to include edit tracking:

```go
type Comment struct {
    ID        int64     `json:"id"`
    PostID    int64     `json:"post_id"`
    UserID    int64     `json:"user_id"`
    Content   string    `json:"content"`
    Image     string    `json:"image,omitempty"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt *time.Time `json:"updated_at,omitempty"`
    IsEdited  bool      `json:"is_edited"`
    Author    User      `json:"author"`
}
```

### Frontend Components

**Enhanced CommentList Component**
- Adds three-dot dropdown menu to each comment
- Implements edit modal with form validation
- Implements delete confirmation modal
- Manages local state for UI interactions
- Shows "edited" indicator for modified comments

**UI State Management**
```javascript
const [openDropdown, setOpenDropdown] = useState(null);
const [editModal, setEditModal] = useState(null);
const [deleteConfirmation, setDeleteConfirmation] = useState(null);
const [editContent, setEditContent] = useState('');
const [editImage, setEditImage] = useState(null);
```

## Data Models

### Comment Data Flow

**Edit Operation:**
1. User clicks edit → Modal opens with current content
2. User modifies content → Form validation occurs
3. User saves → API call with FormData
4. Backend validates ownership → Updates database
5. Frontend receives updated comment → Updates local state
6. UI shows updated content with "edited" indicator

**Delete Operation:**
1. User clicks delete → Confirmation modal opens
2. User confirms → API call to delete endpoint
3. Backend validates ownership → Removes from database
4. Frontend receives success → Removes from local state
5. UI updates to hide deleted comment

### Data Validation

**Frontend Validation:**
- Comment content cannot be empty
- Image file type validation (PNG, JPEG, GIF)
- File size limits (20MB maximum)

**Backend Validation:**
- User ownership verification
- Content sanitization
- Image format validation
- Database constraint validation

## Error Handling

### Frontend Error Handling

**Edit Errors:**
- Network failures → Show error message, keep modal open
- Validation errors → Highlight invalid fields
- Unauthorized access → Show "You can only edit your own comments"

**Delete Errors:**
- Network failures → Show error toast/message
- Unauthorized access → Show "You can only delete your own comments"
- Comment not found → Show "Comment no longer exists"

### Backend Error Handling

**HTTP Status Codes:**
- 200: Successful update
- 204: Successful deletion
- 400: Invalid request data
- 401: Unauthorized (not logged in)
- 403: Forbidden (not comment owner)
- 404: Comment not found
- 500: Internal server error

**Error Response Format:**
```json
{
  "message": "You can only edit your own comments"
}
```

## Testing Strategy

### Backend Testing

**Unit Tests:**
- Service layer methods for update/delete operations
- Store layer database operations
- Handler input validation and error responses

**Integration Tests:**
- End-to-end API endpoint testing
- Database transaction testing
- Authentication and authorization testing

**Test Cases:**
- Successful comment update/delete
- Unauthorized access attempts
- Invalid input handling
- Database constraint violations
- Image upload validation

### Frontend Testing

**Component Tests:**
- Comment dropdown menu interactions
- Edit modal form validation
- Delete confirmation flow
- Error state handling

**Integration Tests:**
- API call success/failure scenarios
- State management during operations
- UI updates after operations

**User Experience Tests:**
- Keyboard navigation
- Mobile responsiveness
- Loading states
- Error message clarity

### Manual Testing Scenarios

1. **Edit Comment Flow:**
   - Open dropdown → Select edit → Modify content → Save changes
   - Verify "edited" indicator appears
   - Test with and without image updates

2. **Delete Comment Flow:**
   - Open dropdown → Select delete → Confirm deletion
   - Verify comment disappears from UI
   - Test cancellation flow

3. **Permission Testing:**
   - Attempt to edit/delete other users' comments
   - Verify appropriate error messages
   - Test with different user roles

4. **Edge Cases:**
   - Very long comment content
   - Special characters and emojis
   - Network interruptions during operations
   - Concurrent edit attempts