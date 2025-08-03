import React, { useState, useEffect } from 'react';
import { profileAPI } from '../../lib/api';

const ProfileFollowing = ({ user }) => {
  const [following, setFollowing] = useState([]);

  useEffect(() => {
    const fetchFollowing = async () => {
      if (user && user.profile_details && user.profile_details.id) {
        try {
          const data = await profileAPI.getFollowing(user.profile_details.id);
          setFollowing(data.user || []);
        } catch (error) {
          console.error('Error fetching following:', error);
        }
      }
    };

    fetchFollowing();
  }, [user]);

  return (
    <div className="space-y-4">
      <div
        className="rounded-xl p-6"
        style={{ backgroundColor: 'var(--primary-background)' }}
      >
        <h3 className="text-xl font-bold mb-4 text-white">Following ({following.length})</h3>
        <div className="grid grid-cols-2 md:grid-cols-3 gap-4">
          {following && following.length > 0 ? (
            following.map((followedUser) => (
              <div
                key={followedUser.id}
                className="rounded-lg p-4 text-center cursor-pointer hover:opacity-80 transition-opacity"
                style={{ backgroundColor: 'var(--secondary-background)' }}
              >
                <img
                  src={profileAPI.fetchProfileImage(followedUser.avatar)}
                  alt={followedUser.first_name}
                  className="w-20 h-20 rounded-full mx-auto mb-3"
                />
                <h4 className="font-medium text-white text-sm">
                  {followedUser.first_name} {followedUser.last_name}
                </h4>
              </div>
            ))
          ) : (
            <p className="text-white">Not following anyone yet. Explore and connect with other users!</p>
          )}
        </div>
      </div>
    </div>
  );
};

export default ProfileFollowing;