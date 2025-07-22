import React from 'react';
const CommunityItem = ({ icon, name, memberCount }) => {
  return (
    <div className="flex items-center gap-3">
      <div
        className="w-10 h-10 rounded-full overflow-hidden"
        style={{ backgroundColor: 'var(--tertiary-text)' }}
      >
        <img src={icon} alt={name} className="w-full h-full object-cover" />
      </div>
      <div className="flex-1">
        <p className="text-sm font-medium" style={{ color: 'var(--primary-text)' }}>{name}</p>
        <p className="text-xs" style={{ color: 'var(--primary-accent)' }}>
          • {memberCount} your friends are in
        </p>
      </div>
    </div>
  );
};
export default CommunityItem;