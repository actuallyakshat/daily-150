import axios from "axios";
import { useEffect, useRef, useState } from "react";
import { useAuth } from "../../store/auth";
import { useNavigate } from "react-router";

export default function Login() {
  const [state, setState] = useState({
    messages: ["Log in to Daily 150", "Enter your username"],
    input: "",
    stage: "username" as "username" | "password" | "finished",
    username: "",
    password: "",
  });

  const { login, user } = useAuth();

  const inputRef = useRef<HTMLInputElement>(null);

  const navigate = useNavigate();

  useEffect(() => {
    inputRef.current?.focus();
  }, [state.messages, state.stage]);

  const handleLogin = async (username: string, password: string) => {
    try {
      // const response = await api.post("/api/login", { username, password });
      await login(username, password);

      setState((prev) => ({
        ...prev,
        messages: [...prev.messages, "Login successful!"],
      }));
    } catch (error) {
      if (axios.isAxiosError(error)) {
        if (error.response?.status === 401) {
          setState((prev) => ({
            ...prev,
            messages: [
              ...prev.messages,
              "Invalid username or password",
              "Enter your username",
            ],
            input: "",
            stage: "username",
          }));
          return;
        }

        setState((prev) => ({
          ...prev,
          messages: [
            ...prev.messages,
            `Error: ${error.response?.status || "Unknown"}`,
            "Enter your username",
          ],
        }));
      } else {
        console.error("Unexpected error:", error);
        setState((prev) => ({
          ...prev,
          messages: [
            ...prev.messages,
            "Unexpected error occurred.",
            "Enter your username",
          ],
        }));
      }
    }
  };

  const handleKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key !== "Enter" || state.input.length === 0) return;

    if (state.stage === "username") {
      setState((prev) => ({
        ...prev,
        username: prev.input,
        stage: "password",
        messages: [...prev.messages, prev.input, "Enter your password"],
        input: "",
      }));

      return;
    }

    if (state.stage === "password") {
      const maskedPassword = "*".repeat(state.input.length);
      const enteredPassword = state.input;

      setState((prev) => ({
        ...prev,
        password: enteredPassword,
        stage: "finished",
        messages: [...prev.messages, maskedPassword, "Logging in..."],
        input: "",
      }));

      handleLogin(state.username, enteredPassword);
      return;
    }

    setState((prev) => ({
      ...prev,
      messages: [...prev.messages, prev.input],
      input: "",
    }));
  };

  useEffect(() => {
    if (user) {
      navigate("/dashboard");
    }
  }, [user, navigate]);

  return (
    <div
      className="p-2 sm:p-3 md:p-4 lg:p-5 flex-grow flex flex-col overflow-hidden"
      onClick={() => inputRef.current?.focus()}
    >
      <div className="flex flex-col overflow-hidden">
        {state.messages.map((message, index) => (
          <div key={index}>
            <p>
              <span className="mr-2 flex-shrink-0">$</span>
              {message}
            </p>
          </div>
        ))}
        <div>
          <span className="mr-2 flex-shrink-0">$</span>
          <input
            autoComplete="new-password"
            spellCheck={false}
            ref={inputRef}
            type="text"
            disabled={state.stage === "finished"}
            value={state.input}
            autoCorrect="off"
            className="absolute opacity-0 left-0"
            onKeyDown={handleKeyDown}
            // inputMode="none"
            onChange={(e) =>
              setState((prev) => ({ ...prev, input: e.target.value }))
            }
          />

          {/* Visible Text + Cursor */}
          <span className="whitespace-pre">
            {state.stage === "password"
              ? "*".repeat(state.input.length)
              : state.input}
            <span className="text-white animate-blink">█</span>
          </span>
        </div>
      </div>
    </div>
  );
}
