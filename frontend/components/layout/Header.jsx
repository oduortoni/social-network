import React, { useState, useRef, useEffect } from 'react';
import { HomeIcon, BellIcon, UsersIcon, MessageCircleIcon, SearchIcon, ChevronDownIcon, LogOutIcon } from 'lucide-react';
import { handleLogout } from '../../lib/auth';

const Header = () => {
  const [dropdownOpen, setDropdownOpen] = useState(false);
  const profileRef = useRef(null);

  // Close dropdown when clicking outside
  useEffect(() => {
    function handleClickOutside(event) {
      if (profileRef.current && !profileRef.current.contains(event.target)) {
        setDropdownOpen(false);
      }
    }
    document.addEventListener("mousedown", handleClickOutside);
    return () => document.removeEventListener("mousedown", handleClickOutside);
  }, []);

  return (
    <div className="flex items-center justify-between w-full px-4 py-2">
      
      {/* Logo and Search Bar */}
      <div className="flex items-center gap-4">
        <div className="bg-white rounded-full p-1 cursor-pointer">
          <div className="w-6 h-6 rounded-full" style={{ backgroundColor: 'var(--primary-background)' }}></div>
        </div>
        <div className="relative">
          <SearchIcon className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 h-4 w-4" />
          <input
            type="text"
            placeholder="Search for people, communities ..."
            className="rounded-full py-1.5 pl-10 pr-4 text-sm focus:outline-none"
            style={{ backgroundColor: 'var(--secondary-background)', color: 'var(--primary-text)' }}
          />
        </div>
      </div>
      
      {/* Navigation Icons */}
      <div className="flex items-center gap-6">
        {/* Home Icon */}
        <div className="flex flex-col items-center cursor-pointer" onClick={() => { /* TODO: handle Home click */ }}>
          <HomeIcon className="w-6 h-6" style={{ color: 'var(--primary-accent)' }} />
          <span className="text-xs" style={{ color: 'var(--primary-text)' }}>Home</span>
        </div>
        {/* Notifications Icon */}
        <div className="flex flex-col items-center cursor-pointer" onClick={() => { /* TODO: handle Notifications click */ }}>
          <BellIcon className="w-6 h-6" style={{ color: 'var(--primary-text)' }} />
          <span className="text-xs" style={{ color: 'var(--primary-text)' }}>Notifications</span>
        </div>
        
        {/* Groups Icon */}
        <div className="flex flex-col items-center cursor-pointer" onClick={() => { /* TODO: handle Groups click */ }}>
          <UsersIcon className="w-6 h-6" style={{ color: 'var(--primary-text)' }} />
          <span className="text-xs" style={{ color: 'var(--primary-text)' }}>Groups</span>
        </div>
        
        {/* Chats Icon */}
        <div className="flex flex-col items-center cursor-pointer" onClick={() => { /* TODO: handle Chats click */ }}>
          <MessageCircleIcon className="w-6 h-6" style={{ color: 'var(--primary-text)' }} />
          <span className="text-xs" style={{ color: 'var(--primary-text)' }}>Chats</span>
        </div>
      </div>
      
      {/* Profile Dropdown */}
      <div className="relative" ref={profileRef}>
        <button
          className="flex items-center gap-2 focus:outline-none"
          onClick={() => setDropdownOpen((open) => !open)}
        >
          <img src="https://randomuser.me/api/portraits/men/30.jpg" alt="Profile" className="w-8 h-8 rounded-full" />
          <span className="text-sm font-medium">Zeus</span>
          <ChevronDownIcon className="w-4 h-4" />
        </button>
        {dropdownOpen && (
          <div className="absolute right-0 mt-2 w-40 bg-white rounded shadow-lg z-10">
            <button
              className="flex items-center gap-2 px-4 py-2 w-full text-left bg-white rounded transition-colors cursor-pointer hover:bg-blue-100"
              style={{ color: 'var(--primary-background)' }}
              onClick={handleLogout}
            >
              <LogOutIcon className="w-4 h-4" />
              Logout
            </button>
          </div>
        )}
      </div>
    </div>
  );
};

export default Header;