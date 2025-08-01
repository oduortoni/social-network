import React from 'react';
import PostCreation from '../posts/PostCreation';
import PostList from '../posts/PostList';

const ProfilePosts = ({ user }) => {
  return (
    <div className="space-y-4">
      {/* Post Creation */}
      <PostCreation user={user} />
      
      {/* User's Posts */}
      <PostList user={user} profileView={true} />
    </div>
  );
};

export default ProfilePosts;