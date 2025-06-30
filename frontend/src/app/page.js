"use client";

import { useEffect, useState } from "react";

export default function Home() {
  const [data, setData] = useState(null);
  const [error, setError] = useState(null);

  useEffect(() => {
    const apiUrl = process.env.NEXT_PUBLIC_API_URL;

    fetch(`${apiUrl}/`)
      .then((res) => {
        if (!res.ok) throw new Error(`API returned ${res.status}`);
        return res.json();
      })
      .then(setData)
      .catch((err) => {
        console.error(err);
        setError("Failed to load API status");
      });
  }, []);

  return (
    <div className="grid grid-rows-[20px_1fr_20px] items-center justify-items-center min-h-screen p-8 pb-20 gap-16 sm:p-20 font-[family-name:var(--font-geist-sans)]">
      <main className="flex flex-col gap-[32px] row-start-2 items-center sm:items-start">
        <h1 className="text-2xl font-bold">API Status</h1>

        {error && <p className="text-red-600">{error}</p>}
        {data ? (
          <div>
            <p><strong>Status:</strong> {data.status}</p>
            <p><strong>Message:</strong> {data.message}</p>
          </div>
        ) : !error ? (
          <p>Loading...</p>
        ) : null}
      </main>

      <footer className="row-start-3 flex gap-[24px] flex-wrap items-center justify-center">
        &copy; 2025 tajjjjr
      </footer>
    </div>
  );
}
