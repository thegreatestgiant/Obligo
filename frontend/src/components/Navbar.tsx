import { useNavigate, Link } from "react-router-dom";
import { useToast } from "./Toast";
import { useAuth } from "../auth/AuthContext";

function Navbar() {
  const navigate = useNavigate();
  const showToast = useToast();
  const { checkAuth } = useAuth();

  const logout = async (e: { preventDefault: () => void }) => {
    e.preventDefault();

    try {
      const response = await fetch("http://localhost:1234/logout", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        // body: JSON.stringify({}),
        credentials: "include",
      });

      if (response.ok) {
        showToast("Logged out");
        await checkAuth();
      } else {
        showToast("Couldn't log out 😥");
      }
    } catch (error) {
      console.error("Network Error:", error);
      showToast("Network issues 😥");
    }
  };

  return (
    <nav className="bg-slate-900 border-b border-slate-800 p-4 flex justify-between items-center">
      <div className="text-xl font-bold text-white">MaaserCalc</div>
      <div className="flex gap-4">
        <Link to="/dashboard" className="text-slate-300 hover:text-white">
          Dashboard
        </Link>
        <Link to="/settings" className="text-slate-300 hover:text-white">
          Settings
        </Link>
        <button onClick={logout} className="text-red-400 hover:text-red-300">
          Logout
        </button>
      </div>
    </nav>
  );
}

export default Navbar;
