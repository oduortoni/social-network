import React from 'react';
import { Lock, UserPlus } from 'lucide-react';
import ProfilePosts from './ProfilePosts';
import ProfileAbout from './ProfileAbout';
import ProfileFollowers from './ProfileFollowers';
import ProfileFollowing from './ProfileFollowing';
import ProfilePhotos from './ProfilePhotos';
import ProfileVideos from './ProfileVideos';

const ProfileMainContent = ({ user, currentUser, isOwnProfile, activeTab }) => {
        
  
  const profileDetails = user?.profile_details || {};
  const isPrivateAccount = profileDetails.profile === false;
  const followStatus = profileDetails.followbtnstatus;
  
  // Check if user can view private content
  const canViewPrivateContent = isOwnProfile ||
    (isPrivateAccount && followStatus === 'following') ||
    !isPrivateAccount;

   console.log("view account",canViewPrivateContent)
  // Private Account Message Component
  const PrivateAccountMessage = () => (
    <div className="flex flex-col items-center justify-center py-16 px-8 text-center">
      <div className="w-24 h-24 rounded-full flex items-center justify-center mb-6"
           style={{ backgroundColor: 'var(--secondary-background)' }}>
        <Lock className="w-12 h-12" style={{ color: 'var(--secondary-text)' }} />
      </div>
      
      <h3 className="text-2xl font-bold text-white mb-4">
        This Account is Private
      </h3>
      
      <p className="text-lg mb-6" style={{ color: 'var(--secondary-text)' }}>
        Follow @{profileDetails.nickname} to see their posts, photos, and other content.
      </p>
      
      {followStatus === 'follow' && (
        <div className="flex items-center gap-2 px-6 py-3 rounded-lg"
             style={{ backgroundColor: 'var(--primary-accent)', color: 'white' }}>
          <UserPlus className="w-5 h-5" />
          <span>Send Follow Request</span>
        </div>
      )}
      
      {followStatus === 'pending' && (
        <div className="flex items-center gap-2 px-6 py-3 rounded-lg border"
             style={{ borderColor: 'var(--border-color)', color: 'var(--secondary-text)' }}>
          <span>Follow Request Sent</span>
        </div>
      )}
    </div>
  );

  const renderContent = () => {
    // If it's a private account and user can't view content, show privacy message
    if (isPrivateAccount && !canViewPrivateContent) {
      return <PrivateAccountMessage />;
    }

    // Otherwise render normal content
    switch (activeTab) {
      case 'posts':
        return <ProfilePosts user={user} currentUser={currentUser} isOwnProfile={isOwnProfile} />;
      case 'about':
        return <ProfileAbout user={user} currentUser={currentUser} isOwnProfile={isOwnProfile} />;
      case 'followers':
        return <ProfileFollowers user={user} currentUser={currentUser} isOwnProfile={isOwnProfile} />;
      case 'following':
        return <ProfileFollowing user={user} currentUser={currentUser} isOwnProfile={isOwnProfile} />;
      case 'photos':
        return <ProfilePhotos user={user} currentUser={currentUser} isOwnProfile={isOwnProfile} />;
      case 'videos':
        return <ProfileVideos user={user} currentUser={currentUser} isOwnProfile={isOwnProfile} />;
      default:
        return <ProfilePosts user={user} currentUser={currentUser} isOwnProfile={isOwnProfile} />;
    }
  };

  return (
    <div className="w-full">
      {renderContent()}
    </div>
  );
};

export default ProfileMainContent;
