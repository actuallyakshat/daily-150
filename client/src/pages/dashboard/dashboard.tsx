import React, { useEffect } from "react";
import { Link, useNavigate } from "react-router";
import { useAuth } from "../../store/auth";
import { JournalEntry } from "../../types/types";
export default function Dashboard({
  setSelectedEntry,
}: {
  setSelectedEntry: React.Dispatch<React.SetStateAction<JournalEntry | null>>;
}) {
  const { user, isLoading } = useAuth();
  const navigate = useNavigate();

  useEffect(() => {
    setSelectedEntry(null);
  }, [setSelectedEntry]);

  if (!isLoading && !user) {
    navigate("/login");
    return null;
  }

  if (!user) return null;

  return (
    <div className="flex-1 p-5">
      <h2 className="text-sm text-zinc-400">
        You have made{" "}
        {user?.journal_entries && user?.journal_entries?.length == 1
          ? "1 entry "
          : user?.journal_entries?.length + " entries "}
        yet. each cell here represents a day.
      </h2>
      <div className="flex items-center gap-1.5 mt-4">
        {user?.journal_entries
          ?.sort((a, b) => Number(new Date(b.date)) - Number(new Date(a.date))) // Correct sorting logic
          .map((entry: JournalEntry) => (
            <Link
              to={`/entry/${entry.ID}`}
              key={entry.ID}
              onClick={() => setSelectedEntry(entry)}
              className="size-4 bg-zinc-400 hover:bg-white transition-colors relative group cursor-pointer p-2"
            >
              <span className="absolute left-0 -top-[200%] text-sm hidden group-hover:block text-white bg-zinc-900 py-1 px-3 rounded transition-all duration-300">
                {new Date(entry.date).toLocaleDateString()}
              </span>
            </Link>
          ))}
      </div>
    </div>
  );
}
