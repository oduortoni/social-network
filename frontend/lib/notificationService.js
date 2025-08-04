class NotificationService {
  constructor() {
    this.handlers = new Map(); // Map<string, Array<function>>
  }

  onNotification(type, handler) {
    if (!this.handlers.has(type)) {
      this.handlers.set(type, []);
    }
    this.handlers.get(type).push(handler);
  }

  removeHandler(type, handler) {
    if (this.handlers.has(type)) {
      const handlers = this.handlers.get(type);
      const index = handlers.indexOf(handler);
      if (index > -1) {
        handlers.splice(index, 1);
      }
    }
  }

  handleNotification(notification) {
    const { type, subtype } = notification;
    
    // Handle user connection/disconnection notifications
    if (type === 'notification') {
      console.log("100000")
      if (subtype === 'user_connected' || subtype === 'user_disconnected') {
        const handlers = this.handlers.get(subtype) || [];
        handlers.forEach(handler => handler(notification));
      }
      if (subtype == 'follow_request') {
        console.log("5000")
        const handlers = this.handlers.get(subtype) || [];
        handlers.forEach(handler => handler(notification));
      }
    }
  }

  async loadOnlineUsers() {
    try {
      const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/api/users/online`, {
        credentials: 'include'
      });
      
      if (response.ok) {
        const data = await response.json();
        return data.online_users || [];
      } else {
        console.error('Failed to load online users:', response.status);
        return [];
      }
    } catch (error) {
      console.error('Failed to load online users:', error);
      return [];
    }
  }
}

export const notificationService = new NotificationService();
