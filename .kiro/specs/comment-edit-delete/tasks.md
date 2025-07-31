# Implementation Plan

- [x] 1. Create database migration for comment edit tracking

  - Create migration file to add `updated_at` column to Comments table
  - Write SQL to update existing comments with created_at as updated_at
  - Test migration up and down scripts
  - _Requirements: 1.6_

- [x] 2. Update Comment model to support edit tracking

  - Add `UpdatedAt` and `IsEdited` fields to Comment struct
  - Implement logic to determine if comment is edited based on timestamps
  - Update JSON tags for proper API response formatting
  - _Requirements: 1.6_

- [ ] 3. Extend backend store interface and implementation

  - [x] 3.1 Add comment management methods to PostStoreInterface

    - Add `UpdateComment(commentID int64, content, imagePath string) (*models.Comment, error)` method signature
    - Add `DeleteComment(commentID int64) error` method signature
    - Add `GetCommentByID(commentID int64) (*models.Comment, error)` method signature
    - _Requirements: 1.5, 2.3_

  - [x] 3.2 Implement comment management methods in PostStore
    - Write `UpdateComment` method with SQL UPDATE query and timestamp handling
    - Write `DeleteComment` method with SQL DELETE query
    - Write `GetCommentByID` method with SQL SELECT query including author join
    - Add proper error handling and transaction management
    - _Requirements: 1.5, 2.3_

- [ ] 4. Extend backend service interface and implementation

  - [x] 4.1 Add comment management methods to PostServiceInterface

    - Add `UpdateComment(commentID, userID int64, content string, imageData []byte, imageMimeType string) (*models.Comment, error)` method signature
    - Add `DeleteComment(commentID, userID int64) error` method signature
    - _Requirements: 1.5, 2.3, 5.3_

  - [ ] 4.2 Implement comment management methods in PostService
    - Write `UpdateComment` method with ownership validation and image handling
    - Write `DeleteComment` method with ownership validation
    - Implement proper error handling for unauthorized access and not found cases
    - Add image processing and storage logic similar to post updates
    - _Requirements: 1.5, 2.3, 5.3, 5.4_

- [ ] 5. Create backend API handlers for comment operations

  - [x] 5.1 Add UpdateComment handler method

    - Parse multipart form data for content and optional image
    - Extract commentID and postID from URL path parameters
    - Validate user authentication and extract userID from context
    - Call service layer UpdateComment method
    - Return updated comment data or appropriate error response
    - _Requirements: 1.1, 1.4, 1.5, 1.7_

  - [x] 5.2 Add DeleteComment handler method
    - Extract commentID and postID from URL path parameters
    - Validate user authentication and extract userID from context
    - Call service layer DeleteComment method
    - Return 204 No Content on success or appropriate error response
    - _Requirements: 2.1, 2.4, 2.5_

- [x] 6. Add API routes for comment operations

  - Register PUT `/posts/{postId}/comments/{commentId}` route with auth middleware
  - Register DELETE `/posts/{postId}/comments/{commentId}` route with auth middleware
  - Ensure routes are properly protected and follow existing patterns
  - _Requirements: 1.1, 2.1_

- [x] 7. Create frontend API functions for comment operations

  - [x] 7.1 Add updateComment function to auth.js

    - Create function that accepts commentId, content, and optional image
    - Build FormData object and make PUT request to comment endpoint
    - Handle response parsing and error cases
    - Return success/error object consistent with existing API functions
    - _Requirements: 1.4, 1.5, 1.7_

  - [x] 7.2 Add deleteComment function to auth.js
    - Create function that accepts commentId
    - Make DELETE request to comment endpoint
    - Handle response and error cases
    - Return success/error object consistent with existing API functions
    - _Requirements: 2.3, 2.5_

- [x] 8. Enhance CommentList component with dropdown menu

  - [x] 8.1 Add three-dot menu icon to each comment

    - Import MoreHorizontalIcon from lucide-react
    - Add menu button positioned next to comment timestamp
    - Implement dropdown state management (openDropdown)
    - Style menu button using global CSS color variables
    - _Requirements: 1.1, 2.1, 3.1, 4.1, 4.4_

  - [x] 8.2 Implement dropdown menu with conditional options
    - Create dropdown container with Edit, Delete, and Follow options
    - Show Edit/Delete options only for user's own comments
    - Show Follow option for all comments
    - Style dropdown using global CSS color variables and hover effects
    - Add click outside handler to close dropdown
    - _Requirements: 1.2, 2.1, 3.2, 5.1, 5.2_

- [x] 9. Implement comment edit functionality

  - [x] 9.1 Create edit modal component

    - Build modal similar to post edit modal with form fields
    - Pre-fill textarea with current comment content
    - Add image upload input for optional image updates
    - Style modal using global CSS color variables
    - _Requirements: 1.3, 4.2, 4.4_

  - [x] 9.2 Add edit form validation and submission

    - Validate that comment content is not empty
    - Handle form submission with updateComment API call
    - Show loading state during submission
    - Handle success and error responses appropriately
    - Update local comment state on successful edit
    - _Requirements: 1.4, 1.5, 1.7_

  - [x] 9.3 Display edited indicator for modified comments
    - Add "edited" text next to timestamp for comments with updated_at different from created_at
    - Style indicator using secondary text color variable
    - Ensure indicator appears immediately after successful edit
    - _Requirements: 1.6_

- [x] 10. Implement comment delete functionality

  - [x] 10.1 Create delete confirmation modal

    - Build confirmation modal similar to post delete modal
    - Display warning message about permanent deletion
    - Add Cancel and Delete buttons with appropriate styling
    - Use warning color variable for delete button
    - _Requirements: 2.2, 4.2, 4.4_

  - [x] 10.2 Handle delete confirmation and execution
    - Call deleteComment API function on confirmation
    - Remove comment from local state on successful deletion
    - Handle error cases with appropriate user feedback
    - Close modal and dropdown after successful deletion
    - _Requirements: 2.3, 2.4, 2.5_

- [x] 12. Write comprehensive tests for comment operations

  - [x] 12.1 Create backend unit tests

    - Write tests for store layer comment CRUD operations
    - Write tests for service layer validation and business logic
    - Write tests for handler layer request/response handling
    - Test authorization and ownership validation
    - _Requirements: 5.4_

  - [x] 12.2 Create frontend component tests
    - Test comment dropdown menu interactions
    - Test edit modal form validation and submission
    - Test delete confirmation flow
    - Test error handling and loading states
    - _Requirements: 1.7, 2.5_

- [ ] 13. Integration testing and bug fixes
  - Test complete edit and delete workflows end-to-end
  - Verify UI consistency with existing post management
  - Test edge cases like network failures and concurrent operations
  - Fix any issues discovered during testing
  - _Requirements: 4.1, 4.2, 4.3, 4.4, 4.5_
