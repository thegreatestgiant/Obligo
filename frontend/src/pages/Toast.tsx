import { useState, useEffect } from "react";

export function useToast() {
  const [toast, setToast] = useState<{
    message: string;
    type: "success" | "error";
  } | null>(null);

  const showToast = (
    message: string,
    type: "success" | "error" = "success",
  ) => {
    setToast({ message, type });
    setTimeout(() => setToast(null), 3000); // Fades out after 3s
  };

  const ToastComponent = toast ? (
    <div
      className={`fixed top-5 right-5 px-6 py-3 rounded-lg shadow-2xl border transition-all duration-300 ease-in-out opacity-100 ${
        toast.type === "success"
          ? "bg-emerald-900 border-emerald-700 text-emerald-100"
          : "bg-red-900 border-red-700 text-red-100"
      }`}
    >
      {toast.message}
    </div>
  ) : null;

  return { showToast, ToastComponent };
}
