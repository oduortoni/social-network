import React, { useState } from 'react';
import UserCircle from '../homepage/UserCircle';
import PostCreation from '../posts/PostCreation';
import PostList from '../posts/PostList';

const Feed = ({ user }) => {
  const [refreshTrigger, setRefreshTrigger] = useState(0);

  const handlePostCreated = (newPost) => {
    // Trigger a refresh of the post list
    setRefreshTrigger(prev => prev + 1);
  };

  return <div className="flex-1 flex flex-col gap-4">
      {/* Stories Section */}
      <div className="flex overflow-x-auto gap-3 pb-2 cursor-pointer">
        <UserCircle image="https://randomuser.me/api/portraits/women/22.jpg" name="Amanda" active={false} highlight="#3f3fd3" />
        <UserCircle image="https://randomuser.me/api/portraits/men/22.jpg" name="John" active={false} highlight="#ff4444" />
        <UserCircle image="https://randomuser.me/api/portraits/men/32.jpg" name="Andrew" active={false} highlight="#3f3fd3" />
        <UserCircle image="https://randomuser.me/api/portraits/women/32.jpg" name="Rosaline" active={false} highlight="#00ff00" />
        <UserCircle image="https://randomuser.me/api/portraits/men/42.jpg" name="Mudreh" active={false} highlight="#ff4444" />
        <UserCircle image="https://randomuser.me/api/portraits/women/42.jpg" name="Juliet" active={false} highlight="#3f3fd3" />
        <UserCircle image="https://randomuser.me/api/portraits/men/52.jpg" name="Bob" active={false} highlight="#3f3fd3" />
        <UserCircle image="https://randomuser.me/api/portraits/men/2.jpg" name="Mudreh" active={false} highlight="#ff4444" />
        <UserCircle image="https://randomuser.me/api/portraits/women/4.jpg" name="Juliet" active={false} highlight="#3f3fd3" />
      </div>

      {/* Post Creation Component */}
      <PostCreation user={user} onPostCreated={handlePostCreated} />
      {/* Posts List */}
      <PostList refreshTrigger={refreshTrigger} />
    </div>;
};

export default Feed;