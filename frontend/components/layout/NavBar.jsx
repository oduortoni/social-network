import Link from "next/link";
import { handleLogout } from '../../lib/auth';
import NotificationCenter from '../notifications/NotificationCenter';

const NavBar = ({ avatar }) => {

  const fallbackAvatar =
    "data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iNDAiIGhlaWdodD0iNDAiIHZpZXdCb3g9IjAgMCA0MCA0MCIgZmlsbD0ibm9uZSIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj4KPGNpcmNsZSBjeD0iMjAiIGN5PSIyMCIgcj0iMjAiIGZpbGw9IiNGM0Y0RjYiLz4KPGNpcmNsZSBjeD0iMjAiIGN5PSIxNiIgcj0iNiIgZmlsbD0iIzlDQTNBRiIvPgo8cGF0aCBkPSJNMzIgMzJDMzIgMjYuNDc3MiAyNy41MjI4IDIyIDIyIDIySDE4QzEyLjQ3NzIgMjIgOCAyNi40NzcyIDggMzJWMzJIMzJWMzJaIiBmaWxsPSIjOUNBM0FGIi8+Cjwvc3ZnPgo=";

  const handleImageError = (e) => {
    e.target.src = fallbackAvatar;
    setErrorCount((prev) => prev + 1);
  };

  return (
    <nav className="w-full bg-white shadow-sm p-4 flex justify-between items-center">
      <div className="flex items-center">
        <img
          src={avatar ? `http://localhost:9000/avatar?avatar=${avatar}` : '/default.png'}
          alt="User Avatar"
          className="w-10 h-10 rounded-full object-cover border-2 border-gray-200"
          onError={handleImageError}
        />
      </div>
      <div className="flex items-center space-x-4">
        <Link href="/me" className="text-gray-700 hover:text-indigo-600 font-medium">
          Dashboard
        </Link>
        <NotificationCenter />
      </div>
      <button
        onClick={handleLogout}
        className="text-red-600 hover:text-red-800 font-medium transition-colors cursor-pointer"
      >
        Logout
      </button>
    </nav>
  );
};

export default NavBar;
