import { useEffect, useState } from "react";
import { apiGet, apiPost } from "../lib/api";

type Program = {
  id: string;
  title: string;
  degree_level: string;
  field: string;
  language: string;
  university_name: string;
  country_code: string;
  tuition_amount?: number;
  tuition_currency?: string;
  has_scholarship: boolean;
};

export default function Programs() {
  const [items, setItems] = useState<Program[]>([]);
  const [page, setPage] = useState(1);
  const [total, setTotal] = useState(0);

  // filters
  const [q, setQ] = useState("");
  const [country, setCountry] = useState("");
  const [level, setLevel] = useState("");

  // shortlist ui state
  const [addingId, setAddingId] = useState<string | null>(null);
  const [added, setAdded] = useState<Record<string, boolean>>({});

  const limit = 10;
  const pages = Math.ceil(total / limit);

  useEffect(() => {
    const params = new URLSearchParams();
    params.set("page", String(page));
    params.set("limit", String(limit));

    if (q) params.set("q", q);
    if (country) params.set("country_code", country);
    if (level) params.set("degree_level", level);


    apiGet<{ items: Program[]; total: number }>(`/programs?${params.toString()}`)
        .then((res) => {
            setItems(res.items);
            setTotal(res.total);
        })
        .catch(console.error);
    }, [q, country, level, page]);

  async function addToShortlist(programId: string) {
    try {
      setAddingId(programId);
      await apiPost("/shortlist", { program_id: programId }); // backend expects program_id
      setAdded((prev) => ({ ...prev, [programId]: true }));
    } catch (e) {
      console.error(e);
      alert("Shortlist-ке қосу кезінде қате шықты. Login жасағаныңды тексер.");
    } finally {
      setAddingId(null);
    }
  }

  return (
    <div className="max-w-6xl mx-auto p-6">
      <h1 className="text-2xl font-semibold mb-4">Programs</h1>

      {/* Filters */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-3 mb-6">
        <input
          className="border rounded px-3 py-2"
          placeholder="Search program"
          value={q}
          onChange={(e) => {
            setPage(1);
            setQ(e.target.value);
          }}
        />

        <select
          className="border rounded px-3 py-2"
          value={country}
          onChange={(e) => {
            setPage(1);
            setCountry(e.target.value);
          }}
        >
          <option value="">All countries</option>
          <option value="DE">Germany</option>
          <option value="US">USA</option>
          <option value="GB">UK</option>
          <option value="FR">France</option>
          <option value="IT">Italy</option>
        </select>

        <select
          className="border rounded px-3 py-2"
          value={level}
          onChange={(e) => {
            setPage(1);
            setLevel(e.target.value);
          }}
        >
          <option value="">All levels</option>
          <option value="bachelor">Bachelor</option>
          <option value="master">Master</option>
        </select>
      </div>

      {/* Results */}
      <div className="space-y-3">
        {items.map((p) => (
          <div key={p.id} className="border rounded-lg p-4 hover:bg-gray-50">
            <div className="flex items-start justify-between gap-4">
              <div>
                <div className="font-medium">
                  {p.title} — {p.university_name}
                </div>
                <div className="text-sm text-gray-600">
                  {p.country_code} · {p.degree_level} · {p.field}
                </div>
                <div className="text-sm mt-1">
                  {p.tuition_amount === undefined || p.tuition_amount === null
                    ? "Tuition not specified"
                    : p.tuition_amount === 0
                        ? "Free"
                        : `${p.tuition_amount} ${p.tuition_currency}`}

                  {p.has_scholarship && " · Scholarship available"}
                </div>
              </div>

              <button
                className="px-3 py-2 border rounded disabled:opacity-50"
                disabled={addingId === p.id || added[p.id]}
                onClick={() => addToShortlist(p.id)}
              >
                {added[p.id]
                  ? "Added"
                  : addingId === p.id
                  ? "Adding..."
                  : "Add to shortlist"}
              </button>
            </div>
          </div>
        ))}
      </div>

      {/* Pagination */}
      {pages > 1 && (
        <div className="flex gap-2 mt-6">
          <button
            className="px-3 py-1 border rounded disabled:opacity-50"
            disabled={page === 1}
            onClick={() => setPage((p) => p - 1)}
          >
            Prev
          </button>
          <span className="px-2 py-1">
            {page} / {pages}
          </span>
          <button
            className="px-3 py-1 border rounded disabled:opacity-50"
            disabled={page === pages}
            onClick={() => setPage((p) => p + 1)}
          >
            Next
          </button>
        </div>
      )}
    </div>
  );
}
