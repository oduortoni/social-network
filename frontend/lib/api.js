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


const fallbackAvatar =
  "data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iNDAiIGhlaWdodD0iNDAiIHZpZXdCb3g9IjAgMCA0MCA0MCIgZmlsbD0ibm9uZSIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj4KPGNpcmNsZSBjeD0iMjAiIGN5PSIyMCIgcj0iMjAiIGZpbGw9IiNGM0Y0RjYiLz4KPGNpcmNsZSBjeD0iMjAiIGN5PSIxNiIgcj0iNiIgZmlsbD0iIzlDQTNBRiIvPgo8cGF0aCBkPSJNMzIgMzJDMzIgMjYuNDc3MiAyNy41MjI4IDIyIDIyIDIySDE4QzEyLjQ3NzIgMjIgOCAyNi40NzcyIDggMzJWMzJIMzJWMzJaIiBmaWxsPSIjOUNBM0FGIi8+Cjwvc3ZnPgo=";

function fetchProfileImage(avatar) {
  if (!avatar) return fallbackAvatar;
  return `${API_BASE}/avatar?avatar=${encodeURIComponent(avatar)}`;
}

function fetchVerifiedBadge() {
  // TODO: Fetch verified badge status from backend
  return null;
}

function fetchFollowers() {
  // TODO: Fetch followers count from backend
  return null;
}

function fetchFollowing() {
  // TODO: Fetch following count from backend
  return null;
}

function fetchProfileStatus() {
  // TODO: Fetch profile status from backend
  return null;
}

function fetchCommunities() {
  // TODO: Fetch communities list from backend
  return null;
}

export const profileAPI = {
  fetchProfileImage,
  fetchVerifiedBadge,
  fetchFollowers,
  fetchFollowing,
  fetchProfileStatus,
  fetchCommunities,
};
