import React from 'react';
import { EyeIcon } from 'lucide-react';

const ProfileGroups = ({ user }) => {
  // Mock groups data
  const groups = Array.from({ length: 6 }, (_, i) => ({
    id: i + 1,
    image: 'cat.jpg',
    title: `Group ${i + 1}`,
    about: 'This is a short description about the group.'
  }));

  return (
    <div className="space-y-4">
      <div
        className="rounded-xl p-6"
        style={{ backgroundColor: 'var(--primary-background)' }}
      >
        <h3 className="text-xl font-bold mb-4 text-white">Groups ({groups.length})</h3>
        <div className="grid grid-cols-2 gap-4">
          {groups.map((group) => (
            <div
              key={group.id}
              className="relative aspect-video rounded-lg overflow-hidden cursor-pointer hover:opacity-90 transition-opacity"
              style={{ backgroundColor: 'var(--secondary-background)' }}
            >
              <img
                src={group.image}
                alt={group.title}
                className="w-full h-full object-cover"
              />
              <div className="absolute inset-0 bg-black/50 flex flex-col justify-end p-4">
                <div className="text-white">
                  <h4 className="text-lg font-semibold">{group.title}</h4>
                  <p className="text-sm opacity-80">{group.about}</p>
                  <button
                    className="mt-2 px-4 py-1.5 flex items-center gap-2 text-sm rounded-lg border"
                    style={{
                      backgroundColor: 'transparent',
                      borderColor: 'var(--border-color)',
                      color: 'var(--primary-text)'
                    }}
                  >
                    <EyeIcon className="w-4 h-4" />
                    View Group
                  </button>
                </div>
              </div>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
};

export default ProfileGroups;
