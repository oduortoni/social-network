import React, { useState } from 'react';
import UserCircle from '../homepage/UserCircle';
import PostCreation from '../posts/PostCreation';
import PostList from '../posts/PostList';

const Feed = ({user = null, connectedUsers = []}) => {
    const users = connectedUsers.filter(u => u.user_id != user.id);
    const [refreshTrigger, setRefreshTrigger] = useState(0);

    const handlePostCreated = (newPost) => {
      // Trigger a refresh of the post list
      setRefreshTrigger(prev => prev + 1);
    };

  return <div className="flex-1 flex flex-col gap-4">
      {/* Content Container with max-width */}
      <div className="w-full max-w-2xl mx-auto">
        {/* Stories Section */}
        {
          users.map((u, index) => {
              return <UserCircle avatar={u.avatar? u.avatar : ''} name={u.nickname} active={false} highlight="#3f3fd3" key={index} />;
          })
        }

        {/* Post Creation Component */}
        <PostCreation user={user} onPostCreated={handlePostCreated} />
        {/* Posts List */}
        <PostList refreshTrigger={refreshTrigger} user={user} />
      </div>
    </div>;
};

export default Feed;