const API_BASE = process.env.NEXT_PUBLIC_API_URL;

const apiCall = async (endpoint, options = {}) => {
  const response = await fetch(`${API_BASE}${endpoint}`, {
    credentials: 'include',
    headers: {
      'Content-Type': 'application/json',
      ...options.headers,
    },
    ...options,
  });
  
  if (!response.ok) {
    throw new Error(`API Error: ${response.status}`);
  }
  
  return response.json();
};

export const chatAPI = {
  // Get private message history
  getPrivateMessages: (userId, limit = 50, offset = 0) =>
    apiCall(`/api/messages/private?user=${userId}&limit=${limit}&offset=${offset}`),
  
  // Get group message history
  getGroupMessages: (groupId, limit = 50, offset = 0) =>
    apiCall(`/api/messages/group?group=${groupId}&limit=${limit}&offset=${offset}`),
  
  // Send group invitation
  sendGroupInvite: (groupId, userId, groupName) =>
    apiCall('/api/groups/invite', {
      method: 'POST',
      body: JSON.stringify({ group_id: groupId, user_id: userId, group_name: groupName }),
    }),
  
  // Get notifications
  getNotifications: (limit = 20, offset = 0) =>
    apiCall(`/api/notifications?limit=${limit}&offset=${offset}`),
  
  // Mark notifications as read
  markNotificationsRead: () =>
    apiCall('/api/notifications/read', { method: 'POST' }),
};
