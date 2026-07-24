import React, { useState } from "react";
import { useNavigate, Link } from "react-router-dom";
import ".././index.css";
import { api } from "../api/clients";

type ErrorType = {
  email?: string;
  username?: string;
  password?: string;
  submit?: string;
};

const INITIAL_FORM_DATA = {
  email: "",
  username: "",
  password: "",
};

function Register() {
  const [formData, setFormData] = useState(INITIAL_FORM_DATA);
  const [errors, setErrors] = useState<ErrorType>({});
  const [isSubmitting, setIsSubmitting] = useState(false);
  const navigate = useNavigate();

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const name = e.target.name as keyof typeof INITIAL_FORM_DATA;
    const value = e.target.value;

    setFormData((prev) => ({ ...prev, [name]: value }));

    // Clear the specific error when the user starts typing again
    if (errors[name]) {
      setErrors((prev) => ({ ...prev, [name]: "" }));
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setErrors({});
    setIsSubmitting(true);

    try {
      await api.register(
        formData.email.trim(),
        formData.username.trim(),
        formData.password,
      );

      navigate("/login", { replace: true });
    } catch (error: any) {
      console.error("Registration Error:", error);
      if (error.message.includes("409")) {
        setErrors({ submit: "Username or email already exists." });
      } else if (error.message === "Unauthorized") {
        setErrors({ submit: "Not authorized to register." });
      } else {
        setErrors({ submit: "Registration failed. Please try again." });
      }
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="min-h-screen relative flex items-center justify-center bg-slate-950 py-12 px-4 sm:px-6 lg:px-8">
      <Link
        to="/"
        className="absolute top-8 left-8 flex items-center gap-2 text-sm font-medium text-slate-400 hover:text-white transition-colors group"
      >
        <svg
          className="w-5 h-5 transform group-hover:-translate-x-1 transition-transform"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            strokeWidth="2"
            d="M10 19l-7-7m0 0l7-7m-7 7h18"
          ></path>
        </svg>
        Back to Home
      </Link>
      {/* The Register Card: Slightly lighter slate, dark borders, soft shadow */}
      <div className="max-w-md w-full space-y-8 bg-slate-900 p-10 rounded-2xl shadow-2xl border border-slate-800">
        {/* Header Section */}
        <div>
          <h2 className="mt-2 text-center text-3xl font-extrabold text-white tracking-tight">
            Create an Account
          </h2>
          <p className="mt-3 text-center text-sm text-slate-400">
            Enter your details to register for Obligo.
          </p>
        </div>

        {/* The Form */}
        <form className="mt-8 space-y-6" onSubmit={handleSubmit}>
          <div className="space-y-5">
            {/* Email Input */}
            <div>
              <label
                htmlFor="email"
                className="block text-sm font-semibold text-slate-300 mb-1"
              >
                Email
              </label>
              <input
                id="email"
                name="email"
                type="email" // Changed to type="email" for built-in browser validation!
                required
                value={formData.email}
                onChange={handleChange}
                placeholder="Enter your email address"
                className={`w-full px-4 py-2.5 rounded-lg border bg-slate-950 text-white placeholder-slate-500 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:bg-slate-900 transition-all duration-200 ${
                  errors.email
                    ? "border-red-500 focus:ring-red-500"
                    : "border-slate-700"
                }`}
              />
              {errors.email && (
                <p className="mt-1.5 text-sm text-red-500 font-medium">
                  {errors.email}
                </p>
              )}
            </div>

            {/* Username Input */}
            <div>
              <label
                htmlFor="username"
                className="block text-sm font-semibold text-slate-300 mb-1"
              >
                Username
              </label>
              <input
                id="username"
                name="username"
                type="text"
                required
                value={formData.username}
                onChange={handleChange}
                placeholder="Choose a username"
                className={`w-full px-4 py-2.5 rounded-lg border bg-slate-950 text-white placeholder-slate-500 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:bg-slate-900 transition-all duration-200 ${
                  errors.username
                    ? "border-red-500 focus:ring-red-500"
                    : "border-slate-700"
                }`}
              />
              {errors.username && (
                <p className="mt-1.5 text-sm text-red-500 font-medium">
                  {errors.username}
                </p>
              )}
            </div>

            {/* Password Input */}
            <div>
              <label
                htmlFor="password"
                className="block text-sm font-semibold text-slate-300 mb-1"
              >
                Password
              </label>
              <input
                id="password"
                name="password"
                type="password"
                required
                value={formData.password}
                onChange={handleChange}
                placeholder="••••••••"
                className={`w-full px-4 py-2.5 rounded-lg border bg-slate-950 text-white placeholder-slate-500 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:bg-slate-900 transition-all duration-200 ${
                  errors.password
                    ? "border-red-500 focus:ring-red-500"
                    : "border-slate-700"
                }`}
              />
              {errors.password && (
                <p className="mt-1.5 text-sm text-red-500 font-medium">
                  {errors.password}
                </p>
              )}
            </div>
          </div>

          {/* Top-Level Submit Error */}
          {errors.submit && (
            <div className="p-3 rounded-lg bg-red-500/10 border border-red-500/20">
              <p className="text-sm text-center text-red-400 font-medium">
                {errors.submit}
              </p>
            </div>
          )}

          {/* Submit Button */}
          <div>
            <button
              type="submit"
              disabled={isSubmitting}
              className="w-full flex justify-center py-2.5 px-4 border border-transparent rounded-lg shadow-sm text-sm font-bold text-white bg-indigo-600 hover:bg-indigo-500 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-offset-slate-900 focus:ring-indigo-500 transition-all duration-200 disabled:opacity-70 disabled:cursor-not-allowed"
            >
              {isSubmitting ? (
                <span className="flex items-center gap-2">
                  <svg
                    className="animate-spin h-4 w-4 text-white"
                    xmlns="http://www.w3.org/2000/svg"
                    fill="none"
                    viewBox="0 0 24 24"
                  >
                    <circle
                      className="opacity-25"
                      cx="12"
                      cy="12"
                      r="10"
                      stroke="currentColor"
                      strokeWidth="4"
                    ></circle>
                    <path
                      className="opacity-75"
                      fill="currentColor"
                      d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                    ></path>
                  </svg>
                  Registering...
                </span>
              ) : (
                "Create Account"
              )}
            </button>
          </div>

          {/* Footer Link */}
          <div className="flex items-center justify-center pt-2">
            <span className="text-sm text-slate-400">
              Already have an account?{" "}
              <Link
                to="/login"
                className="font-semibold text-indigo-400 hover:text-indigo-300 transition-colors"
              >
                Sign in here
              </Link>
            </span>
          </div>
        </form>
      </div>
    </div>
  );
}

export default Register;
