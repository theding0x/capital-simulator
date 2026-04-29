import { useEffect, useState } from "react";

interface GatewayInfo {
  service: string;
  status: string;
  description: string;
  downstream: string[];
}

export default function App() {
  const [info, setInfo] = useState<GatewayInfo | null>(null);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    fetch("/api/v1/info")
      .then((r) => {
        if (!r.ok) throw new Error(`gateway returned ${r.status}`);
        return r.json() as Promise<GatewayInfo>;
      })
      .then(setInfo)
      .catch((e: unknown) => setError(e instanceof Error ? e.message : String(e)));
  }, []);

  return (
    <main className="app">
      <header className="hero">
        <h1>Capital Simulator</h1>
        <p className="subtitle">
          A simulation of an economy modeled chapter by chapter on Karl Marx&rsquo;s
          <em> Capital, Volume I</em>.
        </p>
      </header>

      <section className="card">
        <h2>Gateway status</h2>
        {error && <p className="error">Could not reach api-gateway: {error}</p>}
        {!error && !info && <p>Loading&hellip;</p>}
        {info && (
          <>
            <p>
              <strong>{info.service}</strong> &mdash; {info.status}
            </p>
            <p>{info.description}</p>
            <ul>
              {info.downstream.map((s) => (
                <li key={s}>{s}</li>
              ))}
            </ul>
          </>
        )}
      </section>

      <footer className="footer">
        <small>Scaffold &middot; chapter 0</small>
      </footer>
    </main>
  );
}
