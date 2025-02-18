interface User {
  ID: number;
  username: string;
  JournalEntries: JournalEntry[];
}

interface JournalEntry {
  ID: number;
  userID: number;
  date: Date;
  content: string;
}

export type { User, JournalEntry };
