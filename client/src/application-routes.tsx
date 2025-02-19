import { useEffect, useState } from "react";
import { Route, Routes } from "react-router";
import NotFound from "./components/not-found";
import TerminalWindow from "./components/terminal-window";
import TopBar from "./components/topbar";
import Login from "./pages/auth/login";
import Register from "./pages/auth/register";
import Dashboard from "./pages/dashboard/dashboard";
import Entry from "./pages/entry/entry";
import Landing from "./pages/landing/landing";
import { AuthProvider } from "./providers/auth-provider";
import { JournalEntry } from "./types/types";

export default function ApplicationRoutes() {
  const [selectedEntry, setSelectedEntry] = useState<JournalEntry | null>(null);
  const [allowSave, setAllowSave] = useState(false);
  const [content, setContent] = useState("");
  const [allowNewEntry, setAllowNewEntry] = useState(false);

  useEffect(() => {
    if (!selectedEntry) {
      setAllowSave(false);
    }
  }, [selectedEntry]);

  return (
    <AuthProvider>
      {/* <Navbar /> */}
      <div className="text-white">
        <TerminalWindow>
          <TopBar
            selectedEntry={selectedEntry}
            allowSave={allowSave}
            content={content}
            allowNewEntry={allowNewEntry}
            setAllowSave={setAllowSave}
            setSelectedEntry={setSelectedEntry}
          />
          <Routes>
            <Route path="/" element={<Landing />} />
            <Route path="/login" element={<Login />} />
            <Route path="/register" element={<Register />} />
            <Route
              path="/dashboard"
              element={
                <Dashboard
                  setSelectedEntry={setSelectedEntry}
                  setAllowNewEntry={setAllowNewEntry}
                />
              }
            />

            <Route
              path="/entry"
              element={
                <Entry
                  setAllowSave={setAllowSave}
                  content={content}
                  setContent={setContent}
                  selectedEntry={selectedEntry}
                  setSelectedEntry={setSelectedEntry}
                />
              }
            />

            <Route
              path="/entry/:id"
              element={
                <Entry
                  setAllowSave={setAllowSave}
                  content={content}
                  setContent={setContent}
                  selectedEntry={selectedEntry}
                  setSelectedEntry={setSelectedEntry}
                />
              }
            />

            <Route path="/*" element={<NotFound />} />
          </Routes>
        </TerminalWindow>
      </div>
    </AuthProvider>
  );
}
