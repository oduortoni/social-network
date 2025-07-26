'use client';

import { useState, useEffect } from 'react';
import withAuth from '../../../lib/withAuth';
import MainHomepage from '../../../components/homepage/MainHomepage';
import { wsService } from '../../../lib/websocket';

const Me = ({ user }) => {
  const [connectionStatus, setConnectionStatus] = useState('disconnected');

  useEffect(() => {
    let mounted = true;

    // Set up connection status tracking
    wsService.onMessage('connection_status', (message) => {
      console.log('Connection status update:', message.status);
      if (mounted) {
        setConnectionStatus(message.status);
      }
    });

    // Set up error handling
    wsService.onMessage('error', (error) => {
      console.error('WebSocket error:', error);
      if (mounted) {
        setConnectionStatus('disconnected');
      }
    });

    // Check if user is authenticated before connecting
    const checkAuthAndConnect = async () => {
      try {
        // Debug: Check available cookies
        console.log('Available cookies:', document.cookie);

        // Check if we have a valid session by making a test API call
        const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/me`, {
          credentials: 'include'
        });

        console.log('Authentication check response:', response.status);

        if (response.ok) {
          console.log('User authenticated, attempting WebSocket connection...');
          setConnectionStatus('connecting');
          wsService.connect();
        } else {
          console.error('User not authenticated, cannot connect to WebSocket. Status:', response.status);
          setConnectionStatus('disconnected');
        }
      } catch (error) {
        console.error('Authentication check failed:', error);
        setConnectionStatus('disconnected');
      }
    };

    // Delay connection to ensure component is mounted and server is ready
    const connectTimer = setTimeout(() => {
      if (mounted) {
        checkAuthAndConnect();
      }
    }, 500); // 500ms delay

    return () => {
      mounted = false;
      clearTimeout(connectTimer);
      // Don't disconnect on page leave - only on logout or error
      // wsService.disconnect();
    };
  }, []);

  console.log(user);
  return (
    <div className="min-h-screen">
      <main className="flex flex-col items-center justify-center p-6">
        <MainHomepage user={user} connectionStatus={connectionStatus} />
      </main>
    </div>
  );
};

export default withAuth(Me);