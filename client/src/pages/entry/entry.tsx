import { useEffect, useRef } from "react";
import { useNavigate, useParams } from "react-router";
import useFetch from "../../hooks/use-fetch";
import { JournalEntry } from "../../types/types";
import { formatISO } from "date-fns";
import { useAuth } from "../../store/auth";
interface EntryResponse {
  entry: JournalEntry;
}

export default function Entry({
  setAllowSave,
  content,
  setContent,
  selectedEntry,
  setSelectedEntry,
}: {
  setAllowSave: React.Dispatch<React.SetStateAction<boolean>>;
  content: string;
  setContent: React.Dispatch<React.SetStateAction<string>>;
  selectedEntry: JournalEntry | null;
  setSelectedEntry: React.Dispatch<React.SetStateAction<JournalEntry | null>>;
}) {
  const { id } = useParams();
  const { data, error, loading } = useFetch<EntryResponse>(
    id ? "/api/entry/" + id : null
  );
  const inputRef = useRef<HTMLTextAreaElement>(null);

  const { user, isLoading } = useAuth();
  const navigate = useNavigate();

  useEffect(() => {
    if (data) {
      setSelectedEntry(data.entry);
      setContent(data.entry.content);
    } else if (!id) {
      const todaysCompleteDate = formatISO(new Date(), {
        representation: "complete",
      });
      setSelectedEntry({
        ID: 0,
        date: todaysCompleteDate,
        content: "",
      } as JournalEntry);
    }
  }, [data, setContent, setSelectedEntry, id]);

  useEffect(() => {
    if (!selectedEntry) {
      setContent("");
    }
    inputRef.current?.focus();
  }, [setContent, selectedEntry]);

  useEffect(() => {
    if (
      content !== selectedEntry?.content &&
      content.split(/\s+/).filter(Boolean).length >= 150
    ) {
      setAllowSave(true);
    } else {
      setAllowSave(false);
    }
  }, [content, setAllowSave, selectedEntry]);

  if (!isLoading && !user) {
    setSelectedEntry(null);
    navigate("/login");
    return null;
  }

  if (loading)
    return (
      <div className="flex-1 p-5 animate-pulse space-y-2">
        <span className="w-[80%] bg-zinc-900 rounded-sm h-4 block"></span>
        <span className="w-[50%] bg-zinc-900 rounded-sm h-4 block"></span>
        <span className="w-[70%] bg-zinc-900 rounded-sm h-4 block"></span>
        <span className="w-[40%] bg-zinc-900 rounded-sm h-4 block"></span>
        <span className="w-[50%] bg-zinc-900 rounded-sm h-4 block"></span>
        <span className="w-[60%] bg-zinc-900 rounded-sm h-4 block"></span>
        <span className="w-[30%] bg-zinc-900 rounded-sm h-4 block"></span>
        <span className="w-[80%] bg-zinc-900 rounded-sm h-4 block"></span>
        <span className="w-[70%] bg-zinc-900 rounded-sm h-4 block"></span>
        <span className="w-[80%] bg-zinc-900 rounded-sm h-4 block"></span>
        <span className="w-[50%] bg-zinc-900 rounded-sm h-4 block"></span>
        <span className="w-[20%] bg-zinc-900 rounded-sm h-4 block"></span>
      </div>
    );

  if (error)
    return (
      <div className="flex-1 flex flex-col items-center justify-center gap-3">
        <p className="font-medium text-zinc-400">Error: Could not load entry</p>
        <p className="text-sm">{error}</p>
      </div>
    );

  return (
    <div
      className="flex-1 flex flex-col p-5"
      onClick={() => inputRef.current?.focus()}
    >
      <textarea
        autoCorrect="off"
        autoCapitalize="off"
        spellCheck="false"
        ref={inputRef}
        value={content}
        className="w-full outline-none border-none resize-none flex-1 noscrollbar absolute left-0 opacity-0"
        style={{
          wordWrap: "break-word",
          whiteSpace: "pre-wrap",
        }}
        onChange={(e) => setContent(e.target.value)}
      />
      <pre className="whitespace-pre-wrap">
        {content}
        <span className="text-white animate-blink">█</span>
      </pre>
    </div>
  );
}
