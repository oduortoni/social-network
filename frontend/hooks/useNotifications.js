import { useState, useEffect } from 'react';
import { wsService } from '../lib/websocket';

// Simple notification hook - just tracks unread count
export const useSimpleNotifications = () => {
  const [unreadCount, setUnreadCount] = useState(0);
  const [notifications, setNotifications] = useState([]);

  useEffect(() => {
    // Handle incoming notifications
    const handleNotification = (notification) => {

      // Add to notifications list
      setNotifications(prev => [notification, ...prev.slice(0, 19)]); // Keep last 20

      // Increment unread count
      setUnreadCount(prev => prev + 1);
    };

    // Listen for WebSocket notifications
    wsService.onMessage('notification', handleNotification);

    // No cleanup needed since connection is managed elsewhere
  }, []);

  const markAllAsRead = () => {
    setUnreadCount(0);
  };

  const clearNotifications = () => {
    setNotifications([]);
    setUnreadCount(0);
  };

  return {
    unreadCount,
    notifications,
    markAllAsRead,
    clearNotifications
  };
};

export default useSimpleNotifications;
