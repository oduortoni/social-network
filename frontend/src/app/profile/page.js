'use client';

import { useState, useEffect } from 'react';
import withAuth from '../../../lib/withAuth';
import ProfilePage from '../../../components/profile/ProfilePage';

const Profile = ({ user }) => {
  return (
    <div className="min-h-screen">
      <main className="flex flex-col items-center justify-center">
        <ProfilePage user={user} />
      </main>
    </div>
  );
};

export default withAuth(Profile);