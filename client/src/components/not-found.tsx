import { Link } from "react-router";

export default function NotFound() {
  return (
    <div className="flex-1 flex flex-col gap-2 items-center justify-center">
      <h2 className="text-3xl">woah, seems like you are lost.</h2>
      <Link to="/dashboard" className="text-lime-500 animate-pulse">
        go back to dashboard
      </Link>
    </div>
  );
}
