# Frontend WebSocket Testing Guide

## Overview

This guide covers testing strategies for the frontend WebSocket implementation, including unit tests, integration tests, and end-to-end testing for React components with real-time functionality. Feel free to create the files needed to test these features as highlighted within the document

## Testing Architecture

### Test Structure
```
frontend/
├── __tests__/
│   ├── components/
│   │   ├── chat/
│   │   │   ├── ChatInterface.test.jsx
│   │   │   └── MessageInput.test.jsx
│   │   └── notifications/
│   │       └── NotificationCenter.test.jsx
│   ├── lib/
│   │   ├── websocket.test.js
│   │   └── api.test.js
│   └── __mocks__/
│       ├── websocket.js
│       └── api.js
```

## Unit Testing

### WebSocket Service Testing

```javascript
// __tests__/lib/websocket.test.js
import { WebSocketService } from '../../lib/websocket';

// Mock WebSocket
global.WebSocket = jest.fn(() => ({
  close: jest.fn(),
  send: jest.fn(),
  addEventListener: jest.fn(),
  removeEventListener: jest.fn(),
  readyState: WebSocket.OPEN
}));

describe('WebSocketService', () => {
  let wsService;

  beforeEach(() => {
    wsService = new WebSocketService();
    jest.clearAllMocks();
  });

  test('should connect to WebSocket', () => {
    wsService.connect();
    expect(global.WebSocket).toHaveBeenCalledWith(
      expect.stringContaining('/ws')
    );
  });

  test('should register message handlers', () => {
    const handler = jest.fn();
    wsService.onMessage('test_type', handler);
    
    // Simulate message
    wsService.handleMessage({ type: 'test_type', data: 'test' });
    expect(handler).toHaveBeenCalledWith({ type: 'test_type', data: 'test' });
  });

  test('should send messages when connected', () => {
    wsService.ws = { 
      readyState: WebSocket.OPEN, 
      send: jest.fn() 
    };
    
    wsService.sendMessage('private', 'Hello', 123);
    
    expect(wsService.ws.send).toHaveBeenCalledWith(
      JSON.stringify({
        type: 'private',
        content: 'Hello',
        to: 123,
        timestamp: expect.any(Number)
      })
    );
  });

  test('should not send messages when disconnected', () => {
    const consoleSpy = jest.spyOn(console, 'warn').mockImplementation();
    wsService.ws = { readyState: WebSocket.CLOSED };
    
    wsService.sendMessage('private', 'Hello', 123);
    
    expect(consoleSpy).toHaveBeenCalledWith(
      expect.stringContaining('WebSocket not connected')
    );
    consoleSpy.mockRestore();
  });

  test('should handle reconnection', () => {
    jest.useFakeTimers();
    wsService.reconnect();
    
    expect(wsService.reconnectAttempts).toBe(1);
    
    jest.advanceTimersByTime(2000);
    expect(global.WebSocket).toHaveBeenCalled();
    
    jest.useRealTimers();
  });
});
```

### API Service Testing

```javascript
// __tests__/lib/api.test.js
import { chatAPI } from '../../lib/api';

// Mock fetch
global.fetch = jest.fn();

describe('chatAPI', () => {
  beforeEach(() => {
    fetch.mockClear();
  });

  test('should fetch private messages', async () => {
    const mockMessages = [
      { id: 1, content: 'Hello', sender_id: 1, receiver_id: 2 }
    ];
    
    fetch.mockResolvedValueOnce({
      ok: true,
      json: async () => mockMessages
    });

    const result = await chatAPI.getPrivateMessages(2, 10, 0);
    
    expect(fetch).toHaveBeenCalledWith(
      expect.stringContaining('/api/messages/private?user=2&limit=10&offset=0'),
      expect.objectContaining({
        credentials: 'include'
      })
    );
    expect(result).toEqual(mockMessages);
  });

  test('should handle API errors', async () => {
    fetch.mockResolvedValueOnce({
      ok: false,
      status: 500
    });

    await expect(chatAPI.getPrivateMessages(2)).rejects.toThrow('API Error: 500');
  });
});
```

## Component Testing

### ChatInterface Component Testing

```javascript
// __tests__/components/chat/ChatInterface.test.jsx
import React from 'react';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import ChatInterface from '../../../components/chat/ChatInterface';
import { wsService } from '../../../lib/websocket';

// Mock the WebSocket service
jest.mock('../../../lib/websocket', () => ({
  wsService: {
    connect: jest.fn(),
    disconnect: jest.fn(),
    onMessage: jest.fn(),
    sendMessage: jest.fn(),
    isConnected: jest.fn(() => true)
  }
}));

// Mock the API
jest.mock('../../../lib/api', () => ({
  chatAPI: {
    getPrivateMessages: jest.fn(() => Promise.resolve([])),
    getGroupMessages: jest.fn(() => Promise.resolve([]))
  }
}));

describe('ChatInterface', () => {
  const mockUser = { id: 1, name: 'Test User' };

  beforeEach(() => {
    jest.clearAllMocks();
  });

  test('should render chat interface', () => {
    render(<ChatInterface user={mockUser} />);
    
    expect(screen.getByText('Chats')).toBeInTheDocument();
    expect(screen.getByPlaceholderText('Type a message...')).toBeInTheDocument();
  });

  test('should connect to WebSocket on mount', () => {
    render(<ChatInterface user={mockUser} />);
    
    expect(wsService.onMessage).toHaveBeenCalledWith('private', expect.any(Function));
    expect(wsService.onMessage).toHaveBeenCalledWith('group', expect.any(Function));
    expect(wsService.onMessage).toHaveBeenCalledWith('broadcast', expect.any(Function));
    expect(wsService.onMessage).toHaveBeenCalledWith('notification', expect.any(Function));
  });

  test('should send message when form is submitted', () => {
    render(<ChatInterface user={mockUser} />);
    
    // Set up active chat
    const privateButton = screen.getByText('User 1');
    fireEvent.click(privateButton);
    
    // Type message
    const input = screen.getByPlaceholderText('Type a message...');
    fireEvent.change(input, { target: { value: 'Hello test' } });
    
    // Send message
    const sendButton = screen.getByText('Send');
    fireEvent.click(sendButton);
    
    expect(wsService.sendMessage).toHaveBeenCalledWith('private', 'Hello test', 1);
    expect(input.value).toBe('');
  });

  test('should display connection status', async () => {
    render(<ChatInterface user={mockUser} />);
    
    // Simulate connection status change
    const statusHandler = wsService.onMessage.mock.calls.find(
      call => call[0] === 'connection_status'
    )[1];
    
    statusHandler({ status: 'connecting' });
    
    await waitFor(() => {
      expect(screen.getByText('Connecting to chat...')).toBeInTheDocument();
    });
  });

  test('should disable input when disconnected', async () => {
    wsService.isConnected.mockReturnValue(false);
    
    render(<ChatInterface user={mockUser} />);
    
    const input = screen.getByPlaceholderText('Waiting for connection...');
    expect(input).toBeDisabled();
  });
});
```

### NotificationCenter Component Testing

```javascript
// __tests__/components/notifications/NotificationCenter.test.jsx
import React from 'react';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import NotificationCenter from '../../../components/notifications/NotificationCenter';
import { wsService } from '../../../lib/websocket';
import { chatAPI } from '../../../lib/api';

jest.mock('../../../lib/websocket');
jest.mock('../../../lib/api');

describe('NotificationCenter', () => {
  beforeEach(() => {
    jest.clearAllMocks();
    chatAPI.getNotifications.mockResolvedValue([]);
  });

  test('should render notification bell', () => {
    render(<NotificationCenter />);
    expect(screen.getByText('🔔')).toBeInTheDocument();
  });

  test('should show unread count badge', async () => {
    chatAPI.getNotifications.mockResolvedValue([
      { id: 1, message: 'Test notification', is_read: false }
    ]);

    render(<NotificationCenter />);
    
    await waitFor(() => {
      expect(screen.getByText('1')).toBeInTheDocument();
    });
  });

  test('should handle real-time notifications', async () => {
    render(<NotificationCenter />);
    
    // Get the notification handler
    const notificationHandler = wsService.onMessage.mock.calls.find(
      call => call[0] === 'notification'
    )[1];
    
    // Simulate receiving a notification
    notificationHandler({
      type: 'notification',
      message: 'New friend request',
      timestamp: Date.now()
    });
    
    await waitFor(() => {
      expect(screen.getByText('1')).toBeInTheDocument();
    });
  });

  test('should open dropdown when clicked', () => {
    render(<NotificationCenter />);
    
    const bell = screen.getByText('🔔');
    fireEvent.click(bell);
    
    expect(screen.getByText('Notifications')).toBeInTheDocument();
  });
});
```

## Integration Testing

### Full WebSocket Flow Testing

```javascript
// __tests__/integration/websocket-flow.test.jsx
import React from 'react';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import ChatInterface from '../../components/chat/ChatInterface';

// Use real WebSocket service but with mocked WebSocket
const mockWebSocket = {
  send: jest.fn(),
  close: jest.fn(),
  readyState: WebSocket.OPEN,
  addEventListener: jest.fn(),
  removeEventListener: jest.fn()
};

global.WebSocket = jest.fn(() => mockWebSocket);

describe('WebSocket Integration', () => {
  test('should handle complete message flow', async () => {
    const user = { id: 1, name: 'Test User' };
    render(<ChatInterface user={user} />);
    
    // Wait for connection
    await waitFor(() => {
      expect(global.WebSocket).toHaveBeenCalled();
    });
    
    // Simulate connection open
    const wsInstance = global.WebSocket.mock.results[0].value;
    wsInstance.onopen();
    
    // Wait for connection status update
    await waitFor(() => {
      expect(screen.getByTitle('WebSocket connected')).toBeInTheDocument();
    });
    
    // Send a message
    const input = screen.getByPlaceholderText('Type a message...');
    fireEvent.change(input, { target: { value: 'Test message' } });
    
    const sendButton = screen.getByText('Send');
    fireEvent.click(sendButton);
    
    // Verify message was sent
    expect(mockWebSocket.send).toHaveBeenCalledWith(
      expect.stringContaining('Test message')
    );
  });
});
```

## Mocking Strategies

### WebSocket Service Mock

```javascript
// __mocks__/websocket.js
export const wsService = {
  connect: jest.fn(),
  disconnect: jest.fn(),
  sendMessage: jest.fn(),
  onMessage: jest.fn(),
  isConnected: jest.fn(() => true),
  
  // Test utilities
  simulateMessage: (type, data) => {
    const handlers = wsService._handlers.get(type) || [];
    handlers.forEach(handler => handler(data));
  },
  
  simulateConnectionStatus: (status) => {
    wsService.simulateMessage('connection_status', { status });
  },
  
  _handlers: new Map()
};

// Override onMessage to store handlers for testing
wsService.onMessage.mockImplementation((type, handler) => {
  if (!wsService._handlers.has(type)) {
    wsService._handlers.set(type, []);
  }
  wsService._handlers.get(type).push(handler);
});
```

### API Service Mock

```javascript
// __mocks__/api.js
export const chatAPI = {
  getPrivateMessages: jest.fn(() => Promise.resolve([])),
  getGroupMessages: jest.fn(() => Promise.resolve([])),
  sendGroupInvite: jest.fn(() => Promise.resolve()),
  getNotifications: jest.fn(() => Promise.resolve([])),
  markNotificationsRead: jest.fn(() => Promise.resolve())
};
```

## End-to-End Testing

### Cypress E2E Tests

```javascript
// cypress/integration/websocket.spec.js
describe('WebSocket Real-time Features', () => {
  beforeEach(() => {
    cy.visit('/me');
    cy.login('test@example.com', 'password');
  });

  it('should connect to WebSocket and show connection status', () => {
    cy.get('[data-testid="connection-indicator"]')
      .should('have.class', 'bg-green-500');
  });

  it('should send and receive messages', () => {
    // Open chat with user
    cy.get('[data-testid="chat-user-1"]').click();
    
    // Send message
    cy.get('[data-testid="message-input"]')
      .type('Hello from Cypress test{enter}');
    
    // Verify message appears
    cy.get('[data-testid="message-list"]')
      .should('contain', 'Hello from Cypress test');
  });

  it('should show real-time notifications', () => {
    // Simulate notification from another user
    cy.window().then((win) => {
      win.wsService.simulateMessage('notification', {
        type: 'notification',
        subtype: 'friend_request',
        message: 'You have a new friend request'
      });
    });
    
    // Check notification appears
    cy.get('[data-testid="notification-badge"]')
      .should('contain', '1');
  });
});
```

## Performance Testing

### WebSocket Connection Performance

```javascript
// __tests__/performance/websocket-performance.test.js
describe('WebSocket Performance', () => {
  test('should handle multiple rapid messages', async () => {
    const wsService = new WebSocketService();
    const messageHandler = jest.fn();
    
    wsService.onMessage('test', messageHandler);
    
    // Simulate rapid messages
    const messages = Array.from({ length: 100 }, (_, i) => ({
      type: 'test',
      content: `Message ${i}`
    }));
    
    const start = performance.now();
    
    messages.forEach(msg => wsService.handleMessage(msg));
    
    const end = performance.now();
    
    expect(messageHandler).toHaveBeenCalledTimes(100);
    expect(end - start).toBeLessThan(100); // Should process 100 messages in <100ms
  });
});
```

## Test Configuration

### Jest Setup

```javascript
// jest.config.js
module.exports = {
  testEnvironment: 'jsdom',
  setupFilesAfterEnv: ['<rootDir>/src/setupTests.js'],
  moduleNameMapping: {
    '^@/(.*)$': '<rootDir>/src/$1'
  },
  collectCoverageFrom: [
    'src/**/*.{js,jsx}',
    '!src/**/*.test.{js,jsx}',
    '!src/index.js'
  ]
};
```

### Test Setup

```javascript
// src/setupTests.js
import '@testing-library/jest-dom';

// Mock WebSocket globally
global.WebSocket = jest.fn(() => ({
  close: jest.fn(),
  send: jest.fn(),
  addEventListener: jest.fn(),
  removeEventListener: jest.fn(),
  readyState: WebSocket.OPEN
}));

// Mock environment variables
process.env.NEXT_PUBLIC_API_URL = 'http://localhost:9000';
```

This testing guide provides comprehensive coverage for testing WebSocket functionality in React components, ensuring reliable real-time features.
