'use client';

import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import withAuth from '../../../lib/withAuth';
import Header from '../../../components/layout/Header';
import { chatAPI, profileAPI } from '../../../lib/api';
import { MessageCircleIcon, ArrowLeftIcon } from 'lucide-react';

const ChatsPage = ({ user }) => {
  const [messageableUsers, setMessageableUsers] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const router = useRouter();

  useEffect(() => {
    loadMessageableUsers();
  }, []);

  const loadMessageableUsers = async () => {
    try {
      setLoading(true);
      const users = await chatAPI.getMessageableUsers();
      setMessageableUsers(users || []);
    } catch (error) {
      console.error('Failed to load messageable users:', error);
      setError('Failed to load users. Please try again.');
    } finally {
      setLoading(false);
    }
  };

  const handleUserSelect = (selectedUser) => {
    // Navigate to chat interface with the selected user
    router.push(`/chats/${selectedUser.id}?nickname=${encodeURIComponent(selectedUser.nickname)}`);
  };

  const handleBackToHome = () => {
    router.push('/me');
  };

  if (loading) {
    return (
      <div className="min-h-screen" style={{ backgroundColor: 'var(--primary-background)' }}>
        <Header user={user} />
        <div className="flex items-center justify-center h-96">
          <div className="text-center">
            <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-500 mx-auto mb-4"></div>
            <p style={{ color: 'var(--primary-text)' }}>Loading users...</p>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen" style={{ backgroundColor: 'var(--primary-background)' }}>
      <Header user={user} />
      
      <div className="max-w-4xl mx-auto p-6">
        {/* Header with back button */}
        <div className="flex items-center mb-6">
          <button
            onClick={handleBackToHome}
            className="flex items-center gap-2 px-4 py-2 rounded-lg hover:bg-gray-100 transition-colors"
            style={{ color: 'var(--primary-text)' }}
          >
            <ArrowLeftIcon className="w-5 h-5" />
            Back to Home
          </button>
        </div>

        {/* Page Title */}
        <div className="flex items-center gap-3 mb-8">
          <MessageCircleIcon className="w-8 h-8" style={{ color: 'var(--primary-accent)' }} />
          <h1 className="text-3xl font-bold" style={{ color: 'var(--primary-text)' }}>
            Start a Conversation
          </h1>
        </div>

        {/* Error State */}
        {error && (
          <div className="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded mb-6">
            {error}
            <button
              onClick={loadMessageableUsers}
              className="ml-4 underline hover:no-underline"
            >
              Try again
            </button>
          </div>
        )}

        {/* Users List */}
        {messageableUsers.length > 0 ? (
          <div className="grid gap-4">
            <p className="text-sm mb-4" style={{ color: 'var(--secondary-text)' }}>
              Select a user to start chatting. You can message mutual followers and users with public profiles.
            </p>
            
            {messageableUsers.map((messageableUser) => (
              <div
                key={messageableUser.id}
                onClick={() => handleUserSelect(messageableUser)}
                className="flex items-center gap-4 p-4 rounded-lg border cursor-pointer hover:shadow-md transition-all duration-200 hover:scale-[1.02]"
                style={{ 
                  backgroundColor: 'var(--secondary-background)',
                  borderColor: 'var(--border-color)',
                }}
              >
                {/* User Avatar */}
                <img
                  src={profileAPI.fetchProfileImage(messageableUser.avatar || '')}
                  alt={`${messageableUser.nickname}'s avatar`}
                  className="w-12 h-12 rounded-full object-cover"
                  onError={(e) => {
                    e.target.src = "data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iNDAiIGhlaWdodD0iNDAiIHZpZXdCb3g9IjAgMCA0MCA0MCIgZmlsbD0ibm9uZSIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj4KPGNpcmNsZSBjeD0iMjAiIGN5PSIyMCIgcj0iMjAiIGZpbGw9IiNGM0Y0RjYiLz4KPGNpcmNsZSBjeD0iMjAiIGN5PSIxNiIgcj0iNiIgZmlsbD0iIzlDQTNBRiIvPgo8cGF0aCBkPSJNMzIgMzJDMzIgMjYuNDc3MiAyNy41MjI4IDIyIDIyIDIySDE4QzEyLjQ3NzIgMjIgOCAyNi40NzcyIDggMzJWMzJIMzJWMzJaIiBmaWxsPSIjOUNBM0FGIi8+Cjwvc3ZnPgo=";
                  }}
                />
                
                {/* User Info */}
                <div className="flex-1">
                  <h3 className="font-semibold text-lg" style={{ color: 'var(--primary-text)' }}>
                    {messageableUser.nickname}
                  </h3>
                  <p className="text-sm" style={{ color: 'var(--secondary-text)' }}>
                    Click to start chatting
                  </p>
                </div>

                {/* Chat Icon */}
                <MessageCircleIcon 
                  className="w-6 h-6" 
                  style={{ color: 'var(--primary-accent)' }} 
                />
              </div>
            ))}
          </div>
        ) : (
          !loading && (
            <div className="text-center py-12">
              <MessageCircleIcon 
                className="w-16 h-16 mx-auto mb-4 opacity-50" 
                style={{ color: 'var(--secondary-text)' }} 
              />
              <h3 className="text-xl font-semibold mb-2" style={{ color: 'var(--primary-text)' }}>
                No users available to message
              </h3>
              <p style={{ color: 'var(--secondary-text)' }}>
                You can message mutual followers and users with public profiles.
                <br />
                Try following some users or making your profile public to expand your network.
              </p>
            </div>
          )
        )}
      </div>
    </div>
  );
};

export default withAuth(ChatsPage);
