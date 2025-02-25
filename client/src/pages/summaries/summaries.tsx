import { Link, useNavigate } from "react-router";
import { useAuth } from "../../store/auth";

export default function Summaries() {
  const { user, isLoading } = useAuth();
  const navigate = useNavigate();

  if (!user) return null;

  if (!isLoading && !user) {
    navigate("/login");
  }

  return (
    <div className="flex-1 p-5">
      <h1 className="font-medium text-xl">Weekly Summaries</h1>
      <p className="text-zinc-400 text-sm">
        A list of weekly summaries for {user?.username}. Updated every monday.
      </p>
      <div className="mt-4">
        {user?.summaries.map((summary) => (
          <Link
            key={summary.ID}
            to={`/summary/${summary.ID}`}
            className="p-2 text-sm aspect-square bg-zinc-300 hover:bg-white transition-colors text-black size-8 flex items-center justify-center"
          >
            {summary.week_number}
          </Link>
        ))}
        {user?.summaries.length === 0 && (
          <p className="text-zinc-400 text-sm mt-4">
            No weekly summary generated yet. Keep journaling and we will
            generate one by this monday.
          </p>
        )}
      </div>
    </div>
  );
}
