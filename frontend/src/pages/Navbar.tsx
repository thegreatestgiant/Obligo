import { Link } from "react-router-dom";

function Navbar() {
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
        <button className="text-red-400 hover:text-red-300">Logout</button>
      </div>
    </nav>
  );
}

export default Navbar;
