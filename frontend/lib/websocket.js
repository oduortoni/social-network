class WebSocketService {
  constructor() {
    this.ws = null;
    this.reconnectAttempts = 0;
    this.maxReconnectAttempts = 5;
    this.messageHandlers = new Map(); // Map<string, Array<function>>
    this.isConnecting = false;
    this.shouldReconnect = true;
    this.reconnectTimeout = null;
  }

  connect() {
    // Prevent multiple connection attempts
    if (this.isConnecting || (this.ws && this.ws.readyState === WebSocket.CONNECTING)) {
      return;
    }

    // Don't reconnect if we're already connected
    if (this.ws && this.ws.readyState === WebSocket.OPEN) {
      return;
    }

    this.isConnecting = true;
    this.shouldReconnect = true;

    // Extract session_id from cookies
    const getSessionId = () => {
      const cookies = document.cookie.split(';');
      for (let cookie of cookies) {
        const [name, value] = cookie.trim().split('=');
        if (name === 'session_id') {
          return value;
        }
      }
      return null;
    };

    const sessionId = getSessionId();
    let wsUrl = process.env.NEXT_PUBLIC_API_URL?.replace('http', 'ws') + '/ws';

    // Add session ID as query parameter if available
    if (sessionId) {
      wsUrl += `?session_id=${sessionId}`;
    }

    try {
      // Note: WebSocket constructor doesn't support credentials option
      // Cookies should be automatically included for same-origin requests
      this.ws = new WebSocket(wsUrl);

      this.ws.onopen = () => {
        this.reconnectAttempts = 0;
        this.isConnecting = false;
        this.handleMessage({ type: 'connection_status', status: 'connected' });
      };

      this.ws.onmessage = (event) => {
        try {
          const message = JSON.parse(event.data);
          this.handleMessage(message);
        } catch (error) {
          console.error('Failed to parse WebSocket message:', error);
        }
      };

      this.ws.onclose = (event) => {
        this.isConnecting = false;
        this.handleMessage({ type: 'connection_status', status: 'disconnected' });

        // Log specific close codes for debugging
        if (event.code === 1006) {
          console.error('WebSocket closed abnormally - possible authentication failure');
        } else if (event.code === 1002) {
          console.error('WebSocket closed due to protocol error');
        } else if (event.code === 1011) {
          console.error('WebSocket closed due to server error');
        }

        // Only reconnect if it wasn't a manual disconnect
        if (this.shouldReconnect && event.code !== 1000) {
          this.reconnect();
        }
      };

      this.ws.onerror = (error) => {
        console.error('WebSocket error:', error);
        this.isConnecting = false;
        this.handleMessage({ type: 'connection_status', status: 'disconnected' });
      };
    } catch (error) {
      console.error('Failed to create WebSocket:', error);
      this.isConnecting = false;
      this.reconnect();
    }
  }

  sendMessage(type, content, to = null, groupId = null) {
    if (this.ws?.readyState === WebSocket.OPEN) {
      const message = {
        type,
        content,
        ...(to && { to }),
        ...(groupId && { group_id: groupId }),
        timestamp: Date.now()
      };
      try {
        this.ws.send(JSON.stringify(message));
      } catch (error) {
        console.error('Failed to send WebSocket message:', error);
      }
    } else {
      console.warn('WebSocket not connected, message not sent:', { type, content });
    }
  }

  onMessage(type, handler) {
    if (!this.messageHandlers.has(type)) {
      this.messageHandlers.set(type, []);
    }
    this.messageHandlers.get(type).push(handler);
  }

  removeMessageHandler(type, handler) {
    const handlers = this.messageHandlers.get(type);
    if (handlers) {
      const index = handlers.indexOf(handler);
      if (index > -1) {
        handlers.splice(index, 1);
      }
    }
  }

  handleMessage(message) {
    const handlers = this.messageHandlers.get(message.type);
    if (handlers) {
      handlers.forEach(handler => {
        try {
          handler(message);
        } catch (error) {
          console.error('Error in message handler:', error);
        }
      });
    }
  }

  disconnect() {
    this.shouldReconnect = false;

    // Clear any pending reconnection
    if (this.reconnectTimeout) {
      clearTimeout(this.reconnectTimeout);
      this.reconnectTimeout = null;
    }

    if (this.ws) {
      this.ws.close(1000, 'Manual disconnect');
      this.ws = null;
    }
  }

  reconnect() {
    if (!this.shouldReconnect || this.reconnectAttempts >= this.maxReconnectAttempts) {
      return;
    }

    this.reconnectAttempts++;
    const delay = Math.min(1000 * Math.pow(2, this.reconnectAttempts - 1), 10000); // Exponential backoff, max 10s

    this.handleMessage({ type: 'connection_status', status: 'connecting' });

    this.reconnectTimeout = setTimeout(() => {
      this.connect();
    }, delay);
  }

  isConnected() {
    return this.ws?.readyState === WebSocket.OPEN;
  }

  // Force disconnect (for logout or errors)
  forceDisconnect() {
    this.shouldReconnect = false;
    this.disconnect();
  }
}

export const wsService = new WebSocketService();

// Global function to disconnect WebSocket on logout
export const disconnectWebSocketOnLogout = () => {
  wsService.forceDisconnect();
};