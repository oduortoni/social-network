import React, { useEffect, useState } from 'react';
import ActivityItem from '../homepage/ActivityItem';

const dummyActivities = [
  {
    image: 'https://randomuser.me/api/portraits/men/72.jpg',
    name: 'Vitaliy',
    action: 'sent a friend request',
    time: '3 min ago',
    isGroup: false,
    isPartial: false
  },
  {
    image: 'https://randomuser.me/api/portraits/men/62.jpg',
    name: 'Maksym',
    action: 'sent a friend request',
    time: '4 hrs ago',
    isGroup: false,
    isPartial: false
  },
  {
    image: 'https://randomuser.me/api/portraits/men/52.jpg',
    name: 'Evgeniy',
    action: 'sent a friend request',
    time: '7 hrs ago',
    isGroup: false,
    isPartial: false
  },
  {
    image: 'https://randomuser.me/api/portraits/women/32.jpg',
    name: 'Rosaline',
    action: 'sent a friend request',
    time: '1 hr ago',
    isGroup: false,
    isPartial: false
  },
  {
    image: 'https://randomuser.me/api/portraits/women/42.jpg',
    name: 'UX designers group',
    action: '',
    time: '12 hrs ago',
    isGroup: true,
    isPartial: false
  }
];

const ActivitySidebar = () => {
  const [activities, setActivities] = useState(dummyActivities);

  useEffect(() => {
    // Simulate fetching from backend
    // to be replaced with the real fetch logic
    // fetch('/api/activities').then(...)
    const timer = setTimeout(() => {
      // Example: setActivities(fetchedData)
      // I am using dummy data for now
    }, 1000);
    return () => clearTimeout(timer);
  }, []);

  return (
    <div className="w-72">
      <div className="bg-[#101b70] rounded-xl p-4">
        <h3 className="font-bold mb-4">Recent activity</h3>
        <div className="flex flex-col gap-4">
          {activities.map((activity, idx) => (
            <ActivityItem key={idx} {...activity} />
          ))}
        </div>
      </div>
    </div>
  );
};
export default ActivitySidebar;