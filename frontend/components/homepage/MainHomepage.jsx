import React, { useState } from 'react';
import Header from '../layout/Header';
import ProfileSidebar from '../layout/ProfileSidebar';
import Feed from '../layout/Feed';
import ActivitySidebar from '../layout/ActivitySidebar';
import UserListModal from '../chat/UserListModal';

const MainHomepage = ({ user, profile, connectionStatus, connectedUsers = [] }) => {
  const [showUserListModal, setShowUserListModal] = useState(false);
  return (
    <div className="w-2/3 flex flex-col text-white">
      <Header user={user} />
      <div className="flex flex-1 w-full max-w-7xl mx-auto gap-4 p-4">
        <ProfileSidebar profile={profile} connectionStatus={connectionStatus} />
        <div className="flex-1 flex flex-col">
          <Feed user={user} connectedUsers={connectedUsers} />
        </div>
        <ActivitySidebar user={user} />
      </div>
      {/* Start Conversation Button */}
      <button
        className="fixed bottom-8 right-8 bg-blue-500 text-white p-4 rounded-full shadow-lg hover:bg-blue-600 transition-colors"
        onClick={() => setShowUserListModal(true)}
      >
        Start Conversation
      </button>

      {showUserListModal && (
        <UserListModal user={user} onClose={() => setShowUserListModal(false)} />
      )}
    </div>
  );
};

export default MainHomepage; 