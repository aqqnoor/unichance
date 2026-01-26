import { useEffect, useState } from "react";
import { useParams, Link } from "react-router-dom";
import { apiGet } from "../lib/api";

type UniversityLink = {
  id: string;
  link_type: string;
  url: string;
  title?: string | null;
  is_official: boolean;
  priority: number;
  source_code?: string | null;
  last_verified_at?: string | null;
};

type ProgramLite = {
  id: string;
  title: string;
  degree_level: string;
  field: string;
  language: string;
  tuition_amount?: number;
  tuition_currency?: string | null;
  has_scholarship: boolean;
};

type UniversityDTO = {
  id: string;
  name: string;
  country_code: string;
  city?: string | null;
  website?: string | null;
  qs_rank?: number | null;
  the_rank?: number | null;
  data_updated_at?: string | null;
  links: UniversityLink[];
  programs: ProgramLite[];
};

function niceType(t: string) {
  const map: Record<string, string> = {
    website: "Website",
    admissions: "Admissions",
    scholarships: "Scholarships",
    requirements: "Requirements",
    ranking: "Ranking",
    news: "News",
  };
  return map[t] || t;
}

export default function University() {
  const { id } = useParams();
  const [data, setData] = useState<UniversityDTO | null>(null);
  const [loading, setLoading] = useState(true);
  const [err, setErr] = useState("");

  useEffect(() => {
    if (!id) return;
    setLoading(true);
    setErr("");

    apiGet<UniversityDTO>(`/universities/${id}`)
      .then(setData)
      .catch((e) => setErr(e instanceof Error ? e.message : String(e)))
      .finally(() => setLoading(false));
  }, [id]);

  if (loading) {
    return <div className="max-w-5xl mx-auto px-4 py-8">Loading...</div>;
  }
  if (err) {
    return (
      <div className="max-w-5xl mx-auto px-4 py-8">
        <div className="p-3 rounded border border-red-300 bg-red-50 text-red-700">
          {err}
        </div>
      </div>
    );
  }
  if (!data) return null;

  return (
    <div className="max-w-5xl mx-auto px-4 py-8">
      <div className="mb-6">
        <Link to="/search" className="text-primary-600 hover:underline">
          ← Back to Search
        </Link>
      </div>

      <div className="card mb-6">
        <h1 className="text-3xl font-bold text-gray-900">{data.name}</h1>
        <p className="text-gray-600 mt-2">
          {(data.city ? `${data.city}, ` : "")}
          {data.country_code}
        </p>

        <div className="mt-3 flex flex-wrap gap-3 text-sm text-gray-700">
          {data.website && (
            <a className="text-primary-600 hover:underline" href={data.website} target="_blank" rel="noreferrer">
              Official site
            </a>
          )}
          {typeof data.qs_rank === "number" && data.qs_rank > 0 && (
            <span>QS: #{data.qs_rank}</span>
          )}
          {typeof data.the_rank === "number" && data.the_rank > 0 && (
            <span>THE: #{data.the_rank}</span>
          )}
        </div>
      </div>

      <div className="grid md:grid-cols-2 gap-6">
        {/* LINKS */}
        <div className="card">
          <h2 className="text-xl font-semibold mb-4">Useful links</h2>

          {data.links.length === 0 ? (
            <p className="text-gray-600">No links yet.</p>
          ) : (
            <div className="space-y-3">
              {data.links.map((l) => (
                <div key={l.id} className="p-3 rounded border border-gray-200">
                  <div className="flex items-center justify-between gap-3">
                    <div className="font-medium text-gray-900">
                      {niceType(l.link_type)}
                      {l.is_official ? (
                        <span className="ml-2 text-xs px-2 py-0.5 rounded-full bg-green-100 text-green-800">
                          official
                        </span>
                      ) : null}
                    </div>
                    <a
                      className="text-primary-600 hover:underline text-sm"
                      href={l.url}
                      target="_blank"
                      rel="noreferrer"
                    >
                      open
                    </a>
                  </div>

                  {l.title && <div className="text-sm text-gray-600 mt-1">{l.title}</div>}
                  {l.source_code && (
                    <div className="text-xs text-gray-500 mt-1">source: {l.source_code}</div>
                  )}
                </div>
              ))}
            </div>
          )}
        </div>

        {/* PROGRAMS */}
        <div className="card">
          <h2 className="text-xl font-semibold mb-4">Programs</h2>

          {data.programs.length === 0 ? (
            <p className="text-gray-600">No programs found.</p>
          ) : (
            <div className="space-y-3">
              {data.programs.map((p) => (
                <div key={p.id} className="p-3 rounded border border-gray-200">
                  <div className="font-semibold text-gray-900">{p.title}</div>
                  <div className="text-sm text-gray-600">
                    {p.degree_level} • {p.field} • {p.language}
                    {p.has_scholarship ? " • scholarship" : ""}
                  </div>
                  {typeof p.tuition_amount === "number" && (
                    <div className="text-sm text-gray-600 mt-1">
                      Tuition: {p.tuition_amount.toLocaleString()} {p.tuition_currency || ""}
                    </div>
                  )}
                </div>
              ))}
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
