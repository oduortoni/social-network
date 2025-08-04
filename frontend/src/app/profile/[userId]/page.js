'use client';

import { useState, useEffect } from 'react';
import { useParams } from 'next/navigation';
import withAuth from '../../../../lib/withAuth';
import ProfilePage from '../../../../components/profile/ProfilePage';

const UserProfile = ({ user }) => {
  const params = useParams();
  const userId = params.userId;

  return (
    <div className="min-h-screen">
      <main className="flex flex-col items-center justify-center">
        <ProfilePage user={user} userId={userId} />
      </main>
    </div>
  );
};

export default withAuth(UserProfile);