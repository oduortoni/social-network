"use client"
import React from 'react';
import { MapPinIcon, CalendarIcon, LinkIcon, BriefcaseIcon, MailIcon } from 'lucide-react';

const ProfileAbout = ({ user }) => {
  const profileDetails = user?.profile_details || {};

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
  <CalendarIcon className="w-5 h-5" style={{ color: 'var(--secondary-text)' }} />
  <span style={{ color: 'var(--primary-text)' }}>
    <strong>Date of Birth:</strong> {user.profile_details.dateofbirth.split('-').reverse().join('-')}
  </span>
</div>

<div className="flex items-center gap-3">
  <MailIcon className="w-5 h-5" style={{ color: 'var(--secondary-text)' }} />
  <span style={{ color: 'var(--primary-text)' }}>
    <strong>Email:</strong> {user.profile_details.email}
  </span>
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
          {profileDetails.about || "This user hasn't written a bio yet."}
        </p>
      </div>
    </div>
  );
}
export default ProfileAbout;