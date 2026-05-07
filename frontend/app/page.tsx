"use client";

import { useState } from "react";

export default function Page() {
  const [result, setResult] = useState<unknown>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  async function fetchData() {
    setLoading(true);
    setError(null);
    setResult(null);
    try {
      const res = await fetch("/api/postgres-version");
      if (!res.ok) throw new Error(`HTTP ${res.status}: ${res.statusText}`);
      const data = await res.json();
      setResult(data);
    } catch (err) {
      setError((err as Error).message);
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className="min-h-screen flex flex-col bg-slate-50 text-slate-800 font-sans">

      {/* Header */}
      <header className="w-full bg-white border-b border-slate-200 shadow-sm px-8 py-4 flex items-center justify-between">
        <div className="flex items-center gap-3">
          <span className="text-lg font-semibold text-slate-800 tracking-tight">
            ft_transcendence
          </span>
        </div>
      </header>

      {/* Main */}
      <main className="flex flex-col items-center flex-1 px-6 py-16">
        <div className="w-full max-w-xl">

          <div className="mb-10 text-center">
            <h1 className="text-2xl font-bold text-slate-800 mb-2">API Test</h1>
            <p className="text-slate-500 text-sm">Fetch data from the backend and database</p>
          </div>

          {/* Button */}
          <div className="flex justify-center mb-8">
            <button
              onClick={fetchData}
              disabled={loading}
              className="inline-flex items-center gap-2 px-6 py-2.5 bg-indigo-600 hover:bg-indigo-700 active:scale-95 disabled:opacity-60 disabled:cursor-not-allowed text-white font-semibold rounded-lg shadow-sm transition-all duration-200"
            >
              {loading ? (
                <>
                  <span className="w-3.5 h-3.5 border-2 border-white/30 border-t-white rounded-full animate-spin inline-block" />
                  Fetching…
                </>
              ) : (
                "Fetch Data"
              )}
            </button>
          </div>

          {/* Error */}
          {error && (
            <div className="bg-red-50 border border-red-200 rounded-xl p-5">
              <p className="text-xs font-bold tracking-widest uppercase text-red-500 mb-2">Error</p>
              <p className="text-red-600 text-sm">{error}</p>
            </div>
          )}

          {/* Result */}
          {result !== null && (
            <div className="bg-white border border-slate-200 rounded-xl p-5 shadow-sm">
              <p className="text-xs font-bold tracking-widest uppercase text-indigo-500 mb-1">Response</p>
              <p className="text-xs text-slate-400 mb-4">
                Served by the backend · queried from the database
              </p>
              <pre className="text-sm leading-relaxed text-slate-700 whitespace-pre-wrap break-words overflow-x-auto bg-slate-50 rounded-lg p-4">
                {JSON.stringify(result, null, 2)}
              </pre>
            </div>
          )}

          {/* Hint */}
          {!loading && result === null && !error && (
            <p className="text-center text-slate-400 text-sm">
              Press the button to fetch from{" "}
              <code className="bg-slate-100 px-1.5 py-0.5 rounded text-indigo-500 text-xs font-mono">
                /api/test
              </code>
            </p>
          )}
        </div>
      </main>
    </div>
  );
}
