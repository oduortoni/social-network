import React from 'react';

const ProfilePhotos = ({ user, currentUser, isOwnProfile }) => {
  // Mock photos data
  const photos = Array.from({ length: 15 }, (_, i) => ({
    id: i + 1,
    url: 'cat.jpg',
    caption: `Photo ${i + 1}`
  }));

  return (
    <div className="space-y-4">
      <div
        className="rounded-xl p-6"
        style={{ backgroundColor: 'var(--primary-background)' }}
      >
        <h3 className="text-xl font-bold mb-4 text-white">Photos ({photos.length})</h3>
        <div className="grid grid-cols-3 gap-2">
          {photos.map((photo) => (
            <div
              key={photo.id}
              className="aspect-square rounded-lg overflow-hidden cursor-pointer hover:opacity-80 transition-opacity"
              style={{ backgroundColor: 'var(--secondary-background)' }}
            >
              <img
                src={photo.url} // Replace with actual photo URL  
                alt={photo.caption}
                className="w-full h-full object-cover"
              />
            </div>
          ))}
        </div>
      </div>
    </div>
  );
};

export default ProfilePhotos;