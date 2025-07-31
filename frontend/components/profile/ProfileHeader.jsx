
import { useState } from 'react';

const ProfileHeader = () => {
  const [followBtnStatus, setFollowBtnStatus] = useState('follow');
  const [messageBtnStatus, setMessageBtnStatus] = useState('visible');
  const [editBtnStatus, setEditBtnStatus] = useState('visible');

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
    <div className="bg-white shadow-sm p-4 rounded-lg">
      <div className="flex items-center space-x-4">
        <img src="/cat.jpg" alt="User Avatar" className="w-24 h-24 rounded-full" />
        <div>
          <h1 className="text-2xl font-bold">User Name</h1>
          <p className="text-gray-600">user@gmail.com</p>
          <div className="flex space-x-4 mt-2">
            <p><span className="font-bold">100</span> Posts</p>
            <p><span className="font-bold">500</span> Followers</p>
            <p><span className="font-bold">200</span> Following</p>
          </div>
        </div>
      </div>
      <div className="mt-4 flex space-x-2">
        {followBtnStatus !== 'hide' && (
          <button onClick={handleFollow} className="bg-blue-500 text-white px-4 py-2 rounded-lg">
            {followBtnStatus === 'follow' && 'Follow'}
            {followBtnStatus === 'pending' && 'Follow request sent'}
            {followBtnStatus === 'following' && 'Following'}
          </button>
        )}
        {messageBtnStatus === 'visible' && (
          <button className="bg-gray-200 px-4 py-2 rounded-lg">Message</button>
        )}
        {editBtnStatus === 'visible' && (
          <button className="bg-gray-200 px-4 py-2 rounded-lg">Edit</button>
        )}
      </div>
    </div>
  );
};

export default ProfileHeader;
