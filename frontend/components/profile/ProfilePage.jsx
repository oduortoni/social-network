import React, { useState, useEffect } from 'react';
import Header from '../layout/Header';
import ProfileCover from './ProfileCover';
import ProfileNavigation from './ProfileNavigation';
import ProfileMainContent from './ProfileMainContent';

const ProfilePage = ({ user }) => {
  const [activeTab, setActiveTab] = useState('posts');
  const [profileData, setProfileData] = useState(null);

  useEffect(() => {
    // TODO: Fetch profile data from backend
    setProfileData(user);
  }, [user]);

  return (
    <div className="w-2/3 flex flex-col text-white px-4 py-2 mt-4">
      <Header user={user}/>
      <div className="flex-1 w-full mt-4">
        {/* Cover and Profile Info Section */}
        <ProfileCover user={profileData || user} />
        
        {/* Navigation Tabs */}
        <ProfileNavigation activeTab={activeTab} setActiveTab={setActiveTab} />
        
        {/* Main Content Area */}
        <div className="w-full px-12 py-6">
          <ProfileMainContent user={profileData || user} activeTab={activeTab} />
        </div>
      </div>
    </div>
  );
};

export default ProfilePage;
