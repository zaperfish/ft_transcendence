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
    <div className="min-h-screen flex flex-col">

      {/* Header */}
      <header className="w-full bg-surface border-b border-border px-xl py-md flex items-center justify-between">
        <div className="flex items-center gap-3">
          <span className="text-lg font-semibold text-text-primary tracking-tight">
            ft_transcendence
          </span>
        </div>
      </header>

      {/* Main */}
      <main className="flex flex-col items-center flex-1 px-lg py-3xl">
        <div className="w-full max-w-xl">

          <div className="mb-xl text-center">
            <h1 className="text-2xl font-bold text-text-primary mb-sm">API Test</h1>
            <p className="text-text-secondary text-sm">Fetch data from the backend and database</p>
          </div>

          {/* Button */}
          <div className="flex justify-center mb-xl">
            <button
              onClick={fetchData}
              disabled={loading}
              className="inline-flex items-center gap-2 px-lg py-sm bg-primary hover:bg-primary-dim active:scale-95 disabled:opacity-60 disabled:cursor-not-allowed text-surface font-semibold rounded-md transition-all duration-200"
            >
              {loading ? (
                <>
                  <span className="w-3.5 h-3.5 border-2 border-surface/30 border-t-surface rounded-full animate-spin inline-block" />
                  Fetching…
                </>
              ) : (
                "Fetch Data"
              )}
            </button>
          </div>

          {/* Error */}
          {error && (
            <div className="bg-surface border border-error rounded-md p-md mb-md">
              <p className="text-xs font-bold tracking-widest uppercase text-error mb-sm">Error</p>
              <p className="text-error text-sm">{error}</p>
            </div>
          )}

          {/* Result */}
          {result !== null && (
            <div className="bg-surface border border-border rounded-md p-md">
              <p className="text-xs font-bold tracking-widest uppercase text-primary mb-sm">Response</p>
              <p className="text-xs text-text-tertiary mb-md">
                Served by the backend · queried from the database
              </p>
              <pre className="text-sm leading-relaxed text-text-secondary whitespace-pre-wrap break-words overflow-x-auto bg-surface-container rounded-md p-md">
                {JSON.stringify(result, null, 2)}
              </pre>
            </div>
          )}

          {/* Hint */}
          {!loading && result === null && !error && (
            <p className="text-center text-text-tertiary text-sm mt-md">
              Press the button to fetch from{" "}
              <code className="bg-surface-container px-xs py-xs rounded text-primary text-xs font-mono">
                /api/test
              </code>
            </p>
          )}
        </div>
      </main>
    </div>
  );
}
