interface KPICardProps {
  title: string;
  value: string;
  subtext?: string;
  percent?: string; // e.g., "71.2%"
  intent: "earned" | "donated" | "obligation";
}

export default function KPICard({
  title,
  value,
  subtext,
  percent,
  intent,
}: KPICardProps) {
  // Mapping intents to styles
  const styles = {
    earned: "border-emerald-500 text-emerald-400",
    donated: "border-indigo-500 text-indigo-300",
    obligation: "border-red-500 text-red-500/70",
  };

  return (
    <div
      className={`bg-slate-900 border-4 p-8 rounded-3xl relative overflow-hidden ${styles[intent]}`}
    >
      {/* 1. Donated Wavy Fill */}
      {intent === "donated" && percent && (
        <div
          className="absolute bottom-0 left-0 w-full bg-indigo-500/20 transition-all duration-1000"
          style={{ height: `${percent}` }}
        >
          <div className="absolute top-[-10px] w-full h-4 bg-indigo-500/30 rounded-[50%] animate-pulse"></div>
        </div>
      )}

      {/* 2. Obligation Laser Base (Perfectly flush) */}
      {intent === "obligation" && percent && (
        <div
          className="absolute bottom-0 left-0 h-2 bg-red-500 shadow-[0_0_20px_rgba(239,68,68,0.8)] z-20 rounded-br-2xl"
          style={{ width: `${percent}` }}
        ></div>
      )}

      <div className="relative z-10">
        <p
          className={
            "text-xs font-black uppercase tracking-[0.2em] mb-2 text-slate-400"
          }
        >
          {title}
        </p>
        <p className="text-4xl font-black">{value}</p>
        {subtext && (
          <p className="text-xs text-slate-500 mt-2 font-medium">{subtext}</p>
        )}
      </div>
    </div>
  );
}
