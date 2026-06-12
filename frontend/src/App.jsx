import { Routes, Route } from "react-router-dom";
import Login from "./pages/Login";
import Register from "./pages/Register";
import Home from "./pages/Home";
import Dashboard from "./pages/Dashboard";
import { AuthProvider } from "./AuthContext";
import ProtectedRoute from "./ProtectedRoute";
import Settings from "./pages/Settings";

function App() {
  return (
    <AuthProvider>
      <Routes>
        {/* Public Routes */}
        <Route path="/" element={<Home />} />
        <Route path="/login" element={<Login />} />
        <Route path="/register" element={<Register />} />

        {/* Protected Routes (The Bouncer) */}
        <Route element={<ProtectedRoute />}>
          {/* Anything inside here REQUIRES a valid cookie */}
          <Route path="/dashboard" element={<Dashboard />} />
          <Route path="/settings" element={<Settings />} />
        </Route>
      </Routes>
    </AuthProvider>
  );
}

export default App;
