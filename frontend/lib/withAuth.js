
import { useEffect, useState } from 'react';
import { profileAPI } from './api';

const withAuth = (WrappedComponent) => {
  return function AuthenticatedComponent(props) {
    const [user, setUser] = useState(null);
    const [profile, setProfile] = useState(null);
    const [loading, setLoading] = useState(true);
    const [isAuthenticated, setIsAuthenticated] = useState(false);

  useEffect(() => {
  const checkAuth = async () => {
    try {
      const res = await fetch('http://localhost:9000/me', {
        method: 'GET',
        credentials: 'include',
      });

      if (res.ok) {
        const userData = await res.json();
        setUser(userData);
        setIsAuthenticated(true);

        // Now it's safe to fetch the profile using userData.id
        const profileRes = await profileAPI.getProfile(userData.id);
        setProfile(profileRes); // assuming you have setProfile in scope
      } else {
        window.location.href = '/';
      }
    } catch (error) {
      console.error('Auth check failed:', error);
      window.location.href = '/';
    } finally {
      setLoading(false);
    }
  };

  checkAuth();
}, []);

    if (loading) {
      return <div>Loading...</div>;
    }

    if (!isAuthenticated) {
      return <div>Redirecting...</div>;
    }

    return <WrappedComponent {...props} user={user} profile={profile} />;
  };
};

export default withAuth;