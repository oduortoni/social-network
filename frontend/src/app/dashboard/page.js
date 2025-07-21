'use client';

import withAuth from '../../../lib/withAuth';
import NavBar from '../../../components/layout/NavBar';
import ChatInterface from '../../../components/chat/ChatInterface';
import NotificationCenter from '../../../components/notifications/NotificationCenter';

const Dashboard = ({ user }) => {
  return (
    <div className="min-h-screen bg-gray-100">
      <NavBar avatar={user?.avatar} />
      
      <main className="container mx-auto p-6">
        <div className="mb-6">
          <h1 className="text-3xl font-bold text-gray-800">Dashboard</h1>
          <p className="text-gray-600">Welcome back, {user?.name || user?.email}!</p>
        </div>
        
        {/* Chat Interface */}
        <div className="mb-8">
          <h2 className="text-xl font-semibold mb-4">Messages</h2>
          <ChatInterface user={user} />
        </div>
        
        {/* Other dashboard content */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          <div className="bg-white p-6 rounded-lg shadow">
            <h3 className="font-semibold mb-2">Recent Activity</h3>
            <p className="text-gray-600">No recent activity</p>
          </div>
          
          <div className="bg-white p-6 rounded-lg shadow">
            <h3 className="font-semibold mb-2">Groups</h3>
            <p className="text-gray-600">No groups joined</p>
          </div>
          
          <div className="bg-white p-6 rounded-lg shadow">
            <h3 className="font-semibold mb-2">Friends</h3>
            <p className="text-gray-600">No friends added</p>
          </div>
        </div>
      </main>
    </div>
  );
};

export default withAuth(Dashboard);