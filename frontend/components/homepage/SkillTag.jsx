import React from 'react';
const SkillTag = ({ label }) => {
  return (
    <div
      className="px-3 py-1 rounded-lg text-xs border"
      style={{
        backgroundColor: 'var(--primary-background)',
        borderColor: 'var(--tertiary-text)',
        color: 'var(--primary-text)',
      }}
    >
      {label}
    </div>
  );
};
export default SkillTag;