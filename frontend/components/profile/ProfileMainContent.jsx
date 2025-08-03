import React from 'react';
import ProfilePosts from './ProfilePosts';
import ProfileAbout from './ProfileAbout';
import ProfileFollowers from './ProfileFollowers';
import ProfileFollowing from './ProfileFollowing';
import ProfilePhotos from './ProfilePhotos';
import ProfileVideos from './ProfileVideos';

const ProfileMainContent = ({ user, activeTab }) => {
  console.log('ProfileMainContent rendered with user:', user, 'and activeTab:', activeTab);
  const renderContent = () => {
    switch (activeTab) {
      case 'posts':
        return <ProfilePosts user={user} />;
      case 'about':
        return <ProfileAbout user={user} />;
      case 'followers':
        return <ProfileFollowers user={user} />;
      case 'following':
        return <ProfileFollowing user={user} />;
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
