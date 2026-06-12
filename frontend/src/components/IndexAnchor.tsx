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
  const customTooltip = {
    backgroundColor: "#000",
    borderColor: "#ef4444",
    borderWidth: "2px",
    color: "#fff",
    borderRadius: "0px",
  };

  return (
    <div className="border-2 border-slate-800 p-6 bg-slate-900">
      <h2 className="text-lg font-black uppercase mb-6 text-white">{title}</h2>
      <ResponsiveContainer width="100%" height={350}>
        <ComposedChart data={data}>
          <CartesianGrid
            strokeDasharray="3 3"
            stroke="#1e293b"
            vertical={false}
          />
          <XAxis dataKey="month" stroke="#64748b" />
          <YAxis stroke="#64748b" tickFormatter={formatCompactCurrency} />

          {/* STICKY UX: triggers on the axis slice, not just the data point */}
          <Tooltip
            contentStyle={customTooltip}
            shared={true}
            formatter={(value: any) => formatCurrency(Number(value))}
          />

          <Legend iconType="circle" />
          <Bar dataKey="earned" name="Income" fill="#10b981" barSize={30} />
          <Bar dataKey="donated" name="Donated" fill="#6366f1" barSize={30} />
          <Line
            type="step"
            dataKey="target"
            name="Obligation"
            stroke="#ef4444"
            strokeWidth={3}
            dot={false}
          />
        </ComposedChart>
      </ResponsiveContainer>
    </div>
  );
}
