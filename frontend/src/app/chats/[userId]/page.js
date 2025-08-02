'use client';

import { useState, useEffect } from 'react';
import { useRouter, useParams, useSearchParams } from 'next/navigation';
import withAuth from '../../../../lib/withAuth';
import Header from '../../../../components/layout/Header';
import ChatInterface from '../../../../components/chat/ChatInterface';
import { wsService } from '../../../../lib/websocket';
import { ArrowLeftIcon, MessageCircleIcon } from 'lucide-react';

const ChatPage = ({ user }) => {
  const [connectionStatus, setConnectionStatus] = useState('disconnected');
  const [targetUser, setTargetUser] = useState(null);
  const router = useRouter();
  const params = useParams();
  const searchParams = useSearchParams();

  const userId = params.userId;
  const nickname = searchParams.get('nickname');

  useEffect(() => {
    // Set target user info
    if (userId && nickname) {
      setTargetUser({
        id: parseInt(userId),
        nickname: decodeURIComponent(nickname)
      });
    }

    let mounted = true;

    // Set up WebSocket connection status tracking
    wsService.onMessage('connection_status', (message) => {
      if (mounted) {
        setConnectionStatus(message.status);
      }
    });

    // Connect to WebSocket if not already connected
    if (!wsService.isConnected()) {
      setConnectionStatus('connecting');
      wsService.connect();
    } else {
      setConnectionStatus('connected');
    }

    return () => {
      mounted = false;
    };
  }, [userId, nickname]);

  const handleBackToChats = () => {
    router.push('/chats');
  };

  const handleBackToHome = () => {
    router.push('/me');
  };

  if (!targetUser) {
    return (
      <div className="min-h-screen" style={{ backgroundColor: 'var(--primary-background)' }}>
        <Header user={user} />
        <div className="flex items-center justify-center h-96">
          <div className="text-center">
            <p style={{ color: 'var(--primary-text)' }}>Loading chat...</p>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen" style={{ backgroundColor: 'var(--primary-background)' }}>
      <Header user={user} />
      
      <div className="max-w-6xl mx-auto p-6">
        {/* Navigation Header */}
        <div className="flex items-center justify-between mb-6">
          <div className="flex items-center gap-4">
            <button
              onClick={handleBackToChats}
              className="flex items-center gap-2 px-3 py-2 rounded-lg hover:bg-gray-100 transition-colors"
              style={{ color: 'var(--primary-text)' }}
            >
              <ArrowLeftIcon className="w-4 h-4" />
              Back to Chats
            </button>
            
            <div className="h-6 w-px bg-gray-300"></div>
            
            <button
              onClick={handleBackToHome}
              className="text-sm px-3 py-2 rounded-lg hover:bg-gray-100 transition-colors"
              style={{ color: 'var(--secondary-text)' }}
            >
              Home
            </button>
          </div>

          {/* Connection Status */}
          <div className="flex items-center gap-2">
            <div className={`w-3 h-3 rounded-full ${
              connectionStatus === 'connected' ? 'bg-green-500' :
              connectionStatus === 'connecting' ? 'bg-yellow-500' : 'bg-red-500'
            }`}></div>
            <span className="text-sm capitalize" style={{ color: 'var(--secondary-text)' }}>
              {connectionStatus}
            </span>
          </div>
        </div>

        {/* Chat Header */}
        <div className="flex items-center gap-3 mb-6 p-4 rounded-lg" style={{ backgroundColor: 'var(--secondary-background)' }}>
          <MessageCircleIcon className="w-6 h-6" style={{ color: 'var(--primary-accent)' }} />
          <div>
            <h1 className="text-xl font-semibold" style={{ color: 'var(--primary-text)' }}>
              Chat with {targetUser.nickname}
            </h1>
            <p className="text-sm" style={{ color: 'var(--secondary-text)' }}>
              {connectionStatus === 'connected' 
                ? 'Connected - messages will be delivered in real-time' 
                : connectionStatus === 'connecting'
                ? 'Connecting to chat server...'
                : 'Disconnected - messages may not be delivered'
              }
            </p>
          </div>
        </div>

        {/* Chat Interface */}
        <div className="bg-white rounded-lg shadow-lg overflow-hidden">
          <ChatInterface
            user={user}
            connectionStatus={connectionStatus}
            initialChat={{ type: 'private', id: targetUser.id, nickname: targetUser.nickname }}
            showSidebar={false}
          />
        </div>

        {/* Help Text */}
        <div className="mt-4 text-center">
          <p className="text-sm" style={{ color: 'var(--secondary-text)' }}>
            Messages are automatically saved and will be delivered via WebSocket if the recipient is online.
          </p>
        </div>
      </div>
    </div>
  );
};

export default withAuth(ChatPage);
