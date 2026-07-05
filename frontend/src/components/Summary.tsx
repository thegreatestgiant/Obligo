import { useAuth } from "../auth/AuthContext";
import Loading from "./Loading";
import { formatCurrency } from "../utils/formatter";
import GoalWidget from "./GoalWidget";
import IndexAnchor from "./IndexAnchor";
import KPICard from "./KPICard";
import { useEffect, useState } from "react";
import { api } from "../api/clients";

type Monthly = {
  month: string;
  donated: number;
  earned: number;
  target: number;
};

export default function Summary() {
  const { user, isLoading } = useAuth();

  const [monthlies, setMonthlies] = useState<Monthly[]>([]);

  if (isLoading) return <Loading />;

  // Prepare data safely
  const totalEarned = user?.Total_Earned ?? 0;
  const totalDonated = user?.Total_Donated ?? 0;
  const totalOwed = user?.Total_Owed ?? 0;
  const remaining = Math.max(0, totalOwed - totalDonated);
  const percentFulfilled = user?.Percent_Fulfilled?.toFixed(1) ?? "0.0";
  const percentOwed =
    totalOwed > 0 ? ((remaining / totalOwed) * 100).toFixed(1) : "0.0";

  const getMonthly = async () => {
    const res = await api.getMonthlySummary();
    if (res.ok) {
      if (res.status === 204) {
        setMonthlies([]);
        return;
      }
      const text = await res.text();
      const data = text ? JSON.parse(text) : [];
      setMonthlies(data || []);
    }
  };

  useEffect(() => {
    // 1. Fetch data initially when the component loads
    getMonthly();

    // 2. Define what happens when we hear the event
    const handleLedgerUpdate = () => {
      getMonthly();
    };

    // 3. Start listening for the custom event we created in Ledger.tsx
    window.addEventListener("ledger-updated", handleLedgerUpdate);

    // 4. Clean up the listener if the user navigates away from the component
    return () => {
      window.removeEventListener("ledger-updated", handleLedgerUpdate);
    };
  }, []); // The empty array ensures this listener is only attached once

  return (
    <div className="space-y-12">
      <div className="flex items-center gap-2 mb-6">
        {user && (
          <span className="px-3 py-1 text-xs font-bold uppercase tracking-wider bg-indigo-500/10 text-indigo-400 border border-indigo-500/20 rounded-full">
            Current Goal: {user.Donation_Percent}% Maaser
          </span>
        )}
      </div>
      {/* TIER 1: PUNK KPI ROW */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        <KPICard
          intent="earned"
          title="Total Earned"
          value={formatCurrency(totalEarned)}
        />
        <KPICard
          intent="donated"
          title="Total Donated"
          value={formatCurrency(totalDonated)}
          percent={percentFulfilled + "%"}
        />
        <KPICard
          intent="obligation"
          title="Obligation"
          value={formatCurrency(totalOwed)}
          subtext={`${formatCurrency(remaining)} remaining`}
          percent={percentOwed + "%"}
        />
      </div>
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-8 min-h-0">
        <GoalWidget totalTarget={totalOwed} totalDonated={totalDonated} />
        {/* The data flows natively into the chart right here */}
        <IndexAnchor data={monthlies} title="Cumulative Cashflow Overview" />
      </div>
    </div>
  );
}
