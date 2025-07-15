"use client";

import withAuth from '../../../lib/withAuth';
import NavBar from '../../../components/layout/NavBar';

const Me = ({ user }) => {
  return (
    <div className="min-h-screen bg-gray-100">
      <NavBar avatar={user?.avatar} />
      <main className="flex flex-col items-center justify-center p-6">
        <div className="bg-white shadow-md rounded-2xl p-8 w-full max-w-md text-center mt-8">
          <h1 className="text-2xl font-semibold text-gray-800 mb-4">
            Welcome to your profile,{" "}
            <span className="text-indigo-600">{user?.name || user?.email}</span>!
          </h1>
          <p className="text-gray-500">This is your personal dashboard.</p>
        </div>
      </main>
    </div>
  );
};

export default withAuth(Me);