export default function InputArea() {
  return (
    <div className="p-2 sm:p-3 md:p-4 lg:p-5 flex-grow flex flex-col overflow-hidden">
      <div className="flex flex-grow overflow-hidden">
        <span className="mr-2 flex-shrink-0">$</span>
        <textarea
          className="w-full bg-transparent border-none outline-none resize-none noscrollbar flex-grow"
          style={{
            wordWrap: "break-word",
            whiteSpace: "pre-wrap",
          }}
        />
      </div>
    </div>
  );
}
