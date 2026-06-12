import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ResponsiveContainer,
} from "recharts";
import { useAuth } from "../auth/AuthContext";
import Loading from "./Loading";

export default function Summary() {
  const { user, isLoading } = useAuth();

  if (isLoading) {
    return <Loading />;
  }

  const chartData = [
    {
      name: "Financial Overview",
      Earned: user?.Total_Earned ?? 0,
      Obligation: user?.Total_Owed ?? 0,
      Donated: user?.Total_Donated ?? 0,
    },
  ];

  return (
    <div className="space-y-8">
      {/* Recharts Section */}
      <div className="bg-slate-900 border border-slate-800 p-6 rounded-xl shadow-sm min-h-[400px] mb-8">
        <h2 className="text-lg font-semibold mb-6">Financial Comparison</h2>

        <ResponsiveContainer width="100%" height={350}>
          <BarChart
            data={chartData}
            margin={{ top: 5, right: 30, left: 20, bottom: 5 }}
          >
            <CartesianGrid strokeDasharray="3 3" stroke="#334155" />
            <XAxis dataKey="name" stroke="#94a3b8" />
            <YAxis stroke="#94a3b8" />
            <Tooltip
              cursor={{ fill: "#1e293b" }}
              contentStyle={{
                backgroundColor: "#0f172a",
                borderColor: "#1e293b",
                color: "#fff",
                borderRadius: "8px",
              }}
            />
            <Legend wrapperStyle={{ paddingTop: "20px" }} />
            <Bar dataKey="Earned" fill="#10b981" radius={[4, 4, 0, 0]} />
            <Bar dataKey="Obligation" fill="#ef4444" radius={[4, 4, 0, 0]} />
            <Bar dataKey="Donated" fill="#6366f1" radius={[4, 4, 0, 0]} />
          </BarChart>
        </ResponsiveContainer>
      </div>

      {/* Progress Bar Section */}
      <div className="bg-slate-900 border border-slate-800 p-6 rounded-xl shadow-sm mb-8">
        <div className="flex justify-between items-end mb-2">
          <div>
            <p className="text-sm font-medium text-slate-400">
              Maaser Goal Progress
            </p>
            <p className="text-xl font-semibold text-white mt-1">
              {user?.Percent_Fulfilled?.toFixed(1) ?? "0.0"}% Fulfilled
            </p>
          </div>
          <p className="text-sm font-medium text-slate-500">
            Total Target: $
            {user?.Total_Owed?.toLocaleString(undefined, {
              minimumFractionDigits: 2,
              maximumFractionDigits: 2,
            }) ?? "0.00"}
          </p>
        </div>

        <div className="w-full bg-slate-950 rounded-full h-3 border border-slate-800 overflow-hidden mt-3">
          <div
            className="bg-indigo-500 h-3 rounded-full transition-all duration-1000 ease-out"
            style={{
              width: `${Math.min(user?.Percent_Fulfilled ?? 0, 100)}%`,
            }}
          ></div>
        </div>
      </div>
    </div>
  );
}
