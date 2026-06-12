import React, { createContext, useContext, useState, useEffect } from "react";
import { api } from "../api/clients";

type UserSummary = {
  Total_Owed: number;
  Donation_Percent: number;
  Percent_Fulfilled: number;
  Remaining_Obligation: number;
  Total_Donated: number;
  Total_Earned: number;
};

type AuthContextType = {
  user: UserSummary | null;
  isLoading: boolean;
  checkAuth: (silent?: boolean) => Promise<void>;
};

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<UserSummary | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  const checkAuth = async (silent: boolean = false) => {
    if (!silent) setIsLoading(true);
    try {
      const response = await api.getSummary();

      if (response.ok) {
        const data = await response.json();
        setUser(data);
      } else {
        setUser(null);
      }
    } catch (error) {
      console.error("Auth check failed:", error);
      setUser(null);
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    checkAuth();
  }, []);

  useEffect(() => {
    const handleAuthExpired = () => {
      console.warn("Session expired. Logging out globally.");
      setUser(null); // This triggers ProtectedRoute to redirect to /login
    };

    window.addEventListener("auth-expired", handleAuthExpired);

    // Cleanup the listener when the provider unmounts
    return () => window.removeEventListener("auth-expired", handleAuthExpired);
  }, []);
  return (
    <AuthContext.Provider value={{ user, isLoading, checkAuth }}>
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error("useAuth must be used within an AuthProvider");
  }
  return context;
}
