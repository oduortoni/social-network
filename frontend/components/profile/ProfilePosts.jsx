import React from 'react';
import PostCreation from '../posts/PostCreation';
import PostList from '../posts/PostList';

const ProfilePosts = ({ user, currentUser, isOwnProfile }) => {
  console.log('ProfilePosts rendered with user:', user, 'and posts:', user?.posts);
  
  return (
    <div className="space-y-4">
      {/* Post Creation - Only show for own profile */}
      {isOwnProfile && <PostCreation user={currentUser || user} />}
      
      {/* User's Posts */}
      {user?.posts && (
        <PostList user={user} posts={user.posts} profileView={true} />
      )}
      
      {/* Fallback if no posts */}
      {!user?.posts || user.posts.length === 0 ? (
        <div className="text-center py-8" style={{ color: 'var(--secondary-text)' }}>
          {isOwnProfile ? "You haven't posted anything yet." : "No posts available."}
        </div>
      ) : null}
    </div>
  );
};

export default ProfilePosts;