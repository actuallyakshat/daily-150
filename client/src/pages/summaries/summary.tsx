import { useParams } from "react-router";
import useFetch from "../../hooks/use-fetch";
import { Summary } from "../../types/types";
import ReactMarkdown from "react-markdown";

interface SummaryResponse {
  summary: Summary;
}

export default function SummaryPage() {
  const { id } = useParams();
  const { data, error, loading } = useFetch<SummaryResponse>(
    id ? "/api/summary/" + id : null
  );

  if (loading)
    return (
      <div className="flex-1 p-5 animate-pulse space-y-2">
        <div className="max-w-[300px] w-full h-5 bg-zinc-500 animate-pulse"></div>
        <div className="w-full md:max-w-[60%] h-[45%] bg-zinc-500 animate-pulse"></div>
      </div>
    );
  if (error) return <div>Error: {error}</div>;
  if (!data) return <div>No data</div>;

  console.log(data.summary.week_number);

  return (
    <div className="flex-1 p-5">
      <h1 className="font-bold text-xl">
        Summary for week: {data.summary.week_number}
      </h1>
      <div className="mt-3">
        <ReactMarkdown>{data.summary.summary}</ReactMarkdown>
      </div>
    </div>
  );
}
