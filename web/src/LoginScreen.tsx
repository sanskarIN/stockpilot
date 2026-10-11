import { useState, type FormEvent } from "react";

export function LoginScreen({ onLogin }: { onLogin: (email: string, password: string) => Promise<void> }) {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [pending, setPending] = useState(false);
  const [error, setError] = useState("");

  async function submit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setError("");
    setPending(true);
    try {
      await onLogin(email.trim(), password);
      setPassword("");
    } catch (caught) {
      setError(caught instanceof Error ? caught.message : "Sign in failed. Please try again.");
    } finally {
      setPending(false);
    }
  }

  return (
    <main className="auth-page">
      <section className="auth-card" aria-labelledby="login-title">
        <a className="brand auth-brand" href="/" aria-label="StockPilot home">
          <span className="brand-mark" aria-hidden="true">SP</span>
          <span>
            <strong>StockPilot</strong>
            <small>Inventory control</small>
          </span>
        </a>

        <div className="auth-copy">
          <p className="eyebrow">Secure workspace</p>
          <h1 id="login-title">Sign in to StockPilot</h1>
          <p className="muted">Use the account created by your StockPilot administrator.</p>
        </div>

        <form className="auth-form" onSubmit={submit}>
          <label>
            <span>Email</span>
            <input
              type="email"
              name="email"
              autoComplete="username"
              inputMode="email"
              value={email}
              onChange={(event) => setEmail(event.target.value)}
              required
              maxLength={254}
              disabled={pending}
            />
          </label>
          <label>
            <span>Password</span>
            <input
              type="password"
              name="password"
              autoComplete="current-password"
              value={password}
              onChange={(event) => setPassword(event.target.value)}
              required
              minLength={12}
              maxLength={72}
              disabled={pending}
            />
          </label>
          {error && <p className="form-error" role="alert">{error}</p>}
          <button className="primary-button auth-submit" type="submit" disabled={pending || !email.trim() || !password}>
            {pending ? "Signing in…" : "Sign in"}
          </button>
        </form>

        <div className="auth-help">
          <p>First installation? An operator must run the documented administrator bootstrap command before sign-in.</p>
          <p><strong>Made by the Sanskar</strong> · <a href="mailto:supportramsandesh@gmail.com">Support</a></p>
        </div>
      </section>
    </main>
  );
}
