import React, { useState, useEffect } from 'react';
import { MoreHorizontalIcon, ThumbsUpIcon, ThumbsDownIcon, MessageCircleIcon, Globe, Users, Lock, Edit, Trash2, UserPlus } from 'lucide-react';
import { fetchPosts, deletePost, updatePost } from '../../lib/auth';
import VerifiedBadge from '../homepage/VerifiedBadge';
import CommentForm from './CommentForm';
import CommentList from './CommentList';
import ReactionButtons from './ReactionButtons';

const PostList = ({ refreshTrigger, user, posts: initialPosts }) => {
  const [posts, setPosts] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [expandedComments, setExpandedComments] = useState(new Set());
  const [newComments, setNewComments] = useState({});
  const [openDropdown, setOpenDropdown] = useState(null);
  const [deleteConfirmation, setDeleteConfirmation] = useState(null);
  const [editModal, setEditModal] = useState(null);
  const [editContent, setEditContent] = useState('');
  const [editImage, setEditImage] = useState(null);
  const [editLoading, setEditLoading] = useState(false);

  const loadPosts = async () => {
    setLoading(true);
    setError('');

    try {
      const result = await fetchPosts();

      if (result.success) {
        setPosts(result.data || []);
      } else {
        setError(result.error);
      }
    } catch (error) {
      setError('Failed to load posts');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    if (initialPosts) {
      setPosts(initialPosts);
      setLoading(false);
    } else {
      loadPosts();
    }
  }, [refreshTrigger, initialPosts]);

  const getPrivacyIcon = (privacy) => {
    switch (privacy) {
      case 'public':
        return Globe;
      case 'almost_private':
        return Users;
      case 'private':
        return Lock;
      default:
        return Globe;
    }
  };

  const getPrivacyLabel = (privacy) => {
    switch (privacy) {
      case 'public':
        return 'Public';
      case 'almost_private':
        return 'Followers';
      case 'private':
        return 'Private';
      default:
        return 'Public';
    }
  };

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
    } else {
      const days = Math.floor(diffInSeconds / 86400);
      return `${days}d ago`;
    }
  };

  // Toggle comments visibility for a post
  const toggleComments = (postId) => {
    setExpandedComments(prev => {
      const newSet = new Set(prev);
      if (newSet.has(postId)) {
        newSet.delete(postId);
      } else {
        newSet.add(postId);
      }
      return newSet;
    });
  };

  // Handle new comment creation
  const handleCommentCreated = (postId, comment) => {
    setNewComments(prev => ({
      ...prev,
      [postId]: comment
    }));

    // Clear the new comment after a short delay to allow CommentList to process it
    setTimeout(() => {
      setNewComments(prev => {
        const updated = { ...prev };
        delete updated[postId];
        return updated;
      });
    }, 100);
  };

  // Handle dropdown toggle
  const toggleDropdown = (postId) => {
    setOpenDropdown(openDropdown === postId ? null : postId);
  };

  // Handle edit post
  const handleEditPost = (post) => {
    setEditModal(post.id);
    setEditContent(post.content);
    setEditImage(null);
    setOpenDropdown(null);
  };

  // Handle edit post submission
  const handleEditSubmit = async (postId) => {
    if (!editContent.trim()) return;

    setEditLoading(true);
    try {
      const result = await updatePost(postId, editContent, editImage);
      if (result.success) {
        // Update the post in the local state
        setPosts(prev => prev.map(post =>
          post.id === postId ? result.data : post
        ));
        setEditModal(null);
        setEditContent('');
        setEditImage(null);
      } else {
        console.error('Failed to update post:', result.error);
        // You could add a toast notification here
      }
    } catch (error) {
      console.error('Error updating post:', error);
    } finally {
      setEditLoading(false);
    }
  };

  // Handle delete post
  const handleDeletePost = async (postId) => {
    try {
      const result = await deletePost(postId);
      if (result.success) {
        // Remove the post from the local state
        setPosts(prev => prev.filter(post => post.id !== postId));
        setDeleteConfirmation(null);
        setOpenDropdown(null);
      } else {
        console.error('Failed to delete post:', result.error);
        // You could add a toast notification here
      }
    } catch (error) {
      console.error('Error deleting post:', error);
    }
  };

  // Close dropdown when clicking outside
  useEffect(() => {
    const handleClickOutside = (event) => {
      // Check if the clicked element is part of any dropdown
      const isDropdownClick = event.target.closest('.dropdown-container');
      if (!isDropdownClick && openDropdown !== null) {
        setOpenDropdown(null);
      }
    };

    document.addEventListener('mousedown', handleClickOutside);
    return () => {
      document.removeEventListener('mousedown', handleClickOutside);
    };
  }, [openDropdown]);

  if (loading) {
    return (
      <div className="flex justify-center items-center py-8">
        <div style={{ color: 'var(--secondary-text)' }}>Loading posts...</div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="rounded-xl p-4 mb-4" style={{ backgroundColor: 'rgba(var(--danger-color-rgb), 0.2)', border: '1px solid var(--warning-color)' }}>
        <div style={{ color: 'var(--warning-color)' }} className="text-center">{error}</div>
        <button
          onClick={loadPosts}
          className="mt-2 w-full py-2 px-4 rounded-lg transition-colors"
          style={{ backgroundColor: 'var(--warning-color)', color: 'var(--primary-text)' }}
          onMouseOver={(e) => e.currentTarget.style.opacity = '0.8'}
          onMouseOut={(e) => e.currentTarget.style.opacity = '1'}
        >
          Try Again
        </button>
      </div>
    );
  }

  if (posts.length === 0) {
    return (
      <div className="rounded-xl p-8 text-center" style={{ backgroundColor: 'var(--primary-background)' }}>
        <div style={{ color: 'var(--secondary-text)' }} className="mb-2">No posts yet</div>
        <div className="text-sm" style={{ color: 'var(--secondary-text)' }}>Be the first to share something!</div>
      </div>
    );
  }

  return (
    <div className="space-y-4">
      {posts.map((post) => (
        <div key={post.id} className="rounded-xl p-4" style={{ backgroundColor: 'var(--primary-background)' }}>
          {/* Post Header */}
          <div className="flex justify-between mb-3">
            <div className="flex items-center gap-2">
              <div className="relative">
                <img
                  src={post.author?.avatar && post.author.avatar !== "no profile photo" ? `http://localhost:9000/avatar?avatar=${post.author.avatar}` : "http://localhost:9000/avatar?avatar=user-profile-circle-svgrepo-com.svg"}
                  alt={post.author?.nickname || `${post.author?.first_name || ''} ${post.author?.last_name || ''}`.trim() || 'User'}
                  className="w-10 h-10 rounded-full"
                />
                <div className="absolute -bottom-1 -right-1">
                  <VerifiedBadge />
                </div>
              </div>
              <div>
                <div className="flex items-center gap-2">
                  <span className="font-medium break-words" style={{ color: 'var(--primary-text)' }}>
                    {post.author?.nickname || `${post.author?.first_name || ''} ${post.author?.last_name || ''}`.trim() || 'User'}
                  </span>
                  <div className="flex items-center gap-1 text-xs" style={{ color: 'var(--secondary-text)' }}>
                    {React.createElement(getPrivacyIcon(post.privacy), { className: "w-3 h-3" })}
                    <span>{getPrivacyLabel(post.privacy)}</span>
                  </div>
                </div>
                <div className="text-xs" style={{ color: 'var(--secondary-text)' }}>
                  {formatDate(post.created_at)}
                  {post.is_edited && (
                    <span className="ml-2 text-xs" style={{ color: 'var(--secondary-text)' }}>
                      • edited
                    </span>
                  )}
                </div>
              </div>
            </div>
            {/* Post Actions Dropdown */}
            <div className="relative dropdown-container">
              <button
                onClick={() => toggleDropdown(post.id)}
                className="transition-colors"
                style={{ color: 'var(--secondary-text)' }}
                onMouseOver={(e) => e.currentTarget.style.color = 'var(--primary-text)'}
                onMouseOut={(e) => e.currentTarget.style.color = 'var(--secondary-text)'}
              >
                <MoreHorizontalIcon className="w-5 h-5" />
              </button>

              {/* Dropdown Menu */}
              {openDropdown === post.id && (
                <div
                  className="absolute right-0 top-8 w-48 rounded-lg shadow-lg z-10"
                  style={{ backgroundColor: 'var(--primary-background)', border: '1px solid var(--border-color)' }}
                >
                  <div className="py-1">
                    {user && user.id === post.user_id ? (
                      <>
                        {/* Edit Option */}
                        <button
                          onClick={(e) => {
                            e.stopPropagation();
                            handleEditPost(post);
                          }}
                          className="w-full px-4 py-2 text-left text-sm flex items-center gap-2 transition-colors"
                          style={{ color: 'var(--primary-text)' }}
                          onMouseEnter={(e) => e.currentTarget.style.backgroundColor = 'var(--hover-background)'}
                          onMouseLeave={(e) => e.currentTarget.style.backgroundColor = 'transparent'}
                        >
                          <Edit className="w-4 h-4" />
                          Edit
                        </button>

                        {/* Delete Option */}
                        <button
                          onClick={(e) => {
                            e.stopPropagation();
                            setDeleteConfirmation(post.id);
                            setOpenDropdown(null); // Close dropdown when opening modal
                          }}
                          className="w-full px-4 py-2 text-left text-sm flex items-center gap-2 transition-colors"
                          style={{ color: 'var(--warning-color)' }}
                          onMouseEnter={(e) => e.currentTarget.style.backgroundColor = 'var(--hover-background)'}
                          onMouseLeave={(e) => e.currentTarget.style.backgroundColor = 'transparent'}
                        >
                          <Trash2 className="w-4 h-4" />
                          Delete
                        </button>
                      </>
                    ) : (
                      <>
                        {/* Follow Option */}
                        <button
                          onClick={() => handleFollowUser(post.user_id)}
                          className="w-full px-4 py-2 text-left text-sm flex items-center gap-2 transition-colors"
                          style={{ color: 'var(--primary-text)' }}
                          onMouseEnter={(e) => e.currentTarget.style.backgroundColor = 'var(--hover-background)'}
                          onMouseLeave={(e) => e.currentTarget.style.backgroundColor = 'transparent'}
                        >
                          <UserPlus className="w-4 h-4" />
                          Follow
                        </button>
                      </>
                    )}
                  </div>
                </div>
              )}
            </div>
          </div>

          {/* Post Content */}
          <div className="mb-4">
            <p className="whitespace-pre-wrap break-words overflow-wrap-anywhere" style={{ color: 'var(--primary-text)' }}>{post.content}</p>

            {/* Post Image */}
            {post.image && (
              <div className="mt-3">
                <img
                  src={`http://localhost:9000/avatar?avatar=${post.image}`}
                  alt="Post image"
                  className="max-w-full rounded-lg"
                  onError={(e) => {
                    e.target.style.display = 'none';
                  }}
                />
              </div>
            )}
          </div>

          {/* Post Actions */}
          <div className="flex items-center justify-between pt-3 border-t" style={{ borderColor: 'var(--border-color)' }}>
            <div className="flex items-center gap-4">
              <ReactionButtons post={post} user={user} />
              <button
                onClick={() => toggleComments(post.id)}
                className="flex items-center gap-2 text-sm py-1.5 px-3 rounded-lg transition-colors"
                style={{ color: 'var(--secondary-text)' }}
                onMouseOver={(e) => e.currentTarget.style.backgroundColor = 'var(--hover-background)'}
                onMouseOut={(e) => e.currentTarget.style.backgroundColor = 'transparent'}>
                <MessageCircleIcon className="w-4 h-4" />
                <span>{expandedComments.has(post.id) ? 'Hide Comments' : 'Comment'}</span>
              </button>
            </div>
          </div>

          {/* Comments Section */}
          {expandedComments.has(post.id) && (
            <div className="mt-4 pt-3 border-t space-y-4" style={{ borderColor: 'var(--border-color)' }}>
              {/* Comment Form */}
              <CommentForm
                postId={post.id}
                user={user}
                onCommentCreated={(comment) => handleCommentCreated(post.id, comment)}
              />

              {/* Comments List */}
              <CommentList
                postId={post.id}
                newComment={newComments[post.id]}
                user={user}
              />
            </div>
          )}
        </div>
      ))}

      {/* Delete Confirmation Modal */}
      {deleteConfirmation && (
        <div className="fixed inset-0 flex items-center justify-center" style={{ zIndex: 9999, backgroundColor: 'rgba(0, 0, 0, 0.3)' }}>
          <div
            className="rounded-lg p-6 max-w-md w-full mx-4"
            style={{ backgroundColor: 'var(--primary-background)', border: '1px solid var(--border-color)' }}
          >
            <h3 className="text-lg font-semibold mb-4" style={{ color: 'var(--primary-text)' }}>Delete Post</h3>
            <p className="mb-6" style={{ color: 'var(--primary-text)' }}>Are you sure you want to delete this post? This action cannot be undone.</p>

            <div className="flex gap-3 justify-end">
              <button
                onClick={() => setDeleteConfirmation(null)}
                className="px-4 py-2 rounded-lg transition-colors"
                style={{ backgroundColor: 'var(--secondary-background)', color: 'var(--primary-text)' }}
                onMouseEnter={(e) => e.currentTarget.style.backgroundColor = 'var(--hover-background)'}
                onMouseLeave={(e) => e.currentTarget.style.backgroundColor = 'var(--secondary-background)'}
              >
                Cancel
              </button>
              <button
                onClick={() => handleDeletePost(deleteConfirmation)}
                className="px-4 py-2 rounded-lg transition-colors"
                style={{ backgroundColor: 'var(--warning-color)', color: 'var(--primary-text)' }}
                onMouseEnter={(e) => e.currentTarget.style.opacity = '0.8'}
                onMouseLeave={(e) => e.currentTarget.style.opacity = '1'}
              >
                Delete
              </button>
            </div>
          </div>
        </div>
      )}

      {/* Edit Post Modal */}
      {editModal && (
        <div className="fixed inset-0 flex items-center justify-center" style={{ zIndex: 9999, backgroundColor: 'rgba(0, 0, 0, 0.3)' }}>
          <div
            className="rounded-lg p-6 max-w-md w-full mx-4"
            style={{ backgroundColor: 'var(--primary-background)', border: '1px solid var(--border-color)' }}
          >
            <h3 className="text-lg font-semibold mb-4" style={{ color: 'var(--primary-text)' }}>Edit Post</h3>

            <div className="mb-4">
              <textarea
                value={editContent}
                onChange={(e) => setEditContent(e.target.value)}
                placeholder="What's on your mind?"
                className="w-full p-3 rounded-lg resize-none break-words overflow-wrap-anywhere"
                style={{
                  backgroundColor: 'var(--secondary-background)',
                  border: '1px solid var(--border-color)',
                  minHeight: '100px',
                  color: 'var(--primary-text)'
                }}
                rows={4}
              />
            </div>

            <div className="mb-4">
              <input
                type="file"
                accept="image/*"
                onChange={(e) => setEditImage(e.target.files[0])}
                className="w-full"
                style={{ backgroundColor: 'var(--secondary-background)', color: 'var(--primary-text)' }}
              />
            </div>

            <div className="flex gap-3 justify-end">
              <button
                onClick={() => {
                  setEditModal(null);
                  setEditContent('');
                  setEditImage(null);
                }}
                className="px-4 py-2 rounded-lg transition-colors"
                style={{ backgroundColor: 'var(--secondary-background)', color: 'var(--primary-text)' }}
                onMouseEnter={(e) => e.currentTarget.style.backgroundColor = 'var(--hover-background)'}
                onMouseLeave={(e) => e.currentTarget.style.backgroundColor = 'var(--secondary-background)'}
                disabled={editLoading}
              >
                Cancel
              </button>
              <button
                onClick={() => handleEditSubmit(editModal)}
                className="px-4 py-2 rounded-lg transition-colors"
                style={{ backgroundColor: 'var(--primary-accent)', color: 'var(--quinary-text)' }}
                onMouseEnter={(e) => e.currentTarget.style.opacity = '0.8'}
                onMouseLeave={(e) => e.currentTarget.style.opacity = '1'}
                disabled={editLoading || !editContent.trim()}
              >
                {editLoading ? 'Saving...' : 'Save Changes'}
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default PostList;
