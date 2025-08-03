import React from 'react';

const ProfileNavigation = ({ activeTab, setActiveTab }) => {
  const tabs = [
    { id: 'posts', label: 'Posts', count: 1234 },
    { id: 'about', label: 'About' },
    { id: 'followers', label: 'Followers', count: 12 },
    { id: 'following', label: 'Following', count: 8 },
    { id: 'photos', label: 'Photos', count: 89 },
    { id: 'videos', label: 'My Groups', count: 8 }
  ];

  return (
    <div className="border-b" style={{ borderColor: 'var(--border-color)' }}>
      <div className="w-full px-6">
        <div className="flex gap-8">
          {tabs.map((tab) => (
            <button
              key={tab.id}
              onClick={() => setActiveTab(tab.id)}
              className={`py-4 px-2 font-medium cursor-pointer transition-all relative ${
                activeTab === tab.id
                  ? 'text-white'
                  : 'hover:text-white'
              }`}
              style={{ 
                color: activeTab === tab.id ? 'var(--primary-text)' : 'var(--secondary-text)'
              }}
            >
              {tab.label}
              {tab.count && (
                <span className="ml-2 text-sm" style={{ color: 'var(--secondary-text)' }}>
                  ({tab.count})
                </span>
              )}
              {activeTab === tab.id && (
                <div
                  className="absolute bottom-0 left-0 right-0 h-0.5"
                  style={{ backgroundColor: 'var(--primary-accent)' }}
                />
              )}
            </button>
          ))}
        </div>
      </div>
    </div>
  );
};

export default ProfileNavigation;
