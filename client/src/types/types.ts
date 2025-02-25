interface User {
  ID: number;
  username: string;
  journal_entries: JournalEntry[];
  summaries: Summary[];
}

interface JournalEntry {
  ID: number;
  user_id: number;
  date: string;
  content: string;
}

interface Summary {
  ID: number;
  user_id: number;
  week_number: number;
  summary: string;
}

export type { User, JournalEntry, Summary };
