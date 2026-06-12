import { Link } from "react-router-dom";

function Home() {
  return (
    <div className="min-h-screen bg-slate-950 text-slate-300 font-sans selection:bg-indigo-500/30">
      {/* Navigation Bar */}
      <nav className="container mx-auto px-6 py-6 flex justify-between items-center">
        <div className="text-2xl font-extrabold text-white tracking-tight flex items-center gap-2">
          {/* A simple placeholder logo icon */}
          <svg
            className="w-8 h-8 text-indigo-500"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth="2"
              d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
            ></path>
          </svg>
          Maaser<span className="text-indigo-400">Calc</span>
        </div>
        <div className="flex gap-4">
          <Link
            to="/login"
            className="px-5 py-2.5 text-sm font-semibold text-slate-300 hover:text-white transition-colors"
          >
            Sign In
          </Link>
          <Link
            to="/register"
            className="px-5 py-2.5 text-sm font-bold text-white bg-indigo-600 hover:bg-indigo-500 rounded-lg shadow-lg shadow-indigo-500/20 transition-all duration-200"
          >
            Get Started
          </Link>
        </div>
      </nav>

      {/* Hero Section */}
      <main className="container mx-auto px-6 pt-20 pb-32 flex flex-col items-center text-center">
        <div className="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-slate-900 border border-slate-800 text-sm text-indigo-300 mb-8">
          <span className="flex h-2 w-2 rounded-full bg-indigo-500"></span>
          Now available in Beta
        </div>

        <h1 className="text-5xl md:text-7xl font-extrabold text-white tracking-tight max-w-4xl leading-tight">
          Track your giving with <br className="hidden md:block" />
          <span className="text-transparent bg-clip-text bg-gradient-to-r from-indigo-400 to-cyan-400">
            absolute precision.
          </span>
        </h1>

        <p className="mt-8 text-lg md:text-xl text-slate-400 max-w-2xl">
          MaaserCalc simplifies your charitable contributions. Automatically
          calculate your tithes, track your donation history, and manage your
          financial commitments securely.
        </p>

        <div className="mt-10 flex flex-col sm:flex-row gap-4 w-full sm:w-auto">
          <Link
            to="/register"
            className="px-8 py-4 text-base font-bold text-white bg-indigo-600 hover:bg-indigo-500 rounded-xl shadow-lg shadow-indigo-500/20 transition-all duration-200 flex items-center justify-center gap-2"
          >
            Create Free Account
            <svg
              className="w-5 h-5"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth="2"
                d="M14 5l7 7m0 0l-7 7m7-7H3"
              ></path>
            </svg>
          </Link>
          <a
            href="#features"
            className="px-8 py-4 text-base font-bold text-slate-300 bg-slate-900 hover:bg-slate-800 border border-slate-800 rounded-xl transition-all duration-200 flex items-center justify-center"
          >
            Learn More
          </a>
        </div>
      </main>

      {/* Features Grid */}
      <section
        id="features"
        className="border-t border-slate-900 bg-slate-950/50 pt-24 pb-32"
      >
        <div className="container mx-auto px-6">
          <div className="text-center mb-16">
            <h2 className="text-3xl font-bold text-white">
              Everything you need to manage your Maaser
            </h2>
          </div>

          <div className="grid md:grid-cols-3 gap-8 max-w-5xl mx-auto">
            {/* Feature 1 */}
            <div className="bg-slate-900/50 border border-slate-800 p-8 rounded-2xl hover:bg-slate-900 transition-colors">
              <div className="w-12 h-12 bg-indigo-500/10 rounded-xl flex items-center justify-center mb-6 text-indigo-400">
                <svg
                  className="w-6 h-6"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth="2"
                    d="M9 7h6m0 10v-3m-3 3h.01M9 17h.01M9 14h.01M12 14h.01M15 11h.01M12 11h.01M9 11h.01M7 21h10a2 2 0 002-2V5a2 2 0 00-2-2H7a2 2 0 00-2 2v14a2 2 0 002 2z"
                  ></path>
                </svg>
              </div>
              <h3 className="text-xl font-bold text-white mb-3">
                Smart Calculations
              </h3>
              <p className="text-slate-400 leading-relaxed">
                Instantly calculate your 10% or 20% obligations based on your
                income, automatically adjusting for allowable deductions.
              </p>
            </div>

            {/* Feature 2 */}
            <div className="bg-slate-900/50 border border-slate-800 p-8 rounded-2xl hover:bg-slate-900 transition-colors">
              <div className="w-12 h-12 bg-cyan-500/10 rounded-xl flex items-center justify-center mb-6 text-cyan-400">
                <svg
                  className="w-6 h-6"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth="2"
                    d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"
                  ></path>
                </svg>
              </div>
              <h3 className="text-xl font-bold text-white mb-3">
                Donation Tracking
              </h3>
              <p className="text-slate-400 leading-relaxed">
                Keep a detailed ledger of all your charitable giving. Easily
                categorize donations and see your running balance at a glance.
              </p>
            </div>

            {/* Feature 3 */}
            <div className="bg-slate-900/50 border border-slate-800 p-8 rounded-2xl hover:bg-slate-900 transition-colors">
              <div className="w-12 h-12 bg-purple-500/10 rounded-xl flex items-center justify-center mb-6 text-purple-400">
                <svg
                  className="w-6 h-6"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth="2"
                    d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z"
                  ></path>
                </svg>
              </div>
              <h3 className="text-xl font-bold text-white mb-3">
                Secure & Private
              </h3>
              <p className="text-slate-400 leading-relaxed">
                Your financial data is yours. We use modern encryption standards
                to ensure your income and giving history remain totally private.
              </p>
            </div>
          </div>
        </div>
      </section>
    </div>
  );
}

export default Home;
