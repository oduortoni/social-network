import React from 'react';
import { CameraIcon, EditIcon, UserPlusIcon } from 'lucide-react';
import { profileAPI } from '../../lib/api';

const ProfileCover = ({ user }) => {
  return (
    <div className="w-full">
      {/* Cover Photo */}
      <div className="relative h-80 w-full overflow-hidden">
        <div 
          className="absolute inset-0"
          style={{ 
            backgroundColor: 'var(--secondary-background)',
            backgroundImage: 'url(/assets/background_textures/grain.png)',
            backgroundBlendMode: 'overlay'
          }}
        />
        
        {/* Cover Photo Edit Button */}
        <button className="absolute top-4 right-4 p-2 rounded-lg bg-black bg-opacity-50 hover:bg-opacity-70 transition-all">
          <CameraIcon className="w-5 h-5 text-white" />
        </button>
      </div>

      {/* Profile Info Section */}
      <div className="relative -mt-20 px-6 pb-4">
        <div className="flex flex-col md:flex-row md:items-end md:justify-between">
          {/* Profile Picture and Basic Info */}
          <div className="flex flex-col md:flex-row md:items-end gap-4">
            {/* Profile Picture */}
            <div className="relative">
              <div
                className="w-40 h-40 rounded-full p-1"
                style={{ backgroundColor: 'var(--primary-accent)' }}
              >
                <img
                  src={profileAPI.fetchProfileImage(user.avatar || '')}
                  alt="Profile"
                  className="w-full h-full rounded-full object-cover"
                />
              </div>
              <button className="absolute bottom-2 right-2 p-2 rounded-full" style={{ backgroundColor: 'var(--tertiary-text)' }}>
                <CameraIcon className="w-4 h-4 text-white" />
              </button>
            </div>

            {/* User Info */}
            <div className="pb-4 mt-4">
              <h1 className="text-3xl font-bold text-white">
                {user.first_name} {user.last_name}
              </h1>
              <p className="text-lg" style={{ color: 'var(--secondary-text)' }}>
                @{user.nickname}
              </p>
              <p className="mt-2 max-w-md" style={{ color: 'var(--primary-text)' }}>
                {user.about_me || "No bio available"}
              </p>
              
              {/* Stats */}
              <div className="flex gap-6 mt-3">
                <div className="text-center">
                  <div className="font-bold text-white">1,234</div>
                  <div className="text-sm" style={{ color: 'var(--secondary-text)' }}>Posts</div>
                </div>
                <div className="text-center">
                  <div className="font-bold text-white">5,678</div>
                  <div className="text-sm" style={{ color: 'var(--secondary-text)' }}>Followers</div>
                </div>
                <div className="text-center">
                  <div className="font-bold text-white">890</div>
                  <div className="text-sm" style={{ color: 'var(--secondary-text)' }}>Following</div>
                </div>
              </div>
            </div>
          </div>

          {/* Action Buttons */}
          <div className="flex gap-3 mt-4 md:mt-0">
            <button
              className="px-6 py-2 rounded-lg font-medium flex items-center gap-2 cursor-pointer"
              style={{
                backgroundColor: 'var(--tertiary-text)',
                color: 'var(--primary-text)',
              }}
            >
              <EditIcon className="w-4 h-4" />
              Edit Profile
            </button>
            <button
              className="flex items-center gap-2 px-6 py-2 rounded-lg font-medium border cursor-pointer"
              style={{
                backgroundColor: 'transparent',
                borderColor: 'var(--border-color)',
                color: 'var(--primary-text)',
              }}
            >
              <UserPlusIcon className="w-4 h-4" />
              Follow
            </button>
          </div>
        </div>
      </div>
    </div>
  );
};

export default ProfileCover;