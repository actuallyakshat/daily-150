import { Link } from "react-router";

export default function Landing() {
  return (
    <div className="flex-grow flex flex-col gap-2 items-center justify-center h-full px-4">
      <h1 className="text-4xl md:text-5xl font-black text-white">Daily 150</h1>
      <p className="text-zinc-300 text-lg text-center md:text-xl">
        Start your day by journaling one fifty words.
      </p>
      <div className="flex gap-3 mt-2 items-center justify-center animate-pulse text-lime-500">
        <Link
          to="/login"
          className="text-lg md:text-xl font-medium cursor-pointer"
        >
          Get Started
        </Link>
      </div>
    </div>
  );
}
