import React from 'react';
import { profileAPI } from '../../lib/api';


const UserCircle = ({ avatar, name, active = false, highlight = 'var(--primary-accent)' }) => {
  return (
    <div className="flex flex-col items-center gap-1 min-w-[60px]">
      <div className="rounded-full p-[2px]" style={{ background: highlight }}>
        <div className="rounded-full p-[2px]" style={{ backgroundColor: 'var(--secondary-background)' }}>
          <img src={profileAPI.fetchProfileImage(avatar)} alt="Profile" className="w-12 h-12 rounded-full object-cover" />
        </div>
      </div>
      <span className="text-xs" style={{ color: 'var(--primary-text)' }}>{name}</span>
    </div>
  );
};
export default UserCircle;