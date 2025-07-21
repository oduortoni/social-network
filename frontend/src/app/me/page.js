'use client';

import withAuth from '../../../lib/withAuth';
import NavBar from '../../../components/layout/NavBar';
import ChatInterface from '../../../components/chat/ChatInterface';

const Me = ({ user }) => {
  return (
    <div className="min-h-screen bg-gray-100">
      <NavBar avatar={user?.avatar} />

      <main className="container mx-auto p-6">
        {/* Profile Header */}
        <div className="mb-8">
          <div className="bg-white shadow-md rounded-2xl p-8 text-center">
            <h1 className="text-3xl font-bold text-gray-800 mb-2">
              Welcome back, <span className="text-indigo-600">{user?.name || user?.email}</span>!
            </h1>
            <p className="text-gray-600">Your personal dashboard and social hub</p>
          </div>
        </div>

        {/* Chat Interface */}
        <div className="mb-8">
          <h2 className="text-xl font-semibold mb-4">Messages</h2>
          <ChatInterface user={user} />
        </div>

        {/* Feature Roadmap */}
        <div className="mt-8">
          <h2 className="text-xl font-semibold mb-4">Development Roadmap</h2>
          <div className="bg-white p-6 rounded-lg shadow">
            <div className="space-y-4">
              <div className="flex items-center space-x-3">
                <span className="w-8 h-8 bg-gray-300 text-gray-600 rounded-full flex items-center justify-center text-sm font-bold">1</span>
                <span className="text-gray-700 font-medium">Real-time Messaging & Notifications</span>
                <span className="text-sm text-gray-500">(backend integration)</span>
              </div>
              <div className="flex items-center space-x-3">
                <span className="w-8 h-8 bg-gray-300 text-gray-600 rounded-full flex items-center justify-center text-sm font-bold">2</span>
                <span className="text-gray-700 font-medium">Friend System</span>
                <span className="text-sm text-gray-500">(connections, following)</span>
              </div>
              <div className="flex items-center space-x-3">
                <span className="w-8 h-8 bg-gray-300 text-gray-600 rounded-full flex items-center justify-center text-sm font-bold">3</span>
                <span className="text-gray-600">Group Management</span>
                <span className="text-sm text-gray-500">(Create/join groups, member management)</span>
              </div>
              <div className="flex items-center space-x-3">
                <span className="w-8 h-8 bg-gray-300 text-gray-600 rounded-full flex items-center justify-center text-sm font-bold">4</span>
                <span className="text-gray-600">Posts & Feed</span>
                <span className="text-sm text-gray-500">(Social media posts with real-time updates)</span>
              </div>
              <div className="flex items-center space-x-3">
                <span className="w-8 h-8 bg-gray-300 text-gray-600 rounded-full flex items-center justify-center text-sm font-bold">5</span>
                <span className="text-gray-600">Events System</span>
                <span className="text-sm text-gray-500">(Group events with RSVP functionality)</span>
              </div>
            </div>
          </div>
        </div>
      </main>
    </div>
  );
};

export default withAuth(Me);