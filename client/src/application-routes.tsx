import { Route, Routes } from "react-router";
import TerminalWindow from "./components/terminal-window";
import TopBar from "./components/topbar";
import Landing from "./pages/landing/landing";
import { AuthProvider } from "./providers/auth-provider";
import Login from "./pages/auth/login";
import Register from "./pages/auth/register";
import NotFound from "./components/not-found";
import Dashboard from "./pages/dashboard/dashboard";

export default function ApplicationRoutes() {
  return (
    <AuthProvider>
      {/* <Navbar /> */}
      <div className="text-white">
        <TerminalWindow>
          <TopBar />
          <Routes>
            <Route path="/" element={<Landing />} />
            <Route path="/login" element={<Login />} />
            <Route path="/register" element={<Register />} />
            <Route path="/dashboard" element={<Dashboard />} />
            <Route path="/*" element={<NotFound />} />
          </Routes>
        </TerminalWindow>
      </div>
    </AuthProvider>
  );
}
