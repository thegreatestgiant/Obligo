import { createContext, useContext, useState, ReactNode } from "react";

type ToastType = { message: string; type: "success" | "error" };
type ToastContextType = (message: string, type?: "success" | "error") => void;

const ToastContext = createContext<ToastContextType | undefined>(undefined);

export function ToastProvider({ children }: { children: ReactNode }) {
  const [toast, setToast] = useState<ToastType | null>(null);

  const showToast = (
    message: string,
    type: "success" | "error" = "success",
  ) => {
    setToast({ message, type });
    setTimeout(() => setToast(null), 3000);
  };

  return (
    <ToastContext.Provider value={showToast}>
      {children}
      {toast && (
        <div
          onClick={() => setToast(null)}
          className={`fixed bottom-5 left-5 px-6 py-3 rounded-lg shadow-2xl border transition-all duration-300 z-50 cursor-pointer hover:opacity-80 ${
            toast.type === "success"
              ? "bg-emerald-900 border-emerald-700 text-emerald-100"
              : "bg-red-900 border-red-700 text-red-100"
          }`}
        >
          {toast.message}
        </div>
      )}
    </ToastContext.Provider>
  );
}

export const useToast = () => {
  const context = useContext(ToastContext);
  if (!context) throw new Error("useToast must be used within a ToastProvider");
  return context;
};
