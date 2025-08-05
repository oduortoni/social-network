'use client';

import { useState, useEffect } from 'react';

const ClientDate = ({ dateString, className, style, format = 'relative' }) => {
  const [formattedDate, setFormattedDate] = useState('');
  const [mounted, setMounted] = useState(false);

  useEffect(() => {
    setMounted(true);
    
    const formatDate = (dateString) => {
      try {
        const date = new Date(dateString);
        
        // Check if date is valid
        if (isNaN(date.getTime())) {
          return 'Invalid date';
        }

        if (format === 'time') {
          return date.toLocaleTimeString();
        }
        
        if (format === 'date') {
          return date.toLocaleDateString();
        }
        
        // Default relative format
        const now = new Date();
        const diffInSeconds = Math.floor((now - date) / 1000);

        if (diffInSeconds < 60) {
          return 'Just now';
        } else if (diffInSeconds < 3600) {
          const minutes = Math.floor(diffInSeconds / 60);
          return `${minutes}m ago`;
        } else if (diffInSeconds < 86400) {
          const hours = Math.floor(diffInSeconds / 3600);
          return `${hours}h ago`;
        } else if (diffInSeconds < 604800) {
          const days = Math.floor(diffInSeconds / 86400);
          return `${days}d ago`;
        } else {
          return date.toLocaleDateString();
        }
      } catch (error) {
        console.error('Error formatting date:', error);
        return 'Invalid date';
      }
    };

    setFormattedDate(formatDate(dateString));
  }, [dateString, format]);

  if (!mounted) {
    return <span className={className} style={style}>...</span>;
  }

  return <span className={className} style={style}>{formattedDate}</span>;
};

export default ClientDate;