import React, { useState, useEffect } from "react";
import { fetchComments } from "../../lib/auth";

const CommentList = ({ postId, newComment }) => {
  const [comments, setComments] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  // Format date helper
  const formatDate = (dateString) => {
    const date = new Date(dateString);
    const now = new Date();
    const diffInSeconds = Math.floor((now - date) / 1000);

    if (diffInSeconds < 60) {
      return 'Just now';
    } else if (diffInSeconds < 3600) {
      const minutes = Math.floor(diffInSeconds / 60);
      return `${minutes}m ago`;
    } else if (diffInSeconds < 86400) {
      const hours = Math.floor(diffInSeconds / 3600);
      return `${hours}h ago`;
    } else if (diffInSeconds < 604800) {
      const days = Math.floor(diffInSeconds / 86400);
      return `${days}d ago`;
    } else {
      return date.toLocaleDateString();
    }
  };

  // Get display name for user
  const getDisplayName = (author) => {
    if (author?.nickname) {
      return author.nickname;
    }
    const firstName = author?.first_name || '';
    const lastName = author?.last_name || '';
    const fullName = `${firstName} ${lastName}`.trim();
    return fullName || 'User';
  };

  // Load comments
  const loadComments = async () => {
    try {
      setLoading(true);
      const result = await fetchComments(postId);
      
      if (result.success) {
        setComments(result.data || []);
        setError("");
      } else {
        setError(result.error);
      }
    } catch (error) {
      console.error('Error loading comments:', error);
      setError('Failed to load comments');
    } finally {
      setLoading(false);
    }
  };

  // Load comments on mount and when postId changes
  useEffect(() => {
    if (postId) {
      loadComments();
    }
  }, [postId]);

  // Add new comment to the list when one is created
  useEffect(() => {
    if (newComment) {
      setComments(prev => [newComment, ...prev]); // Add to beginning for newest first
    }
  }, [newComment]);

  if (loading) {
    return (
      <div className="text-center py-4" style={{ color: 'var(--secondary-text)' }}>
        <div className="text-sm">Loading comments...</div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="text-center py-4">
        <div className="text-red-400 text-sm">{error}</div>
        <button 
          onClick={loadComments}
          className="text-blue-400 text-sm mt-2 hover:underline"
        >
          Try again
        </button>
      </div>
    );
  }

  if (comments.length === 0) {
    return (
      <div className="text-center py-4" style={{ color: 'var(--secondary-text)' }}>
        <div className="text-sm">No comments yet. Be the first to comment!</div>
      </div>
    );
  }
  console.log("comments...", comments);

  return (
    <div className="space-y-3">
      {comments.map((comment) => (
        <div key={comment.id} className="flex gap-3">
          {/* Comment Author Avatar */}
          <img
            src={comment.author?.avatar && comment.author.avatar !== "no profile photo" 
              ? `http://localhost:9000/avatar?avatar=${comment.author.avatar}` 
              : "http://localhost:9000/avatar?avatar=user-profile-circle-svgrepo-com.svg"
            }
            alt={getDisplayName(comment.author)}
            className="w-8 h-8 rounded-full object-cover flex-shrink-0"
          />

          {/* Comment Content */}
          <div className="flex-1">
            <div 
              className="rounded-lg p-3" 
              style={{ backgroundColor: 'var(--secondary-background)' }}
            >
              {/* Author Name */}
              <div className="font-medium text-white text-sm mb-1">
                {getDisplayName(comment.author)}
              </div>

              {/* Comment Text */}
              {comment.content && (
                <div className="text-white text-sm whitespace-pre-wrap mb-2">
                  {comment.content}
                </div>
              )}

              {/* Comment Image */}
              {comment.image && (
                <div className="mt-2">
                  <img 
                    src={`http://localhost:9000/avatar?avatar=${comment.image}`}
                    alt="Comment image" 
                    className="max-w-full max-h-48 rounded-lg object-cover"
                    onError={(e) => {
                      e.target.style.display = 'none';
                    }}
                  />
                </div>
              )}
            </div>

            {/* Comment Timestamp */}
            <div 
              className="text-xs mt-1 ml-3" 
              style={{ color: 'var(--secondary-text)' }}
            >
              {formatDate(comment.created_at)}
            </div>
          </div>
        </div>
      ))}
    </div>
  );
};

export default CommentList;
