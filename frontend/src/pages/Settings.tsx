import React, { useState } from "react";
import Navbar from "./Navbar";

function Settings() {
  const [percent, setPercent] = useState("");
  const [passwords, setPasswords] = useState({ old: "", new: "" });

  const handleUpdatePercent = async () => {
    await fetch("http://localhost:1234/users/settings", {
      method: "PATCH",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ donation_percentage: parseFloat(percent) }),
      credentials: "include",
    });
    alert("Updated!");
  };

  const handleChangePassword = async () => {
    await fetch("http://localhost:1234/users/change-password", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        old_password: passwords.old,
        new_password: passwords.new,
      }),
      credentials: "include",
    });
    alert("Password updated!");
  };

  return (
    <div className="min-h-screen bg-slate-950 text-white">
      <Navbar />
      <div className="max-w-2xl mx-auto p-10 space-y-12">
        {/* Donation Settings */}
        <section className="bg-slate-900 p-6 rounded-xl border border-slate-800">
          <h2 className="text-xl font-bold mb-4">Donation Percentage</h2>
          <input
            type="number"
            value={percent}
            onChange={(e) => setPercent(e.target.value)}
            className="bg-slate-950 border border-slate-700 p-2 rounded w-full mb-4"
            placeholder="10"
          />
          <button
            onClick={handleUpdatePercent}
            className="bg-indigo-600 px-4 py-2 rounded"
          >
            Update
          </button>
        </section>

        {/* Security Settings */}
        <section className="bg-slate-900 p-6 rounded-xl border border-slate-800">
          <h2 className="text-xl font-bold mb-4">Change Password</h2>
          <input
            type="password"
            placeholder="Old Password"
            onChange={(e) =>
              setPasswords({ ...passwords, old: e.target.value })
            }
            className="bg-slate-950 border border-slate-700 p-2 rounded w-full mb-2"
          />
          <input
            type="password"
            placeholder="New Password"
            onChange={(e) =>
              setPasswords({ ...passwords, new: e.target.value })
            }
            className="bg-slate-950 border border-slate-700 p-2 rounded w-full mb-4"
          />
          <button
            onClick={handleChangePassword}
            className="bg-indigo-600 px-4 py-2 rounded"
          >
            Update Password
          </button>
        </section>
      </div>
    </div>
  );
}

export default Settings;
