import React, { createContext, useContext, useState, useEffect } from "react";

// 1. Update the shape to match your Go 'summary' struct exactly!
// (Assuming your Go struct uses the exact field names for JSON)
type UserSummary = {
  Total_Owed: number; // Change to match Go
  Donation_Percent: number; // Add underscore
  Percent_Fulfilled: number; // Add underscore
  Remaining_Obligation: number; // Add underscore
  Total_Donated: number; // Add underscore
  Total_Earned: number; // Add underscore
};

type AuthContextType = {
  user: UserSummary | null;
  isLoading: boolean;
  checkAuth: () => Promise<void>;
};

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<UserSummary | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  const checkAuth = async () => {
    setIsLoading(true);
    try {
      // 2. Point this to your summary endpoint!
      const response = await fetch("http://localhost:1234/summary", {
        method: "GET",
        credentials: "include",
      });

      if (response.ok) {
        const data = await response.json();
        setUser(data); // Save the summary data to the global state!
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
  const user = context.user;
  console.log(user);
  console.log(user?.Donation_Percent);
  return context;
}
