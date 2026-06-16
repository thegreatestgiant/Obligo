async function fetchAPI(endpoint: string, options: RequestInit = {}) {
  const baseUrl = (window as any).API_URL || "";

  const response = await fetch(`${baseUrl}${endpoint}`, {
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
  register: (email: string, username: string, password: string) =>
    fetchAPI("/register", {
      method: "POST",
      body: JSON.stringify({ email, username, password }),
    }),
  login: (username: string, password: string) =>
    fetchAPI("/login", {
      method: "POST",
      body: JSON.stringify({
        username: username,
        password: password,
      }),
    }),

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

  // Settings
  updatePercent: (percent: number) =>
    fetchAPI("/users/settings", {
      method: "PATCH",
      body: JSON.stringify({ donation_percentage: percent }),
    }),
  changePassword: (old_password: string, new_password: string) =>
    fetchAPI("/users/settings", {
      method: "PATCH",
      body: JSON.stringify({
        old_password: old_password,
        new_password: new_password,
      }),
    }),
};
