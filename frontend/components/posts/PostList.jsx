import React, { useState, useEffect } from 'react';
import { MoreHorizontalIcon, ThumbsUpIcon, ThumbsDownIcon, MessageCircleIcon, Globe, Users, Lock, Edit, Trash2, UserPlus } from 'lucide-react';
import { fetchPosts, deletePost } from '../../lib/auth';
import VerifiedBadge from '../homepage/VerifiedBadge';
import CommentForm from './CommentForm';
import CommentList from './CommentList';

const PostList = ({ refreshTrigger, user }) => {
  const [posts, setPosts] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [expandedComments, setExpandedComments] = useState(new Set());
  const [newComments, setNewComments] = useState({});
  const [openDropdown, setOpenDropdown] = useState(null);
  const [deleteConfirmation, setDeleteConfirmation] = useState(null);

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
    loadPosts();
  }, [refreshTrigger]);

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
        <div className="text-gray-400">Loading posts...</div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="rounded-xl p-4 mb-4 bg-red-500 bg-opacity-20 border border-red-500">
        <div className="text-red-300 text-center">{error}</div>
        <button 
          onClick={loadPosts}
          className="mt-2 w-full py-2 px-4 bg-red-500 text-white rounded-lg hover:bg-red-600 transition-colors"
        >
          Try Again
        </button>
      </div>
    );
  }

  if (posts.length === 0) {
    return (
      <div className="rounded-xl p-8 text-center" style={{ backgroundColor: 'var(--primary-background)' }}>
        <div className="text-gray-400 mb-2">No posts yet</div>
        <div className="text-sm text-gray-500">Be the first to share something!</div>
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
                  <span className="font-medium text-white break-words">
                    {post.author?.nickname || `${post.author?.first_name || ''} ${post.author?.last_name || ''}`.trim() || 'User'}
                  </span>
                  <div className="flex items-center gap-1 text-xs" style={{ color: 'var(--secondary-text)' }}>
                    {React.createElement(getPrivacyIcon(post.privacy), { className: "w-3 h-3" })}
                    <span>{getPrivacyLabel(post.privacy)}</span>
                  </div>
                </div>
                <div className="text-xs" style={{ color: 'var(--secondary-text)' }}>
                  {formatDate(post.created_at)}
                </div>
              </div>
            </div>
            {/* Post Actions Dropdown */}
            <div className="relative dropdown-container">
              <button
                onClick={() => toggleDropdown(post.id)}
                className="text-gray-400 hover:text-white transition-colors"
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
                    {/* Edit Option */}
                    <button
                      className="w-full px-4 py-2 text-left text-sm text-white hover:text-white flex items-center gap-2 transition-colors"
                      style={{ ':hover': { backgroundColor: 'var(--hover-background)' } }}
                      onMouseEnter={(e) => e.target.style.backgroundColor = 'var(--hover-background)'}
                      onMouseLeave={(e) => e.target.style.backgroundColor = 'transparent'}
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
                      className="w-full px-4 py-2 text-left text-sm text-white hover:text-white flex items-center gap-2 transition-colors"
                      style={{ color: 'var(--danger-color)' }}
                      onMouseEnter={(e) => e.target.style.backgroundColor = 'var(--hover-background)'}
                      onMouseLeave={(e) => e.target.style.backgroundColor = 'transparent'}
                    >
                      <Trash2 className="w-4 h-4" />
                      Delete
                    </button>

                    {/* Follow Option */}
                    <button
                      className="w-full px-4 py-2 text-left text-sm text-white hover:text-white flex items-center gap-2 transition-colors"
                      onMouseEnter={(e) => e.target.style.backgroundColor = 'var(--hover-background)'}
                      onMouseLeave={(e) => e.target.style.backgroundColor = 'transparent'}
                    >
                      <UserPlus className="w-4 h-4" />
                      Follow
                    </button>
                  </div>
                </div>
              )}
            </div>
          </div>

          {/* Post Content */}
          <div className="mb-4">
            <p className="text-white whitespace-pre-wrap break-words overflow-wrap-anywhere">{post.content}</p>
            
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
          <div className="flex items-center justify-between pt-3 border-t border-[#3f3fd3]/30">
            <div className="flex items-center gap-4">
              <button className="flex items-center gap-2 text-sm py-1.5 px-3 rounded-lg hover:bg-[#3f3fd3]/20 transition-colors" style={{ color: 'var(--secondary-text)' }}>
                <ThumbsUpIcon className="w-4 h-4" />
                <span>Like</span>
              </button>
              <button className="flex items-center gap-2 text-sm py-1.5 px-3 rounded-lg hover:bg-[#3f3fd3]/20 transition-colors" style={{ color: 'var(--secondary-text)' }}>
                <ThumbsDownIcon className="w-4 h-4" />
                <span>Dislike</span>
              </button>
              <button
                onClick={() => toggleComments(post.id)}
                className="flex items-center gap-2 text-sm py-1.5 px-3 rounded-lg hover:bg-[#3f3fd3]/20 transition-colors"
                style={{ color: 'var(--secondary-text)' }}
              >
                <MessageCircleIcon className="w-4 h-4" />
                <span>{expandedComments.has(post.id) ? 'Hide Comments' : 'Comment'}</span>
              </button>
            </div>
          </div>

          {/* Comments Section */}
          {expandedComments.has(post.id) && (
            <div className="mt-4 pt-3 border-t border-[#3f3fd3]/30 space-y-4">
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
              />
            </div>
          )}
        </div>
      ))}

      {/* Delete Confirmation Modal */}
      {deleteConfirmation && (
        <div className="fixed inset-0 bg-black bg-opacity-30 flex items-center justify-center" style={{ zIndex: 9999 }}>
          <div
            className="rounded-lg p-6 max-w-md w-full mx-4"
            style={{ backgroundColor: 'var(--primary-background)', border: '1px solid var(--border-color)' }}
          >
            <h3 className="text-lg font-semibold text-white mb-4">Delete Post</h3>
            <p className="text-white mb-6">Are you sure you want to delete this post? This action cannot be undone.</p>

            <div className="flex gap-3 justify-end">
              <button
                onClick={() => setDeleteConfirmation(null)}
                className="px-4 py-2 rounded-lg text-white transition-colors"
                style={{ backgroundColor: 'var(--secondary-background)' }}
                onMouseEnter={(e) => e.target.style.backgroundColor = 'var(--hover-background)'}
                onMouseLeave={(e) => e.target.style.backgroundColor = 'var(--secondary-background)'}
              >
                Cancel
              </button>
              <button
                onClick={() => handleDeletePost(deleteConfirmation)}
                className="px-4 py-2 rounded-lg text-white transition-colors"
                style={{ backgroundColor: 'var(--danger-color)' }}
                onMouseEnter={(e) => e.target.style.opacity = '0.8'}
                onMouseLeave={(e) => e.target.style.opacity = '1'}
              >
                Delete
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default PostList;
