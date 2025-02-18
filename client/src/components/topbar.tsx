import { Link } from "react-router";
import { useAuth } from "../store/auth";
import { useEffect } from "react";

export default function TopBar() {
  const { user, logout } = useAuth();
  useEffect(() => {
    console.log(user);
  }, [user]);

  return (
    <div className="h-8 sm:h-12 px-2 sm:px-5 bg-zinc-900 flex items-center justify-between text-xs sm:text-sm">
      <Link
        to={`/entry/${new Date().toLocaleDateString()}`}
        className="hidden sm:block font-medium"
      >
        {new Date().toLocaleDateString()}
      </Link>
      {user && (
        <div className="flex items-center justify-center gap-5">
          <Link
            className="hover:text-lime-500 transition-colors"
            to="/entries/all"
          >
            {user?.username}
          </Link>
          <button
            onClick={logout}
            className="cursor-pointer hover:text-lime-500 transition-colors"
          >
            logout
          </button>
        </div>
      )}
      {!user && (
        <div className="flex items-center justify-center gap-5">
          <Link
            to="/login"
            className="text-zinc-300 hover:text-lime-500 transition-colors text-sm font-medium"
          >
            login
          </Link>
          <Link
            to="/register"
            className="text-zinc-300 text-sm font-medium hover:text-lime-500 transition-colors"
          >
            register
          </Link>
        </div>
      )}
    </div>
  );
}
