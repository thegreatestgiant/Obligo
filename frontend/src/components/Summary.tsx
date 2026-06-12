import { useAuth } from "../auth/AuthContext";
import Loading from "./Loading";
import { formatCurrency } from "../utils/formatter";
import GoalWidget from "./GoalWidget";
import IndexAnchor from "./IndexAnchor";
import KPICard from "./KPICard";

export default function Summary() {
  const { user, isLoading } = useAuth();

  if (isLoading) return <Loading />;

  // Prepare data safely
  const totalEarned = user?.Total_Earned ?? 0;
  const totalDonated = user?.Total_Donated ?? 0;
  const totalOwed = user?.Total_Owed ?? 0;
  const remaining = Math.max(0, totalOwed - totalDonated);
  const percentFulfilled = user?.Percent_Fulfilled?.toFixed(1) ?? "0.0";
  const percentOwed =
    totalOwed > 0 ? ((remaining / totalOwed) * 100).toFixed(1) : "0.0";

  // Mocked data for the trend chart (using the structure your chart component expects)
  const chartData = [
    {
      month: "Current",
      earned: totalEarned,
      donated: totalDonated,
      target: totalOwed,
    },
  ];

  return (
    <div className="space-y-12">
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

      {/* TIER 2: INTERACTIVE & TREND DATA */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
        <GoalWidget totalTarget={totalOwed} totalDonated={totalDonated} />
        <IndexAnchor data={chartData} title="Cumulative Cashflow Overview" />
      </div>
    </div>
  );
}
