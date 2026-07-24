import Navbar from "../components/Navbar";
import Transactions from "../components/Transactions";

export default function TransactionsPage() {
  return (
    <div className="min-h-screen bg-slate-950 text-white flex flex-col">
      <Navbar />

      <main className="flex-1 w-full max-w-7xl mx-auto p-4 md:p-8 space-y-8">
        <div>
          <h1 className="text-2xl font-bold text-white">Transactions</h1>
          <p className="text-sm text-slate-400 mt-1">
            Manage your paychecks and donations.
          </p>
        </div>

        <Transactions />
      </main>
    </div>
  );
}
