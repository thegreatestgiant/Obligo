import {
  ComposedChart,
  Bar,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ResponsiveContainer,
} from "recharts";
import { formatCompactCurrency, formatCurrency } from "../utils/formatter";

interface ChartProps {
  data: any[];
  title: string;
}

export default function IndexAnchor({ data, title }: ChartProps) {
  // 1. Softened the tooltip for a more modern UX
  const customTooltip = {
    backgroundColor: "#0f172a", // Tailwind slate-900
    borderColor: "#334155", // Tailwind slate-700
    borderWidth: "1px",
    color: "#f8fafc",
    borderRadius: "8px",
    padding: "12px",
  };

  // 2. Helper function to turn "January 2025" into "Jan '25"
  const formatMonthAxis = (tickItem: string) => {
    if (!tickItem) return "";
    const parts = tickItem.split(" ");
    if (parts.length === 2) {
      return `${parts[0].substring(0, 3)} '${parts[1].substring(2)}`;
    }
    return tickItem;
  };

  return (
    // 3. Matched the container styling to your Ledger cards (rounded corners, softer border)
    <div className="border border-slate-800 p-6 bg-slate-900 rounded-xl shadow-sm">
      <h2 className="text-lg font-semibold text-white mb-6">{title}</h2>

      <ResponsiveContainer width="100%" height={350}>
        <ComposedChart
          data={data}
          margin={{ top: 10, right: 10, left: -20, bottom: 0 }}
        >
          <CartesianGrid
            strokeDasharray="3 3"
            stroke="#1e293b"
            vertical={false}
          />

          {/* 4. Removed the heavy axis line and applied the date formatter */}
          <XAxis
            dataKey="month"
            stroke="#64748b"
            tickFormatter={formatMonthAxis}
            tick={{ fontSize: 12 }}
            tickMargin={12}
            axisLine={false}
            tickLine={false}
          />

          {/* 5. Removed the heavy Y-axis line for a cleaner "floating" grid look */}
          <YAxis
            stroke="#64748b"
            tickFormatter={formatCompactCurrency}
            tick={{ fontSize: 12 }}
            axisLine={false}
            tickLine={false}
          />

          <Tooltip
            contentStyle={customTooltip}
            shared={true}
            formatter={(value: any) => formatCurrency(Number(value))}
          />

          {/* Added some padding so the legend doesn't crowd the X-axis labels */}
          <Legend iconType="circle" wrapperStyle={{ paddingTop: "16px" }} />

          {/* 6. Added maxBarSize for responsiveness and radius for rounded top corners */}
          <Bar
            dataKey="earned"
            name="Income"
            fill="#10b981"
            maxBarSize={40}
            radius={[4, 4, 0, 0]}
          />
          <Bar
            dataKey="donated"
            name="Donated"
            fill="#6366f1"
            maxBarSize={40}
            radius={[4, 4, 0, 0]}
          />

          {/* 7. Changed type to 'monotone' and added small dots to highlight data points */}
          <Line
            type="monotone"
            dataKey="target"
            name="Obligation"
            stroke="#ef4444"
            strokeWidth={3}
            dot={{ r: 4, fill: "#0f172a", strokeWidth: 2 }}
            activeDot={{ r: 6 }}
          />
        </ComposedChart>
      </ResponsiveContainer>
    </div>
  );
}
