'use client';

import withAuth from '../../../lib/withAuth';
import MainHomepage from '../../../components/homepage/MainHomepage';

const formatDate = (dateString) => {
  const date = new Date(dateString);
  const options = { year: 'numeric', month: 'long', day: 'numeric' };
  return date.toLocaleDateString('en-US', options);
};
const Me = ({ user }) => {
  console.log(user);
  return (
    <div className="min-h-screen">
      <main className="flex flex-col items-center justify-center p-6">
        <MainHomepage user={user} />
      </main>
    </div>
  );
};

export default withAuth(Me);