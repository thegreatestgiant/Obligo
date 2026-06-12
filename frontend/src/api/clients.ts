const BASE_URL = "http://localhost:1234";

// Centralized fetch helper to automatically handle credentials and JSON
async function fetchAPI(endpoint: string, options: RequestInit = {}) {
  const response = await fetch(`${BASE_URL}${endpoint}`, {
    ...options,
    headers: {
      "Content-Type": "application/json",
      ...options.headers,
    },
    credentials: "include", // Ensures our Go cookies are always sent!
  });

  if (response.status === 401) {
    // Shout to the React app that the session is dead
    window.dispatchEvent(new Event("auth-expired"));

    // Throw a specific error so your components can choose to ignore the generic toast
    throw new Error("Unauthorized");
  }

  if (!response.ok) {
    throw new Error(`API Error: ${response.status}`);
  }

  return response;
}

export const api = {
  // Auth & User
  logout: () => fetchAPI("/logout", { method: "POST" }),
  getSummary: () => fetchAPI("/summary", { method: "GET" }),

  // Ledger Entries
  getEntries: () => fetchAPI("/entries", { method: "GET" }),
  createEntry: (
    amount: number,
    type: "paycheck" | "donation",
    description: string,
  ) =>
    fetchAPI("/entries", {
      method: "POST",
      body: JSON.stringify({ amount, ledger_entry: type, description }),
    }),
  getEntry: (id: string) => fetchAPI(`/entries/${id}`, { method: "GET" }),
  deleteEntry: (id: string) => fetchAPI(`/entries/${id}`, { method: "DELETE" }),
};
