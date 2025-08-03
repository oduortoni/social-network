import React from 'react';
import PostCreation from '../posts/PostCreation';
import PostList from '../posts/PostList';

const ProfilePosts = ({ user }) => {
  console.log('ProfilePosts rendered with user:', user, 'and posts:', user?.posts);
  
  return (
    <div className="space-y-4">
      {/* Post Creation */}
      <PostCreation user={user} />
      
      {/* User's Posts */}
     {user?.posts && (
  <PostList user={user} posts={user.posts} profileView={true} />
)}
      {/* Fallback if no posts */}
      {!user?.posts || user.posts.length === 0 ? (
        <div className="text-center text-gray-500">No posts available.</div>
      ) : null}
    </div>
  );
};

export default ProfilePosts;