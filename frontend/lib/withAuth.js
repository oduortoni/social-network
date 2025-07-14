
import { useEffect, useState } from 'react';

const withAuth = (WrappedComponent) => {
  return function AuthenticatedComponent(props) {
    const [user, setUser] = useState(null);
    const [loading, setLoading] = useState(true);
    const [isAuthenticated, setIsAuthenticated] = useState(false);

    useEffect(() => {
      const checkAuth = async () => {
        try {
          const response = await fetch('http://localhost:9000/me', {
            method: 'GET',
            credentials: 'include',
          });

          if (response.ok) {
            const userData = await response.json();
            setUser(userData);
            setIsAuthenticated(true);
          } else {
            window.location.href = '/login';
          }
        } catch (error) {
          console.error('Auth check failed:', error);
          window.location.href = '/login';
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

    return <WrappedComponent {...props} user={user} />;
  };
};

export default withAuth;