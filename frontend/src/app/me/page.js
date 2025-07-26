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
        console.log('Connection status update:', message.status);
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

    // Notifications are now handled by the Header component



    // Since withAuth already ensures authentication, directly connect to WebSocket
    const connectTimer = setTimeout(() => {
      if (mounted) {
        console.log('User is authenticated (via withAuth), attempting WebSocket connection...');
        console.log('Available cookies:', document.cookie);
        setConnectionStatus('connecting');
        wsService.connect();
      }
    }, 500); // 500ms delay

    return () => {
      mounted = false;
      clearTimeout(connectTimer);
    };
  }, []);

  // console.log(user);
  return (
    <div className="min-h-screen">
      <main className="flex flex-col items-center justify-center p-6">
        <MainHomepage user={user} connectionStatus={connectionStatus} />
      </main>
    </div>
  );
};

export default withAuth(Me);