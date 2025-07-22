"use client";

import withAuth from '../../../lib/withAuth';
import MainHomepage from '../../../components/homepage/MainHomepage';

const Me = ({ user }) => {
  return (
    <div className="min-h-screen">
      <main className="flex flex-col items-center justify-center p-6">
        <MainHomepage user={user} />
      </main>
    </div>
  );
};

export default withAuth(Me);