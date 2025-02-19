interface User {
  ID: number;
  username: string;
  journal_entries: JournalEntry[];
}

interface JournalEntry {
  ID: number;
  user_id: number;
  date: string;
  content: string;
}

export type { User, JournalEntry };
