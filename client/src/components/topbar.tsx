import { Link, useLocation, useNavigate } from "react-router";
import api from "../lib/axios";
import { useAuth } from "../store/auth";
import { JournalEntry } from "../types/types";
import { useCallback, useEffect, useState } from "react";

interface TopBarProps {
  selectedEntry: JournalEntry | null;
  allowSave: boolean;
  content: string;
  allowNewEntry: boolean;
  setAllowSave: React.Dispatch<React.SetStateAction<boolean>>; // Add this to your interface
  setSelectedEntry: React.Dispatch<React.SetStateAction<JournalEntry | null>>;
}

export default function TopBar({
  selectedEntry,
  allowSave,
  content,
  allowNewEntry,
  setAllowSave,
  setSelectedEntry,
}: TopBarProps) {
  const { user, logout, refreshUser } = useAuth();
  const [isLoading, setIsLoading] = useState({
    save: false,
    delete: false,
  });

  const navigate = useNavigate();
  const { pathname } = useLocation();

  const onSave = useCallback(async () => {
    try {
      if (!selectedEntry) return;
      setIsLoading((prevState) => ({ ...prevState, save: true }));
      if (selectedEntry.ID === 0) {
        const response = await api.post("/api/entry", {
          content,
          date: selectedEntry.date,
        });
        console.log(response);
        refreshUser();
        setSelectedEntry(response.data.entry);
        return;
      }

      const response = await api.patch("/api/entry/" + selectedEntry.ID, {
        content,
      });
      setSelectedEntry(response.data.entry);
    } catch (error) {
      console.error(error);
    } finally {
      setIsLoading((prevState) => ({ ...prevState, save: false }));
    }
  }, [content, selectedEntry, setSelectedEntry, refreshUser]);

  const deleteEntry = useCallback(async () => {
    try {
      setIsLoading((prevState) => ({ ...prevState, delete: true }));
      if (!selectedEntry) return;
      await api.delete("/api/entry/" + selectedEntry.ID);
      setSelectedEntry(null);
      refreshUser();
      navigate("/dashboard");
    } catch (error) {
      console.error(error);
    } finally {
      setIsLoading((prevState) => ({ ...prevState, delete: false }));
    }
  }, [selectedEntry, setSelectedEntry, navigate, refreshUser]);

  useEffect(() => {
    if (
      content !== selectedEntry?.content &&
      content.split(/\s+/).filter(Boolean).length >= 150
    ) {
      setAllowSave(true);
    } else {
      setAllowSave(false);
    }
  }, [content, selectedEntry, setAllowSave]);

  useEffect(() => {
    const handleKeyDown = async (event: KeyboardEvent) => {
      if (event.ctrlKey && event.key === "s") {
        event.preventDefault();
        if (allowSave) {
          await onSave();
        }
      }
    };

    document.addEventListener("keydown", handleKeyDown);
    return () => document.removeEventListener("keydown", handleKeyDown);
  }, [allowSave, onSave]);

  useEffect(() => {
    if (pathname != "/entry") {
      setSelectedEntry(null);
    }
  }, [pathname, setSelectedEntry]);

  return (
    <div className="min-h-8 sm:min-h-12 px-2 sm:px-5 bg-zinc-900 flex items-center justify-between text-xs sm:text-sm">
      <div className="flex items-center gap-4">
        <p className="hidden sm:block font-medium">
          {selectedEntry
            ? new Date(selectedEntry.date.split("T")[0]).toLocaleDateString()
            : new Date().toLocaleDateString()}
        </p>

        {selectedEntry && (
          <div className="flex items-center justify-center gap-4">
            <p
              style={{
                color:
                  content.split(/\s+/).filter(Boolean).length < 150
                    ? "red"
                    : "#7ccf00",
              }}
            >
              {content.split(/\s+/).filter(Boolean).length} / 150
            </p>
            {selectedEntry.ID !== 0 && (
              <button
                onClick={deleteEntry}
                className="cursor-pointer hover:text-red-500 transition-colors"
              >
                {isLoading.delete ? "deleting..." : "delete"}
              </button>
            )}
          </div>
        )}
      </div>

      {user && (
        <div className="flex items-center justify-center gap-5">
          {allowSave && (
            <button
              onClick={onSave}
              className="cursor-pointer hover:text-lime-500 transition-colors"
            >
              {isLoading.save ? "saving..." : "save"}
            </button>
          )}

          {!selectedEntry && allowNewEntry && (
            <Link
              to="/entry"
              className="cursor-pointer hover:text-lime-500 transition-colors"
            >
              make today's entry
            </Link>
          )}

          <Link
            to="/summaries"
            className="cursor-pointer hover:text-lime-500 transition-colors"
          >
            summaries
          </Link>

          <Link
            className="hover:text-lime-500 transition-colors"
            to="/dashboard"
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
