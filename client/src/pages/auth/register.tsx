import { useEffect, useRef, useState } from "react";
import api from "../../lib/axios";
import axios from "axios";

export default function Register() {
  const [state, setState] = useState({
    messages: [
      "Register today to start your journaling journey",
      "Enter your username",
    ],
    input: "",
    stage: "username" as
      | "username"
      | "password"
      | "confirmPassword"
      | "finished",
    username: "",
    password: "",
    confirmPassword: "",
  });

  const inputRef = useRef<HTMLInputElement>(null);

  useEffect(() => {
    inputRef.current?.focus();
  }, [state.messages, state.stage]);

  const handleRegister = async (username: string, password: string) => {
    console.log("USERNAME IS", username);
    console.log("PASSWORD IS", password);

    try {
      const response = await api.post("/api/register", { username, password });

      console.log(response.status);

      setState((prev) => ({
        ...prev,
        messages: [...prev.messages, "Registration successful!"],
      }));
    } catch (error) {
      if (axios.isAxiosError(error)) {
        if (error.response?.status === 400) {
          setState((prev) => ({
            ...prev,
            messages: [
              ...prev.messages,
              "Username already taken",
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
        stage: "confirmPassword",
        messages: [...prev.messages, maskedPassword, "Confirm your password"],
        input: "",
      }));
      return;
    }

    if (state.stage === "confirmPassword") {
      if (state.input !== state.password) {
        setState((prev) => ({
          ...prev,
          messages: [
            ...prev.messages,
            "Passwords do not match, try again.",
            "Enter your password",
          ],
          input: "",
          stage: "password",
        }));
        return;
      }

      setState((prev) => ({
        ...prev,
        confirmPassword: state.input,
        stage: "finished",
        messages: [
          ...prev.messages,
          "*".repeat(state.input.length),
          "Registering...",
        ],
        input: "",
      }));

      handleRegister(state.username, state.password);
      return;
    }

    setState((prev) => ({
      ...prev,
      messages: [...prev.messages, prev.input],
      input: "",
    }));
  };

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
            autoComplete="off"
            spellCheck={false}
            ref={inputRef}
            type={state.stage === "username" ? "text" : "password"}
            disabled={state.stage === "finished"}
            value={state.input}
            autoCorrect="off"
            className="absolute opacity-0 left-0"
            onKeyDown={handleKeyDown}
            inputMode="none"
            onChange={(e) =>
              setState((prev) => ({ ...prev, input: e.target.value }))
            }
          />

          {/* Visible Text + Cursor */}
          <span className="whitespace-pre">
            {state.input}
            <span className="text-white animate-blink">█</span>
          </span>
        </div>
      </div>
    </div>
  );
}
