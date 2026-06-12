import { useState } from "react";
import { PieChart, Pie, Tooltip, ResponsiveContainer } from "recharts";
import { formatCurrency } from "../utils/formatter";

interface GoalWidgetProps {
  totalTarget: number;
  totalDonated: number;
}

export default function GoalWidget({
  totalTarget,
  totalDonated,
}: GoalWidgetProps) {
  const [isHalfDonut, setIsHalfDonut] = useState(true);

  const remaining = Math.max(0, totalTarget - totalDonated);
  const percentFulfilled =
    totalTarget > 0 ? ((totalDonated / totalTarget) * 100).toFixed(1) : "0.0";

  const pieData = [
    { name: "Donated", value: totalDonated, fill: "#6366f1" },
    { name: "Still Owed", value: remaining, fill: "#ef4444" },
  ];

  const customTooltip = {
    backgroundColor: "#0f172a",
    borderColor: "#1e293b",
    color: "#fff",
    borderRadius: "8px",
  };

  return (
    <div className="bg-slate-900 border border-slate-800 p-8 rounded-2xl shadow-sm h-[450px] flex flex-col relative overflow-hidden">
      <div className="flex justify-between items-center mb-4 z-10">
        <h2 className="text-xl font-semibold text-slate-200 tracking-tight">
          Goal Progress
        </h2>
        <button
          onClick={() => setIsHalfDonut(!isHalfDonut)}
          className="text-xs font-bold uppercase tracking-widest bg-slate-800 hover:bg-slate-700 text-slate-300 px-4 py-2 rounded-lg transition-all active:scale-95"
        >
          {isHalfDonut ? "Full Circle" : "Speedometer"}
        </button>
      </div>

      <div className="flex-grow">
        <ResponsiveContainer width="100%" height="100%">
          <PieChart>
            <Pie
              data={pieData}
              cx="50%"
              cy={isHalfDonut ? "85%" : "50%"}
              startAngle={isHalfDonut ? 180 : 360}
              endAngle={0}
              innerRadius={isHalfDonut ? 120 : 90}
              outerRadius={isHalfDonut ? 165 : 130}
              paddingAngle={4}
              dataKey="value"
              stroke="none"
              animationBegin={0}
              animationDuration={1200}
              animationEasing="ease-in-out"
              fill="fill"
            />
            <Tooltip
              contentStyle={customTooltip}
              formatter={(value: any) => formatCurrency(Number(value))}
            />
          </PieChart>
        </ResponsiveContainer>

        <div
          className="absolute left-0 right-0 text-center pointer-events-none transition-all duration-1000 ease-in-out"
          style={{
            transform: isHalfDonut
              ? "translateY(-100px)"
              : "translateY(-210px)",
            top: "100%",
          }}
        >
          <p className="text-5xl font-black text-white tracking-tighter">
            {percentFulfilled}%
          </p>
          <p className="text-xs font-bold uppercase tracking-[0.2em] text-slate-500 mt-1">
            Fulfilled
          </p>
        </div>
      </div>
    </div>
  );
}
