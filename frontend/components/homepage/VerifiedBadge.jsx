import React from 'react';
const VerifiedBadge = () => {
  return (
    <div className="rounded-full p-[2px]" style={{ backgroundColor: 'var(--tertiary-text)' }}>
      <svg width="12" height="12" viewBox="0 0 24 24" fill="none">
        <path d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" stroke="white" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" />
      </svg>
    </div>
  );
};
export default VerifiedBadge;