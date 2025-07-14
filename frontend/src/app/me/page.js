"use client";

import withAuth from '../../../lib/withAuth';

const Me = ({ user }) => {
  return (
    <div>
      <h1>Welcome to your profile, {user?.name || user?.email}!</h1>
      {/* Add more profile-related components here */}
    </div>
  );
};

export default withAuth(Me);