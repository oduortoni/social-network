import React from 'react';
const UserCircle = ({ image, name, active = false, highlight = '#ffd700' }) => {
  return <div className="flex flex-col items-center gap-1 min-w-[60px]">
      <div className={`rounded-full p-[2px]`} style={{
      background: highlight
    }}>
        <div className="bg-[#101b70] rounded-full p-[2px]">
          <img src={image} alt={name} className="w-12 h-12 rounded-full object-cover" />
        </div>
      </div>
      <span className="text-xs">{name}</span>
    </div>;
};
export default UserCircle;