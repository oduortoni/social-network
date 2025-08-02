import React from 'react';
import { profileAPI } from '../../lib/api';

const ProfileFriends = ({ user }) => {
  // Mock friends data
  const friends = Array.from({ length: 12 }, (_, i) => ({
    id: i + 1,
    first_name: `Friend ${i + 1}`,
    last_name: 'User',
    nickname: `friend${i + 1}`,
    avatar: '',
    mutualFriends: Math.floor(Math.random() * 50)
  }));

  return (
    <div className="space-y-4">
      <div
        className="rounded-xl p-6"
        style={{ backgroundColor: 'var(--primary-background)' }}
      >
        <h3 className="text-xl font-bold mb-4 text-white">Friends ({friends.length})</h3>
        <div className="grid grid-cols-2 md:grid-cols-3 gap-4">
          {friends.map((friend) => (
            <div
              key={friend.id}
              className="rounded-lg p-4 text-center cursor-pointer hover:opacity-80 transition-opacity"
              style={{ backgroundColor: 'var(--secondary-background)' }}
            >
              <img
                src={profileAPI.fetchProfileImage(friend.avatar)}
                alt={friend.first_name}
                className="w-20 h-20 rounded-full mx-auto mb-3"
              />
              <h4 className="font-medium text-white text-sm">
                {friend.first_name} {friend.last_name}
              </h4>
              <p className="text-xs mt-1" style={{ color: 'var(--secondary-text)' }}>
                {friend.mutualFriends} mutual friends
              </p>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
};

export default ProfileFriends;