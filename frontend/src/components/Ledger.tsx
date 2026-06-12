import { useState, useEffect } from "react";
import { useToast } from "./Toast";
import { useAuth } from "../auth/AuthContext";
import { formatCurrency } from "../utils/formatter";

type Entry = {
  id: string;
  amount: number;
  ledger_entry: "paycheck" | "donation";
  description: string;
  transaction_date: string;
};

export default function Ledger() {
  const [entries, setEntries] = useState<Entry[]>([]);
  const [amount, setAmount] = useState("");
  const [type, setType] = useState<"paycheck" | "donation">("paycheck");
  const [description, setDescription] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);

  const showToast = useToast();
  const { checkAuth } = useAuth();

  // Fetch past entries
  const fetchEntries = async () => {
    try {
      const res = await fetch("http://localhost:1234/entries", {
        credentials: "include",
      });
      if (res.ok) {
        if (res.status === 204) {
          setEntries([]);
          return;
        }

        // Safely parse the text
        const text = await res.text();
        const data = text ? JSON.parse(text) : [];
        setEntries(data || []);
      }
    } catch (err) {
      console.error("Failed to fetch entries", err);
    }
  };

  useEffect(() => {
    fetchEntries();
  }, []);

  // Handle adding a new entry
  const handleSubmit = async (e: { preventDefault: () => void }) => {
    e.preventDefault();
    setIsSubmitting(true);

    try {
      const res = await fetch("http://localhost:1234/entries", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          amount: parseFloat(amount),
          ledger_entry: type,
          description: description,
        }),
        credentials: "include",
      });

      if (res.ok) {
        showToast("Entry added successfully!", "success");
        setAmount("");
        setDescription("");

        // Refresh the table and the global dashboard numbers!
        await fetchEntries();
        await checkAuth(true);
      } else {
        showToast("Failed to add entry.", "error");
      }
    } catch (err) {
      showToast("Network error.", "error");
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="grid grid-cols-1 lg:grid-cols-3 gap-8 mt-8">
      {/* LEFT COLUMN: Add Entry Form */}
      <div className="bg-slate-900 border border-slate-800 p-6 rounded-xl shadow-sm h-fit">
        <h2 className="text-lg font-semibold text-white mb-4">
          Log New Transaction
        </h2>
        <form onSubmit={handleSubmit} className="space-y-4">
          {/* Type Toggle */}
          <div className="flex rounded-lg bg-slate-950 p-1 border border-slate-800">
            <button
              type="button"
              onClick={() => setType("paycheck")}
              className={`flex-1 py-2 text-sm font-medium rounded-md transition-all ${
                type === "paycheck"
                  ? "bg-emerald-600 text-white"
                  : "text-slate-400 hover:text-white"
              }`}
            >
              Paycheck
            </button>
            <button
              type="button"
              onClick={() => setType("donation")}
              className={`flex-1 py-2 text-sm font-medium rounded-md transition-all ${
                type === "donation"
                  ? "bg-indigo-600 text-white"
                  : "text-slate-400 hover:text-white"
              }`}
            >
              Donation
            </button>
          </div>

          <div>
            <label className="block text-sm font-medium text-slate-400 mb-1">
              Amount ($)
            </label>
            <input
              type="number"
              step="0.01"
              required
              value={amount}
              onChange={(e) => setAmount(e.target.value)}
              className="w-full bg-slate-950 border border-slate-700 rounded-lg p-2.5 text-white focus:ring-2 focus:ring-indigo-500 outline-none"
              placeholder="0.00"
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-slate-400 mb-1">
              Description (Optional)
            </label>
            <input
              type="text"
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              className="w-full bg-slate-950 border border-slate-700 rounded-lg p-2.5 text-white focus:ring-2 focus:ring-indigo-500 outline-none"
              placeholder="e.g., May Salary, Shul Donation"
            />
          </div>

          <button
            type="submit"
            disabled={isSubmitting}
            className="w-full bg-indigo-600 hover:bg-indigo-500 text-white font-bold py-2.5 rounded-lg transition-colors disabled:opacity-50"
          >
            {isSubmitting ? "Saving..." : "Save Entry"}
          </button>
        </form>
      </div>

      {/* RIGHT COLUMN: Recent History */}
      <div className="lg:col-span-2 bg-slate-900 border border-slate-800 p-6 rounded-xl shadow-sm">
        <h2 className="text-lg font-semibold text-white mb-4">
          Recent History
        </h2>

        {entries.length === 0 ? (
          <div className="text-center py-10 text-slate-500">
            No entries found. Log a paycheck or donation to get started!
          </div>
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full text-left text-sm">
              <thead>
                <tr className="border-b border-slate-800 text-slate-400">
                  <th className="pb-3 font-medium">Date</th>
                  <th className="pb-3 font-medium">Description</th>
                  <th className="pb-3 font-medium text-right">Amount</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-slate-800/50">
                {entries.map((entry, index) => (
                  <tr
                    key={index}
                    className="text-slate-300 hover:bg-slate-800/20 transition-colors"
                  >
                    <td className="py-4">
                      {new Date(entry.transaction_date).toLocaleDateString()}
                    </td>
                    <td className="py-4">
                      <span className="flex items-center gap-2">
                        <span
                          className={`w-2 h-2 rounded-full ${entry.ledger_entry === "paycheck" ? "bg-emerald-500" : "bg-indigo-500"}`}
                        ></span>
                        {entry.description ||
                          (entry.ledger_entry === "paycheck"
                            ? "Income"
                            : "Donation")}
                      </span>
                    </td>
                    <td
                      className={`py-4 text-right font-medium ${entry.ledger_entry === "paycheck" ? "text-emerald-400" : "text-indigo-400"}`}
                    >
                      {entry.ledger_entry === "paycheck" ? "+" : "-"}
                      {formatCurrency(entry.amount)}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </div>
    </div>
  );
}
