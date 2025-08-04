import React, { useEffect, useState, useCallback } from 'react';
import ActivityItem from '../homepage/ActivityItem';
import { wsService } from '../../lib/websocket';
import { notificationService } from '../../lib/notificationService';
import { profileAPI } from '../../lib/api';

const ActivitySidebar = () => {
  const [activities, setActivities] = useState([]);

  const formatTimeSince = (timestamp) => {
    const now = Date.now();
    const diffMs = now - timestamp * 1000; // assuming timestamp is in seconds
    const diffSec = Math.floor(diffMs / 1000);
    const diffMin = Math.floor(diffSec / 60);
    const diffHr = Math.floor(diffMin / 60);
    const diffDay = Math.floor(diffHr / 24);

    if (diffSec < 60) return 'Just now';
    if (diffMin < 60) return `${diffMin} min ago`;
    if (diffHr < 24) return `${diffHr} hr${diffHr > 1 ? 's' : ''} ago`;
    return `${diffDay} day${diffDay > 1 ? 's' : ''} ago`;
  };

  const handleRequest = useCallback((notification) => {
    console.log('Received follow_request notification:', notification);

    const {
      user_name,
      avatar,
      timestamp,
      user_id,
    } = notification;

    const activity = {
      image: profileAPI.fetchProfileImage(avatar || ''),
      name: user_name || 'Unknown User',
      action: 'sent a follow request',
      time: formatTimeSince(timestamp),
      isGroup: false,
      isPartial: false,
      userId: user_id,
    };

    setActivities((prev) => [activity, ...prev]);
  }, []);

  const handleNotification = useCallback((notification) => {
    notificationService.handleNotification(notification);
  }, []);

  useEffect(() => {
    // WebSocket handlers
    wsService.onMessage('notification', handleNotification);
    notificationService.onNotification('follow_request', handleRequest);

    // Initial pending follow requests
    profileAPI.fetchPendingFollowRequests()
      .then((requests) => {

        console.log(requests)
        if (!Array.isArray(requests)) return;

        const formatted = requests.map((req) => ({
          image: profileAPI.fetchProfileImage(req.avatar || ''),
          name: req.user_name || 'Unknown User',
          action: 'sent a follow request',
          time: formatTimeSince(req.timestamp),
          isGroup: false,
          isPartial: false,
          userId: req.user_id,
        }));

        setActivities((prev) => [...formatted, ...prev]);
      })
      .catch((err) => {
        console.error('Failed to load pending follow requests:', err);
      });

    return () => {
      wsService.removeHandler('notification', handleNotification);
    };
  }, [handleRequest, handleNotification]);

  return (
    <div className="w-72">
      <div className="rounded-xl p-4" style={{ backgroundColor: 'var(--primary-background)' }}>
        <h3 className="font-bold mb-4">Recent activity</h3>
        <div className="flex flex-col gap-4 cursor-pointer">
          {activities.length === 0 ? (
            <p className="text-sm text-gray-500">No recent activity</p>
          ) : (
            activities.map((activity, idx) => (
              <ActivityItem key={idx} {...activity} />
            ))
          )}
        </div>
      </div>
    </div>
  );
};

export default ActivitySidebar;
