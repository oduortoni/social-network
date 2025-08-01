import React from 'react';
import ProfilePosts from './ProfilePosts';
import ProfileAbout from './ProfileAbout';
import ProfileFriends from './ProfileFriends';
import ProfilePhotos from './ProfilePhotos';
import ProfileVideos from './ProfileVideos';

const ProfileMainContent = ({ user, activeTab }) => {
  const renderContent = () => {
    switch (activeTab) {
      case 'posts':
        return <ProfilePosts user={user} />;
      case 'about':
        return <ProfileAbout user={user} />;
      case 'friends':
        return <ProfileFriends user={user} />;
      case 'photos':
        return <ProfilePhotos user={user} />;
      case 'videos':
        return <ProfileVideos user={user} />;
      default:
        return <ProfilePosts user={user} />;
    }
  };

  return (
    <div className="w-full">
      {renderContent()}
    </div>
  );
};

export default ProfileMainContent;
