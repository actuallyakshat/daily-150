import React from "react";

export default function TerminalWindow({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <div className="min-h-screen bg-black text-white flex items-center justify-center p-4">
      <div className="w-full max-w-screen-sm sm:max-w-screen-md lg:max-w-screen-lg xl:max-w-screen-xl border border-zinc-800 rounded-lg shadow-lg text-sm sm:text-base md:text-lg h-[90vh] sm:h-[80vh] md:h-[70vh] lg:h-[60vh] overflow-hidden flex flex-col">
        {/* Top Bar */}
        {/* <TopBar wordCount={wordCount} onSave={onSave} /> */}
        {/* Input area */}
        {children}
      </div>
    </div>
  );
}
