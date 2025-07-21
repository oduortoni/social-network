import React from 'react';
const ActivityItem = ({
  image,
  name,
  action,
  time,
  isGroup = false,
  isPartial = false
}) => {
  if (isPartial) {
    return (
      <div className="flex items-center gap-3">
        <img src={image} alt={name} className="w-10 h-10 rounded-full" />
        <div className="flex-1">
          <div className="flex items-center gap-1">
            <span className="font-medium text-sm">{name}</span>
            <span className="text-sm text-[#afafed]">{action}</span>
          </div>
        </div>
      </div>
    );
  }
  return (
    <div className="flex items-start gap-3">
      <div className="relative">
        <img src={image} alt={name} className="w-10 h-10 rounded-full" />
        {isGroup && (
          <div className="absolute -top-1 -right-1 w-4 h-4 bg-[#ff4444] rounded-full flex items-center justify-center text-[8px]">
            3
          </div>
        )}
      </div>
      
      <div className="flex-1">
        <div className="flex items-center gap-1">
          <span className="font-medium text-sm">{name}</span>
          <span className="text-sm text-[#afafed]">{action}</span>
        </div>
        {time && <p className="text-xs text-[#afafed]">{time}</p>}
        <div className="flex justify-between items-center mt-2 gap-2">
          <button className="bg-[#ffd700] text-black px-4 py-1 rounded-full text-xs flex-1">Accept</button>
          <button className="bg-[#101b70] border border-[#3f3fd3]/30 px-4 py-1 rounded-full text-xs flex-1">Decline</button>
        </div>
      </div>
    </div>
  );
};
export default ActivityItem;