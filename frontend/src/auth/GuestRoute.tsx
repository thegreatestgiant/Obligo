import { Navigate, Outlet } from "react-router-dom";

import { useAuth } from "./AuthContext";

export default function GuestRoute() {
  const { user, isLoading } = useAuth();

  // If the app is still talking to the Go server, show a loading screen
  if (isLoading) {
    return (
      <div className="min-h-screen bg-slate-950 flex items-center justify-center">
        <div className="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-indigo-500"></div>
      </div>
    );
  }

  if (user) {
    return <Navigate to="/dashboard" replace />;
  }

  return <Outlet />;
}
