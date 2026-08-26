import { useEffect, useState } from "react";
import { createRoot } from "react-dom/client";
import { Dashboard } from "./Dashboard";
import { LoginScreen } from "./LoginScreen";
import { ProductsScreen } from "./ProductsScreen";
import { APIError, stockpilotAPI } from "./api";
import type { User } from "./types";
import "./styles.css";

type SessionState =
  | { kind: "checking" }
  | { kind: "signed-out" }
  | { kind: "signed-in"; user: User };

type View = "dashboard" | "products";

function App() {
  const [session, setSession] = useState<SessionState>({ kind: "checking" });
  const [view, setView] = useState<View>("dashboard");

  useEffect(() => {
    let active = true;

    void stockpilotAPI.me().then(
      (user) => {
        if (active) setSession({ kind: "signed-in", user });
      },
      (error: unknown) => {
        if (!active) return;
        if (error instanceof APIError && error.status === 401) {
          setSession({ kind: "signed-out" });
          return;
        }
        setSession({ kind: "signed-out" });
      }
    );

    return () => {
      active = false;
    };
  }, []);

  async function login(email: string, password: string) {
    const result = await stockpilotAPI.login(email, password);
    setView("dashboard");
    setSession({ kind: "signed-in", user: result.user });
  }

  async function logout() {
    try {
      await stockpilotAPI.logout();
    } finally {
      setView("dashboard");
      setSession({ kind: "signed-out" });
    }
  }

  function expireSession() {
    setView("dashboard");
    setSession({ kind: "signed-out" });
  }

  if (session.kind === "checking") {
    return (
      <main className="auth-page" aria-busy="true" aria-live="polite">
        <section className="auth-card session-check" aria-labelledby="session-check-title">
          <span className="brand-mark" aria-hidden="true">SP</span>
          <p className="eyebrow">Secure workspace</p>
          <h1 id="session-check-title">Checking your session</h1>
          <p className="muted">StockPilot is verifying your signed-in session.</p>
        </section>
      </main>
    );
  }

  if (session.kind === "signed-out") {
    return <LoginScreen onLogin={login} />;
  }

  if (view === "products") {
    return <ProductsScreen user={session.user} onBack={() => setView("dashboard")} onSessionExpired={expireSession} />;
  }

  return (
    <Dashboard
      user={session.user}
      onLogout={logout}
      onSessionExpired={expireSession}
      onOpenProducts={() => setView("products")}
    />
  );
}

const root = document.getElementById("root");
if (!root) throw new Error("StockPilot root element was not found.");

createRoot(root).render(<App />);
