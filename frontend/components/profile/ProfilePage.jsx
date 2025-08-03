import React, { useState, useEffect, useCallback } from 'react';
import Header from '../layout/Header';
import ProfileCover from './ProfileCover';
import ProfileNavigation from './ProfileNavigation';
import ProfileMainContent from './ProfileMainContent';

import { profileAPI } from '../../lib/api';

const ProfilePage = ({ user, userId }) => {
  const [activeTab, setActiveTab] = useState('posts');
  const [profileData, setProfileData] = useState(null);

  const fetchProfileData = useCallback(async () => {
    const targetUserId = userId || (user && user.id);
    if (targetUserId) {
      try {
        const data = await profileAPI.getProfile(targetUserId);
        setProfileData(data);
      } catch (error) {
        console.error('Error fetching profile data:', error);
      }
    }
  }, [user, userId]);

  useEffect(() => {
    fetchProfileData();
  }, [fetchProfileData]);

  return (
    <div className="w-2/3 flex flex-col text-white px-4 py-2 mt-4">
      <Header user={user}/>
      <div className="flex-1 w-full mt-4">
        {/* Cover and Profile Info Section */}
        <ProfileCover user={profileData} refreshProfile={fetchProfileData} />
        
        {/* Navigation Tabs */}
        <ProfileNavigation activeTab={activeTab} setActiveTab={setActiveTab} />
        
        {/* Main Content Area */}
        <div className="w-full px-12 py-6">
          <ProfileMainContent user={profileData} activeTab={activeTab} />
        </div>
      </div>
    </div>
  );
};

export default ProfilePage;


