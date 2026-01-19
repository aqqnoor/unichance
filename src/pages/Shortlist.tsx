import { useEffect, useState } from "react";
import { apiGet, apiPost } from "../lib/api";


type Item = {
  program_id: string;
  title: string;
  university: string;
  country: string;
  degree: string;
  field: string;
  tuition_amount?: number;
  tuition_currency?: string;
  has_scholarship: boolean;
  qs_rank?: number;
  the_rank?: number;
};

export default function Shortlist() {
  const [items, setItems] = useState<Item[]>([]);

  const load = () =>
    apiGet<{ items: Item[] }>("/shortlist").then(r => setItems(r.items));

  useEffect(() => { load(); }, []);

  return (
    <div className="max-w-6xl mx-auto p-6">
      <h1 className="text-2xl font-semibold mb-4">Shortlist</h1>

      <div className="overflow-auto">
        <table className="w-full border">
          <thead>
            <tr className="bg-gray-100">
              <th className="p-2 border">Program</th>
              <th className="p-2 border">University</th>
              <th className="p-2 border">Tuition</th>
              <th className="p-2 border">Scholarship</th>
              <th className="p-2 border">QS</th>
              <th className="p-2 border"></th>
            </tr>
          </thead>
          <tbody>
            {items.map(i => (
              <tr key={i.program_id}>
                <td className="p-2 border">{i.title}</td>
                <td className="p-2 border">{i.university}</td>
                <td className="p-2 border">
                  {i.tuition_amount ? `${i.tuition_amount} ${i.tuition_currency}` : "—"}
                </td>
                <td className="p-2 border">{i.has_scholarship ? "Yes" : "No"}</td>
                <td className="p-2 border">{i.qs_rank ?? "—"}</td>
                <td className="p-2 border">
                  <button
                    className="text-red-600"
                    onClick={() => apiPost(`/shortlist/${i.program_id}`, {})}
                  >
                    Remove
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}
