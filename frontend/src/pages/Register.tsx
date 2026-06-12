import { useState } from "react";
import { useNavigate, Link } from "react-router-dom";

function Register() {
  type ErrorType = {
    email?: string;
    username?: string;
    password?: string;
    submit?: string;
  };

  const INITIAL_ERRORS: ErrorType = {};

  const [formData, setFormData] = useState(INITIAL_ERRORS);
  const [errors, setErrors] = useState<ErrorType>({});
  const [isSubmitting, setIsSubmitting] = useState(false);
  const navigate = useNavigate();

  const handleChange = (e: { target: { name: any; value: any } }) => {
    const name = e.target.name as keyof ErrorType;
    const value = e.target.value;
    setFormData((prev) => ({ ...prev, [name]: value }));
    // Clear the specific error when the user starts typing again
    if (errors[name]) {
      setErrors((prev) => ({ ...prev, [name]: "" }));
    }
  };

  const handleSubmit = async (formData: FormData) => {
    const email = formData.get("email") as string;
    const username = formData.get("username") as string;
    const password = formData.get("password") as string;
    setErrors(INITIAL_ERRORS);
    setIsSubmitting(true);

    try {
      const response = await fetch("http://localhost:1234/register", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ email, username, password }),
        credentials: "include",
      });

      if (response.ok) {
        navigate("/app", { replace: true });
      } else {
        setErrors((prev) => ({
          ...prev,
          submit: "Invalid username or password.",
        }));
      }
    } catch (error) {
      console.error("Network Error:", error);
      setErrors((prev) => ({
        ...prev,
        submit: "Cannot connect to the server.",
      }));
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    // Outer container: takes up the full screen height, centers the card, adds a soft background color
    <div className="min-h-screen flex items-center justify-center bg-slate-50 py-12 px-4 sm:px-6 lg:px-8">
      {/* The Login Card */}
      <div className="max-w-md w-full space-y-8 bg-white p-10 rounded-2xl shadow-xl border border-slate-100">
        {/* Header Section */}
        <div>
          <h2 className="mt-2 text-center text-3xl font-extrabold text-slate-900 tracking-tight">
            Welcome
          </h2>
          <p className="mt-3 text-center text-sm text-slate-500">
            Please enter your details to Register
          </p>
        </div>

        {/* The Form */}
        <form className="mt-8 space-y-6" action={handleSubmit}>
          <div className="space-y-5">
            <div>
              <label
                htmlFor="email"
                className="block text-sm font-semibold text-slate-700 mb-1"
              >
                Email
              </label>
              <input
                id="email"
                name="email"
                type="text"
                required
                value={formData.email}
                onChange={handleChange}
                placeholder="Enter your email address"
                className={`w-full px-4 py-2.5 rounded-lg border bg-slate-50 text-slate-900 placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:bg-white transition-all duration-200 ${
                  errors.email
                    ? "border-red-500 focus:ring-red-500"
                    : "border-slate-200"
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
                className="block text-sm font-semibold text-slate-700 mb-1"
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
                placeholder="Enter your username"
                className={`w-full px-4 py-2.5 rounded-lg border bg-slate-50 text-slate-900 placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:bg-white transition-all duration-200 ${
                  errors.username
                    ? "border-red-500 focus:ring-red-500"
                    : "border-slate-200"
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
                className="block text-sm font-semibold text-slate-700 mb-1"
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
                className={`w-full px-4 py-2.5 rounded-lg border bg-slate-50 text-slate-900 placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:bg-white transition-all duration-200 ${
                  errors.password
                    ? "border-red-500 focus:ring-red-500"
                    : "border-slate-200"
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
            <div className="p-3 rounded-lg bg-red-50 border border-red-200">
              <p className="text-sm text-center text-red-600 font-medium">
                {errors.submit}
              </p>
            </div>
          )}

          {/* Submit Button */}
          <div>
            <button
              type="submit"
              disabled={isSubmitting}
              className="w-full flex justify-center py-2.5 px-4 border border-transparent rounded-lg shadow-sm text-sm font-bold text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 transition-all duration-200 disabled:opacity-70 disabled:cursor-not-allowed"
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
                "Register"
              )}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
export default Register;
