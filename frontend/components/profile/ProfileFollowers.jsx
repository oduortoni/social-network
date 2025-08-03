import React, { useState, useEffect } from 'react';
import { profileAPI } from '../../lib/api';

const ProfileFollowers = ({ user }) => {
  const [followers, setFollowers] = useState([]);

  useEffect(() => {
    const fetchFollowers = async () => {
      if (user && user.profile_details && user.profile_details.id) {
        try {
          const data = await profileAPI.getFollowers(user.profile_details.id);
          setFollowers(data.user || []);
        } catch (error) {
          console.error('Error fetching followers:', error);
        }
      }
    };

    
    fetchFollowers();
  }, [user]);

  return (
    <div className="space-y-4">
      <div
        className="rounded-xl p-6"
        style={{ backgroundColor: 'var(--primary-background)' }}
      >
        <h3 className="text-xl font-bold mb-4 text-white">Followers ({followers.length})</h3>
        <div className="grid grid-cols-2 md:grid-cols-3 gap-4">
          {followers && followers.length > 0 ? (
            followers.map((follower) => (
              <div
                key={follower.id}
                className="rounded-lg p-4 text-center cursor-pointer hover:opacity-80 transition-opacity"
                style={{ backgroundColor: 'var(--secondary-background)' }}
              >
                <img
                  src={profileAPI.fetchProfileImage(follower.avatar)}
                  alt={follower.first_name}
                  className="w-20 h-20 rounded-full mx-auto mb-3"
                />
                <h4 className="font-medium text-white text-sm">
                  {follower.first_name} {follower.last_name}
                </h4>
                <p className="text-xs mt-1" style={{ color: 'var(--secondary-text)' }}>
                  {follower.mutualFollowers} mutual followers
                </p>
              </div>
            ))
          ) : (
            <p className="text-white">No followers yet. Start sharing your content to attract followers!</p>
          )}
        </div>
      </div>
    </div>
  );
};

export default ProfileFollowers;