
import { useState } from 'react';

const FollowerList = () => {
  const followers = ['User 1', 'User 2', 'User 3'];

  return (
    <div className="mt-8">
      <h2 className="text-xl font-bold mb-4">Followers</h2>
      <ul>
        {followers.map((follower) => (
          <li key={follower} className="flex items-center justify-between py-2">
            <div className="flex items-center space-x-4">
              <img src="/cat.jpg" alt="User Avatar" className="w-12 h-12 rounded-full" />
              <span>{follower}</span>
            </div>
            <FollowButton />
          </li>
        ))}
      </ul>
    </div>
  );
};

const FollowButton = () => {
  const [followBtnStatus, setFollowBtnStatus] = useState('follow');

  const handleFollow = () => {
    if (followBtnStatus === 'follow') {
      setFollowBtnStatus('pending');
    } else if (followBtnStatus === 'pending') {
      setFollowBtnStatus('following');
    } else {
      setFollowBtnStatus('follow');
    }
  };

  return (
    <button onClick={handleFollow} className="bg-blue-500 text-white px-4 py-2 rounded-lg">
      {followBtnStatus === 'follow' && 'Follow'}
      {followBtnStatus === 'pending' && 'Follow request sent'}
      {followBtnStatus === 'following' && 'Following'}
    </button>
  );
};

export default FollowerList;
