import { useState, useEffect } from "react";
import { useToast } from "./Toast";
import { useAuth } from "../auth/AuthContext";
import { formatCurrency } from "../utils/formatter";
import { api } from "../api/clients";

type Entry = {
  transaction_id: string;
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
  const [currentPage, setCurrentPage] = useState(1);
  const itemsPerPage = 8;
  const showToast = useToast();
  const { checkAuth } = useAuth();

  const handleDelete = async (id: string) => {
    try {
      const res = await api.deleteEntry(id);
      if (res.ok) {
        showToast("Transaction deleted.", "success");
        // setCurrentPage(1);
        await fetchEntries(); // Refresh data
        await checkAuth(true); // Sync global dashboard numbers
      } else {
        showToast("Failed to delete.", "error");
      }
    } catch (err: any) {
      if (err.message === "Unauthorized") return;
      showToast("Network error.", "error");
    }
  };

  // Fetch past entries
  const fetchEntries = async () => {
    const res = await api.getEntries();
    if (res.ok) {
      if (res.status === 204) {
        setEntries([]);
        return;
      }
      const text = await res.text();
      const data = text ? JSON.parse(text) : [];
      setEntries(data || []);
    }
  };

  useEffect(() => {
    fetchEntries();
  }, []);

  useEffect(() => {
    const maxPages = Math.max(1, Math.ceil(entries.length / itemsPerPage));
    if (currentPage > maxPages) {
      setCurrentPage(maxPages);
    }
  }, [entries.length, currentPage, itemsPerPage]);
  const handleSubmit = async (e: { preventDefault: () => void }) => {
    e.preventDefault();
    setIsSubmitting(true);

    try {
      const res = await api.createEntry(parseFloat(amount), type, description);

      if (res.ok) {
        showToast("Entry added successfully!", "success");
        setAmount("");
        setDescription("");
        setCurrentPage(1);

        // Refresh the table and the global dashboard numbers!
        await fetchEntries();
        await checkAuth(true);
      } else {
        showToast("Failed to add entry.", "error");
      }
    } catch (err: any) {
      if (err.message === "Unauthorized") return;
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

      {/* RIGHT COLUMN: Recent History (Paginated) */}
      <div className="lg:col-span-2 bg-slate-900 border border-slate-800 p-6 rounded-xl shadow-sm">
        <div className="flex justify-between items-center mb-4">
          <h2 className="text-lg font-semibold text-white">Recent History</h2>
          <div className="flex items-center gap-3">
            <span className="text-xs font-bold text-slate-500 uppercase tracking-widest">
              Page {currentPage} /{" "}
              {Math.max(1, Math.ceil(entries.length / itemsPerPage))}
            </span>
          </div>
        </div>

        {entries.length === 0 ? (
          <div className="text-center py-10 text-slate-500">
            No entries found. Log a paycheck or donation to get started!
          </div>
        ) : (
          <>
            <div className="overflow-x-auto min-h-[400px]">
              <table className="w-full text-left text-sm">
                <thead>
                  <tr className="border-b border-slate-800 text-slate-400">
                    <th className="pb-3 font-medium">Date</th>
                    <th className="pb-3 font-medium">Description</th>
                    <th className="pb-3 font-medium text-right">Amount</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-slate-800/50">
                  {entries
                    .slice(
                      (currentPage - 1) * itemsPerPage,
                      currentPage * itemsPerPage,
                    )
                    .map((entry) => (
                      <tr
                        key={entry.transaction_id}
                        className="group text-slate-300 hover:bg-slate-800/20 transition-colors"
                      >
                        <td className="py-4">
                          {new Date(
                            entry.transaction_date,
                          ).toLocaleDateString()}
                        </td>
                        <td className="py-4">
                          <span className="flex items-center gap-2">
                            <span
                              className={`w-2 h-2 rounded-full ${
                                entry.ledger_entry === "paycheck"
                                  ? "bg-emerald-500"
                                  : "bg-indigo-500"
                              }`}
                            ></span>
                            {entry.description ||
                              (entry.ledger_entry === "paycheck"
                                ? "Income"
                                : "Donation")}
                          </span>
                        </td>
                        <td className="py-4 text-right flex items-center justify-end gap-4">
                          <span
                            className={
                              entry.ledger_entry === "paycheck"
                                ? "text-emerald-400"
                                : "text-indigo-400"
                            }
                          >
                            {entry.ledger_entry === "paycheck" ? "+" : "-"}
                            {formatCurrency(entry.amount)}
                          </span>

                          {/* Trash Can Icon */}
                          <button
                            onClick={() => handleDelete(entry.transaction_id)}
                            className={`text-slate-600 hover:text-red-500 transition-opacity opacity-0 group-hover:opacity-100`}
                            title="Warning: Deleting and recreating an entry changes its historical date to today, which may alter past monthly summaries and recalculate owed amounts using your current donation percentage."
                          >
                            <svg
                              xmlns="http://www.w3.org/2000/svg"
                              className="h-4 w-4"
                              fill="none"
                              viewBox="0 0 24 24"
                              stroke="currentColor"
                            >
                              <path
                                strokeLinecap="round"
                                strokeLinejoin="round"
                                strokeWidth={2}
                                d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
                              />
                            </svg>
                          </button>
                        </td>
                      </tr>
                    ))}
                </tbody>
              </table>
            </div>

            {/* Pagination Controls */}
            <div className="flex justify-between items-center mt-6 pt-4 border-t border-slate-800">
              <button
                onClick={() => setCurrentPage((prev) => Math.max(prev - 1, 1))}
                disabled={currentPage === 1}
                className="px-4 py-2 bg-slate-800 hover:bg-slate-700 text-slate-300 rounded-lg text-sm transition-all disabled:opacity-30"
              >
                Previous
              </button>
              <button
                onClick={() =>
                  setCurrentPage((prev) =>
                    Math.min(
                      prev + 1,
                      Math.ceil(entries.length / itemsPerPage),
                    ),
                  )
                }
                disabled={
                  currentPage >= Math.ceil(entries.length / itemsPerPage)
                }
                className="px-4 py-2 bg-slate-800 hover:bg-slate-700 text-slate-300 rounded-lg text-sm transition-all disabled:opacity-30"
              >
                Next
              </button>
            </div>
          </>
        )}
      </div>
    </div>
  );
}
