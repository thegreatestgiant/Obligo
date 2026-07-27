import { useState, useEffect, useRef } from "react";
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

export default function Transactions() {
  const [entries, setEntries] = useState<Entry[]>([]);
  const [amount, setAmount] = useState("");
  const [type, setType] = useState<"paycheck" | "donation">("paycheck");
  const [description, setDescription] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [originalAmount, setOriginalAmount] = useState<number | null>(null);
  const [originalDescription, setOriginalDescription] = useState<string | null>(
    null,
  );
  const [currentPage, setCurrentPage] = useState(1);
  const itemsPerPage = 8;
  const showToast = useToast();
  const { checkAuth } = useAuth();

  const [editingId, setEditingId] = useState<string | null>(null);
  const editingEntry = entries.find((e) => e.transaction_id === editingId);

  const fileInputRef = useRef<HTMLInputElement>(null);
  const formRef = useRef<HTMLDivElement>(null);
  const amountInputRef = useRef<HTMLInputElement>(null);

  // State for the Delete Confirmation Modal
  const [deleteModalOpen, setDeleteModalOpen] = useState(false);
  const [entryToDelete, setEntryToDelete] = useState<string | null>(null);
  const [dontShowAgain, setDontShowAgain] = useState(false);

  const initiateDelete = (id: string) => {
    const skipWarning = sessionStorage.getItem("skipDeleteWarning") === "true";
    if (skipWarning) {
      executeDelete(id);
    } else {
      setEntryToDelete(id);
      setDeleteModalOpen(true);
    }
  };

  const confirmDelete = () => {
    if (dontShowAgain) {
      sessionStorage.setItem("skipDeleteWarning", "true");
    }
    if (entryToDelete) {
      executeDelete(entryToDelete);
    }
  };

  const cancelDelete = () => {
    setDeleteModalOpen(false);
    setEntryToDelete(null);
    setDontShowAgain(false);
  };

  const executeDelete = async (id: string) => {
    try {
      const res = await api.deleteEntry(id);
      if (res.ok) {
        showToast("Transaction deleted.", "success");
        await fetchEntries();
        await checkAuth(true);
        window.dispatchEvent(new Event("transactions-updated"));

        if (editingId === id) resetForm();
      } else {
        showToast("Failed to delete.", "error");
      }
    } catch (err: any) {
      if (err.message === "Unauthorized") return;
      showToast("Network error.", "error");
    } finally {
      setDeleteModalOpen(false);
      setEntryToDelete(null);
    }
  };

  const handleEditClick = (entry: Entry) => {
    setEditingId(entry.transaction_id);
    setAmount(entry.amount.toString());
    setDescription(entry.description);
    setOriginalAmount(entry.amount);
    setOriginalDescription(entry.description);
    setType(entry.ledger_entry);

    setTimeout(() => {
      formRef.current?.scrollIntoView({ behavior: "smooth", block: "nearest" });
      amountInputRef.current?.focus();
    }, 50);
  };

  const resetForm = () => {
    setEditingId(null);
    setAmount("");
    setOriginalAmount(null);
    setDescription("");
    setOriginalDescription(null);
    setType("paycheck");
  };

  const handleExport = async () => {
    try {
      const res = await api.exportCSV();
      const blob = await res.blob();
      const url = window.URL.createObjectURL(blob);
      const a = document.createElement("a");
      a.href = url;
      a.download = "ledger.csv";
      document.body.appendChild(a);
      a.click();
      a.remove();
      window.URL.revokeObjectURL(url);
    } catch (err) {
      showToast("Failed to export data.", "error");
    }
  };

  const handleImport = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;

    try {
      const res = await api.importCSV(file);
      const data = await res.json();
      showToast(
        `Imported ${data.inserted} rows. (Skipped ${data.skipped} duplicates)`,
        "success",
      );
      await fetchEntries();
      await checkAuth(true);
      window.dispatchEvent(new Event("transactions-updated"));
    } catch (err) {
      showToast("Failed to import data. Check file format.", "error");
    } finally {
      if (fileInputRef.current) fileInputRef.current.value = "";
    }
  };

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

  // Listen for the Escape key to exit edit mode
  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      if (e.key === "Escape" && editingId) {
        resetForm();
      }
    };

    window.addEventListener("keydown", handleKeyDown);
    return () => window.removeEventListener("keydown", handleKeyDown);
  }, [editingId]);

  useEffect(() => {
    fetchEntries();
  }, []);

  useEffect(() => {
    const maxPages = Math.max(1, Math.ceil(entries.length / itemsPerPage));
    if (currentPage > maxPages) {
      setCurrentPage(maxPages);
    }
  }, [entries.length, currentPage, itemsPerPage]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsSubmitting(true);

    try {
      let res;
      if (editingId) {
        const changedAmount =
          parseFloat(amount) !== originalAmount
            ? parseFloat(amount)
            : undefined;
        const changedDescription =
          description !== originalDescription ? description : undefined;

        res = await api.editEntry(editingId, changedAmount, changedDescription);
      } else {
        res = await api.createEntry(parseFloat(amount), type, description);
      }

      if (res.ok) {
        showToast(
          editingId ? "Entry updated!" : "Entry added successfully!",
          "success",
        );

        if (!editingId) {
          setCurrentPage(1);
        }

        resetForm();
        await fetchEntries();
        await checkAuth(true);
        window.dispatchEvent(new Event("transactions-updated"));
      } else {
        showToast("Failed to save entry.", "error");
      }
    } catch (err: any) {
      if (err.message === "Unauthorized") return;
      showToast("Network error.", "error");
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <>
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-8 mt-8">
        {/* LEFT COLUMN: Add Entry Form */}
        <div
          ref={formRef}
          className={`border p-6 rounded-xl shadow-sm h-fit transition-all duration-300 ${
            editingId
              ? "bg-slate-900 border-amber-500/50 shadow-lg shadow-amber-900/10 ring-1 ring-amber-500/20"
              : "bg-slate-900 border-slate-800"
          }`}
        >
          <div className="flex justify-between items-center mb-4">
            <h2 className="text-lg font-semibold text-white flex items-center gap-2">
              {editingId ? "Edit Transaction" : "Log New Transaction"}
              {/* Badge showing date of edited entry */}
              {editingEntry && (
                <span className="text-[11px] font-medium px-2 py-0.5 bg-amber-500/10 text-amber-400 rounded-md flex items-center border border-amber-500/20">
                  <span className="w-1.5 h-1.5 rounded-full bg-amber-500 mr-1.5 animate-pulse"></span>
                  {new Date(editingEntry.transaction_date).toLocaleDateString()}
                </span>
              )}
            </h2>
            {editingId && (
              <button
                type="button"
                onClick={resetForm}
                className="text-sm text-slate-400 hover:text-white transition-colors"
                title="Press Esc to cancel"
              >
                Cancel
              </button>
            )}
          </div>
          <form onSubmit={handleSubmit} className="space-y-4">
            {/* Type Toggle */}
            <div className="flex rounded-lg bg-slate-950 p-1 border border-slate-800">
              <button
                type="button"
                disabled={!!editingId}
                onClick={() => setType("paycheck")}
                className={`flex-1 py-2 text-sm font-medium rounded-md transition-all ${
                  type === "paycheck"
                    ? "bg-emerald-600 text-white"
                    : "text-slate-400 hover:text-white"
                } disabled:opacity-60 disabled:cursor-not-allowed`}
              >
                Paycheck
              </button>
              <button
                type="button"
                disabled={!!editingId}
                onClick={() => setType("donation")}
                className={`flex-1 py-2 text-sm font-medium rounded-md transition-all ${
                  type === "donation"
                    ? "bg-indigo-600 text-white"
                    : "text-slate-400 hover:text-white"
                } disabled:opacity-60 disabled:cursor-not-allowed`}
              >
                Donation
              </button>
            </div>

            <div>
              <label className="block text-sm font-medium text-slate-400 mb-1">
                Amount ($)
              </label>
              <input
                ref={amountInputRef}
                type="number"
                step="0.01"
                required
                value={amount}
                onChange={(e) => setAmount(e.target.value)}
                className="w-full bg-slate-950 border border-slate-700 rounded-lg p-2.5 text-white focus:ring-2 focus:ring-indigo-500 outline-none transition-shadow"
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
                className="w-full bg-slate-950 border border-slate-700 rounded-lg p-2.5 text-white focus:ring-2 focus:ring-indigo-500 outline-none transition-shadow"
                placeholder="e.g., May Salary, Shul Donation"
              />
            </div>

            <button
              type="submit"
              disabled={isSubmitting}
              className={`w-full font-bold py-2.5 rounded-lg transition-colors disabled:opacity-50 text-white ${
                editingId
                  ? "bg-amber-600 hover:bg-amber-500"
                  : "bg-indigo-600 hover:bg-indigo-500"
              }`}
            >
              {isSubmitting
                ? "Saving..."
                : editingId
                  ? "Save Changes"
                  : "Save Entry"}
            </button>
          </form>
        </div>

        {/* RIGHT COLUMN: Recent History (Paginated) */}
        <div className="lg:col-span-2 bg-slate-900 border border-slate-800 p-6 rounded-xl shadow-sm">
          <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center mb-4 gap-4">
            <h2 className="text-lg font-semibold text-white">Recent History</h2>

            <div className="flex items-center gap-3">
              <input
                type="file"
                accept=".csv"
                ref={fileInputRef}
                onChange={handleImport}
                className="hidden"
              />
              <button
                onClick={() => fileInputRef.current?.click()}
                className="px-3 py-1.5 text-xs font-semibold bg-slate-800 hover:bg-slate-700 text-slate-300 rounded-md transition-colors"
              >
                Import CSV
              </button>
              <button
                onClick={handleExport}
                className="px-3 py-1.5 text-xs font-semibold bg-slate-800 hover:bg-slate-700 text-slate-300 rounded-md transition-colors"
              >
                Export CSV
              </button>
              <span className="text-xs font-bold text-slate-500 uppercase tracking-widest border-l border-slate-700 pl-3">
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
                      <th className="pb-3 pl-2 font-medium">Date</th>
                      <th className="pb-3 px-4 font-medium">Description</th>
                      <th className="pb-3 pl-4 pr-0 font-medium text-right">Amount</th>
                      <th className="pb-3 pl-2 pr-2 w-[70px] font-medium text-right text-slate-500 text-[10px] uppercase tracking-wider">Tools</th>
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
                          className={`group transition-colors ${
                            editingId === entry.transaction_id
                              ? "bg-slate-800/80 text-white"
                              : "text-slate-300 hover:bg-slate-800/20"
                          }`}
                        >
                          <td className="py-4 pl-2 whitespace-nowrap">
                            {new Date(
                              entry.transaction_date,
                            ).toLocaleDateString()}
                          </td>
                          <td className="py-4 px-4 min-w-[120px]">
                            <span className="flex items-center gap-2">
                              <span
                                className={`w-2 h-2 rounded-full shrink-0 ${
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
                          <td className="py-4 pl-4 pr-0 text-right whitespace-nowrap">
                            <span
                              className={
                                entry.ledger_entry === "paycheck"
                                  ? "text-emerald-400 font-medium"
                                  : "text-indigo-400 font-medium"
                              }
                            >
                              {entry.ledger_entry === "paycheck" ? "+" : "-"}
                              {formatCurrency(entry.amount)}
                            </span>
                          </td>
                          <td className="py-4 pr-2 pl-2 text-right w-[70px]">
                            <div className="flex items-center justify-end gap-2">
                              {/* Render explicit "EDITING" badge or the standard action buttons */}
                              {editingId === entry.transaction_id ? (
                                <span className="text-[10px] font-bold text-amber-500 tracking-widest bg-amber-500/10 px-2 py-1 rounded">
                                  EDITING
                                </span>
                              ) : (
                                <>
                                  {/* Edit Icon */}
                                  <button
                                    onClick={() => handleEditClick(entry)}
                                    className="text-slate-600 hover:text-amber-500 transition-opacity opacity-100 md:opacity-0 group-hover:opacity-100"
                                    title="Edit Entry"
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
                                        d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"
                                      />
                                    </svg>
                                  </button>
  
                                  {/* Trash Can Icon */}
                                  <button
                                    onClick={() =>
                                      initiateDelete(entry.transaction_id)
                                    }
                                    className="text-slate-600 hover:text-red-500 transition-opacity opacity-100 md:opacity-0 group-hover:opacity-100"
                                    title="Delete Transaction"
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
                                </>
                              )}
                            </div>
                          </td>
                        </tr>
                      ))}
                  </tbody>
                </table>
              </div>

              {/* Pagination Controls */}
              <div className="flex justify-between items-center mt-6 pt-4 border-t border-slate-800">
                <button
                  onClick={() =>
                    setCurrentPage((prev) => Math.max(prev - 1, 1))
                  }
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

      {/* Delete Confirmation Modal */}
      {deleteModalOpen && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm px-4">
          <div className="bg-slate-900 border border-slate-700 rounded-xl p-6 max-w-md w-full shadow-2xl">
            <h3 className="text-xl font-bold text-white mb-2">
              Delete Transaction?
            </h3>
            <p className="text-sm text-slate-400 mb-4">
              Warning: Deleting and recreating an entry changes its historical
              date to today, which may alter past monthly summaries and
              recalculate owed amounts using your current donation percentage.
            </p>

            <label className="flex items-center gap-2 mb-6 cursor-pointer group">
              <input
                type="checkbox"
                className="w-4 h-4 rounded border-slate-700 bg-slate-950 text-red-500 focus:ring-red-500 focus:ring-offset-slate-900"
                checked={dontShowAgain}
                onChange={(e) => setDontShowAgain(e.target.checked)}
              />
              <span className="text-sm text-slate-300 group-hover:text-white transition-colors">
                Don't show me this warning again this session
              </span>
            </label>

            <div className="flex justify-end gap-3">
              <button
                onClick={cancelDelete}
                className="px-4 py-2 rounded-lg text-sm font-medium text-slate-300 bg-slate-800 hover:bg-slate-700 transition-colors"
              >
                Cancel
              </button>
              <button
                onClick={confirmDelete}
                className="px-4 py-2 rounded-lg text-sm font-medium text-white bg-red-600 hover:bg-red-500 transition-colors"
              >
                Delete
              </button>
            </div>
          </div>
        </div>
      )}
    </>
  );
}
