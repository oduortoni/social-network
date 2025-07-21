'use client';

import { useState, useEffect, useRef } from 'react';
import { wsService } from '../../lib/websocket';
import { chatAPI } from '../../lib/api';

const ChatInterface = ({ user }) => {
  const [messages, setMessages] = useState([]);
  const [newMessage, setNewMessage] = useState('');
  const [activeChat, setActiveChat] = useState(null); // { type: 'private', id: userId } or { type: 'group', id: groupId }
  const [connectionStatus, setConnectionStatus] = useState('disconnected'); // 'connecting', 'connected', 'disconnected'
  const messagesEndRef = useRef(null);

  useEffect(() => {
    let mounted = true;

    // Set up message handlers first
    wsService.onMessage('private', handlePrivateMessage);
    wsService.onMessage('group', handleGroupMessage);
    wsService.onMessage('broadcast', handleBroadcastMessage);
    wsService.onMessage('notification', handleNotification);

    // Add connection status tracking
    wsService.onMessage('connection_status', (message) => {
      if (mounted) {
        setConnectionStatus(message.status);
      }
    });

    // Delay connection to ensure component is mounted and server is ready
    const connectTimer = setTimeout(() => {
      if (mounted) {
        setConnectionStatus('connecting');
        wsService.connect();
      }
    }, 500); // 500ms delay

    return () => {
      mounted = false;
      clearTimeout(connectTimer);
      wsService.disconnect();
    };
  }, []);

  useEffect(() => {
    scrollToBottom();
  }, [messages]);

  const handlePrivateMessage = (message) => {
    if (activeChat?.type === 'private' && 
        (message.to === user.id || message.from === activeChat.id)) {
      setMessages(prev => [...prev, message]);
    }
  };

  const handleGroupMessage = (message) => {
    if (activeChat?.type === 'group' && message.group_id === activeChat.id) {
      setMessages(prev => [...prev, message]);
    }
  };

  const handleBroadcastMessage = (message) => {
    // Show broadcast messages in all chats
    setMessages(prev => [...prev, { ...message, isBroadcast: true }]);
  };

  const handleNotification = (notification) => {
    // Handle real-time notifications
    console.log('New notification:', notification);
    // You could show a toast notification here
  };

  const sendMessage = () => {
    if (!newMessage.trim() || !activeChat) return;

    // Check if WebSocket is connected
    if (!wsService.isConnected()) {
      alert('Not connected to chat server. Please wait for connection to be established.');
      return;
    }

    if (activeChat.type === 'private') {
      wsService.sendMessage('private', newMessage, activeChat.id);
    } else if (activeChat.type === 'group') {
      wsService.sendMessage('group', newMessage, null, activeChat.id);
    }

    setNewMessage('');
  };

  const loadChatHistory = async (chatType, chatId) => {
    try {
      let history;
      if (chatType === 'private') {
        history = await chatAPI.getPrivateMessages(chatId);
      } else if (chatType === 'group') {
        history = await chatAPI.getGroupMessages(chatId);
      }
      setMessages(history || []);
      setActiveChat({ type: chatType, id: chatId });
    } catch (error) {
      console.error('Failed to load chat history:', error);
    }
  };

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  };

  return (
    <div className="flex h-96 bg-white rounded-lg shadow-lg">
      {/* Chat Sidebar */}
      <div className="w-1/3 border-r border-gray-200 p-4">
        <div className="flex items-center justify-between mb-4">
          <h3 className="font-semibold">Chats</h3>
          <div className={`w-3 h-3 rounded-full ${
            connectionStatus === 'connected' ? 'bg-green-500' :
            connectionStatus === 'connecting' ? 'bg-yellow-500' : 'bg-red-500'
          }`} title={`WebSocket ${connectionStatus}`}></div>
        </div>
        
        {/* Private Chats */}
        <div className="mb-4">
          <h4 className="text-sm text-gray-600 mb-2">Private Messages</h4>
          <button
            onClick={() => loadChatHistory('private', 1)} // Example user ID
            className="w-full text-left p-2 hover:bg-gray-100 rounded"
          >
            User 1
          </button>
        </div>
        
        {/* Group Chats */}
        <div>
          <h4 className="text-sm text-gray-600 mb-2">Groups</h4>
          <button
            onClick={() => loadChatHistory('group', 1)} // Example group ID
            className="w-full text-left p-2 hover:bg-gray-100 rounded"
          >
            Test Group
          </button>
        </div>
      </div>
      
      {/* Chat Messages */}
      <div className="flex-1 flex flex-col">
        {activeChat ? (
          <>
            {/* Messages Area */}
            <div className="flex-1 p-4 overflow-y-auto">
              {messages.map((message, index) => (
                <div
                  key={index}
                  className={`mb-2 ${
                    message.isBroadcast ? 'text-center text-blue-600' : ''
                  }`}
                >
                  <div className="text-sm text-gray-500">
                    {new Date(message.timestamp * 1000).toLocaleTimeString()}
                  </div>
                  <div className="bg-gray-100 p-2 rounded">
                    {message.content}
                  </div>
                </div>
              ))}
              <div ref={messagesEndRef} />
            </div>
            
            {/* Message Input */}
            <div className="p-4 border-t border-gray-200">
              {connectionStatus !== 'connected' && (
                <div className="mb-2 text-sm text-center">
                  <span className={`${
                    connectionStatus === 'connecting' ? 'text-yellow-600' : 'text-red-600'
                  }`}>
                    {connectionStatus === 'connecting' ? 'Connecting to chat...' : 'Disconnected from chat server'}
                  </span>
                </div>
              )}
              <div className="flex">
                <input
                  type="text"
                  value={newMessage}
                  onChange={(e) => setNewMessage(e.target.value)}
                  onKeyDown={(e) => e.key === 'Enter' && sendMessage()}
                  placeholder={connectionStatus === 'connected' ? "Type a message..." : "Waiting for connection..."}
                  disabled={connectionStatus !== 'connected'}
                  className={`flex-1 p-2 border border-gray-300 rounded-l ${
                    connectionStatus !== 'connected' ? 'bg-gray-100 cursor-not-allowed' : ''
                  }`}
                />
                <button
                  onClick={sendMessage}
                  disabled={connectionStatus !== 'connected'}
                  className={`px-4 py-2 text-white rounded-r ${
                    connectionStatus === 'connected'
                      ? 'bg-blue-500 hover:bg-blue-600'
                      : 'bg-gray-400 cursor-not-allowed'
                  }`}
                >
                  Send
                </button>
              </div>
            </div>
          </>
        ) : (
          <div className="flex-1 flex items-center justify-center text-gray-500">
            Select a chat to start messaging
          </div>
        )}
      </div>
    </div>
  );
};

export default ChatInterface;