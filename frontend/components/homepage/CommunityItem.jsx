import React from 'react';
const CommunityItem = ({ icon, name, memberCount }) => {
  return <div className="flex items-center gap-3">
      <div className="w-10 h-10 rounded-full overflow-hidden bg-[#3f3fd3]">
        <img src={icon} alt={name} className="w-full h-full object-cover" />
      </div>
      <div className="flex-1">
        <p className="text-sm font-medium">{name}</p>
        <p className="text-xs text-[#ffd700]">
          • {memberCount} your friends are in
        </p>
      </div>
    </div>;
};
export default CommunityItem;