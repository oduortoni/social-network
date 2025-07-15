import Link from "next/link";
import { handleLogout } from '../../lib/auth';

const NavBar = () => {
  return (
    <nav className="w-full bg-white shadow-sm p-4 flex justify-between">
      <div className="flex items-center space-x-4">
        <Link href="/dashboard" className="text-gray-700 hover:text-indigo-600 font-medium">
          Dashboard
        </Link>
        <Link href="/me" className="text-gray-700 hover:text-indigo-600 font-medium">
          Profile
        </Link>
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
