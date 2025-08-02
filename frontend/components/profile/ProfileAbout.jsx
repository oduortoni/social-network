import React from 'react';
import { MapPinIcon, CalendarIcon, LinkIcon, BriefcaseIcon } from 'lucide-react';

const ProfileAbout = ({ user }) => {
  return (
    <div className="space-y-6">
      {/* Basic Info */}
      <div
        className="rounded-xl p-6"
        style={{ backgroundColor: 'var(--primary-background)' }}
      >
        <h3 className="text-xl font-bold mb-4 text-white">About</h3>
        <div className="space-y-4">
          <div className="flex items-center gap-3">
            <BriefcaseIcon className="w-5 h-5" style={{ color: 'var(--secondary-text)' }} />
            <span style={{ color: 'var(--primary-text)' }}>Software Developer</span>
          </div>
          <div className="flex items-center gap-3">
            <MapPinIcon className="w-5 h-5" style={{ color: 'var(--secondary-text)' }} />
            <span style={{ color: 'var(--primary-text)' }}>San Francisco, CA</span>
          </div>
          <div className="flex items-center gap-3">
            <CalendarIcon className="w-5 h-5" style={{ color: 'var(--secondary-text)' }} />
            <span style={{ color: 'var(--primary-text)' }}>
              Joined {new Date(user.created_at).toLocaleDateString('en-US', { month: 'long', year: 'numeric' })}
            </span>
          </div>
          <div className="flex items-center gap-3">
            <LinkIcon className="w-5 h-5" style={{ color: 'var(--secondary-text)' }} />
            <a href="#" className="hover:underline" style={{ color: 'var(--tertiary-text)' }}>
              portfolio.example.com
            </a>
          </div>
        </div>
      </div>

      {/* Bio */}
      <div
        className="rounded-xl p-6"
        style={{ backgroundColor: 'var(--primary-background)' }}
      >
        <h3 className="text-xl font-bold mb-4 text-white">Bio</h3>
        <p style={{ color: 'var(--primary-text)' }} className="leading-relaxed">
          {user.about_me || "This user hasn't written a bio yet."}
        </p>
      </div>
    </div>
  );
};

export default ProfileAbout;