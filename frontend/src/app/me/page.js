'use client';

import { useState, useEffect } from 'react';
import withAuth from '../../../lib/withAuth';
import MainHomepage from '../../../components/homepage/MainHomepage';
import { wsService } from '../../../lib/websocket';

const Me = ({ user }) => {
  const [connectionStatus, setConnectionStatus] = useState('disconnected');
  const [connectedUsers, setConnectedUsers] = useState([]);

  useEffect(() => {
    let mounted = true;

    // Load connected users from API
    const loadConnectedUsers = async () => {
      try {
        const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/api/users/online`, {
          credentials: 'include'
        });


        if (response.ok) {
          const data = await response.json();
          setConnectedUsers(data.online_users || []);
        } else {
          const errorText = await response.text();
          console.error('API error:', response.status, errorText);
        }
      } catch (error) {
        console.error('Failed to load connected users:', error);
      }
    };

    // Set up connection status tracking
    wsService.onMessage('connection_status', (message) => {
      if (mounted) {
        setConnectionStatus(message.status);

        // Load connected users immediately when connection is established
        if (message.status === 'connected') {
          loadConnectedUsers();
        }
      }
    });

    // Set up error handling
    wsService.onMessage('error', (error) => {
      console.error('WebSocket error:', error);
      if (mounted) {
        setConnectionStatus('disconnected');
      }
    });

    // Handle user connection/disconnection notifications
    wsService.onMessage('notification', (notification) => {
      if (mounted && notification.type === 'notification') {
        if (notification.subtype === 'user_connected') {
          setConnectedUsers(prev => {
            // Add user if not already in the list
            const userExists = prev.some(u => u.user_id === notification.user_id);
            if (!userExists) {
              return [...prev, {
                user_id: notification.user_id,
                nickname: notification.nickname, // Use consistent field name
                avatar: notification.avatar
              }];
            }
            return prev;
          });
        } else if (notification.subtype === 'user_disconnected') {
          setConnectedUsers(prev =>
            prev.filter(u => u.user_id !== notification.user_id)
          );
        }
      }
    });

    // Since withAuth already ensures authentication, directly connect to WebSocket
    const connectTimer = setTimeout(() => {
      if (mounted) {
        setConnectionStatus('connecting');
        wsService.connect();
      }
    }, 500); // 500ms delay

    return () => {
      mounted = false;
      clearTimeout(connectTimer);
    };
  }, []);

  return (
    <div className="min-h-screen">
      <main className="flex flex-col items-center justify-center p-6">
        <MainHomepage user={user} connectionStatus={connectionStatus} connectedUsers={connectedUsers} />
      </main>
    </div>
  );
};

export default withAuth(Me);