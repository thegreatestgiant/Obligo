import Ledger from "../components/Ledger";
import Navbar from "../components/Navbar";
import Summary from "../components/Summary";

export default function Dashboard() {
  return (
    <div className="min-h-screen bg-slate-950 text-white">
      <Navbar />

      <main className="max-w-7xl mx-auto p-8 space-y-8">
        <div>
          <h1 className="text-2xl font-bold text-white">Dashboard</h1>
          <p className="text-sm text-slate-400 mt-1">
            Your Maaser calculation and donation summary.
          </p>
        </div>

        <Summary />
        <Ledger />
      </main>
    </div>
  );
}
