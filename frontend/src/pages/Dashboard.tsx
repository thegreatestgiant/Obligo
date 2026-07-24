
import Navbar from "../components/Navbar";
import Summary from "../components/Summary";
import { Link } from "react-router-dom";

export default function Dashboard() {
  return (
    <div className="min-h-screen bg-slate-950 text-white flex flex-col">
      <Navbar />

      <main className="flex-1 w-full max-w-7xl mx-auto p-4 md:p-8 space-y-8">
        <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
          <div>
            <h1 className="text-2xl font-bold text-white">Dashboard</h1>
            <p className="text-sm text-slate-400 mt-1">
              Your Maaser calculation and donation summary.
            </p>
          </div>
          <Link
            to="/transactions"
            className="px-5 py-2.5 bg-indigo-600 hover:bg-indigo-500 text-white text-sm font-bold rounded-lg shadow-lg shadow-indigo-500/20 transition-all text-center"
          >
            Manage Transactions
          </Link>
        </div>

        <Summary />
      </main>
    </div>
  );
}
