import { useCallback, useEffect, useRef, useState } from "react";
import { APIError, stockpilotAPI, type ReplenishmentReview, type ReplenishmentReviewOutcome } from "./api";

type Props = { onSessionExpired: () => void };
const outcomes: Array<ReplenishmentReviewOutcome | "all"> = ["all", "accepted", "modified", "dismissed", "expired"];

export function ReplenishmentReviewHistoryPanel({ onSessionExpired }: Props) {
  const [outcome, setOutcome] = useState<ReplenishmentReviewOutcome | "all">("all");
  const [items, setItems] = useState<ReplenishmentReview[]>([]);
  const [state, setState] = useState<"loading" | "ready" | "error">("loading");
  const [error, setError] = useState("");
  const requestVersion = useRef(0);

  const load = useCallback(async () => {
    const version = ++requestVersion.current;
    setState("loading");
    try {
      const response = await stockpilotAPI.listReplenishmentReviews(outcome === "all" ? undefined : outcome, 50, 0);
      if (version !== requestVersion.current) return;
      setItems(response.items ?? []);
      setError("");
      setState("ready");
    } catch (err) {
      if (version !== requestVersion.current) return;
      if (err instanceof APIError && err.status === 401) { onSessionExpired(); return; }
      setError(err instanceof Error ? err.message : "Review history could not be loaded.");
      setState("error");
    }
  }, [onSessionExpired, outcome]);

  useEffect(() => { void load(); }, [load]);

  return <article className="panel wide" aria-labelledby="replenishment-review-history-title" aria-busy={state === "loading"}>
    <div className="panel-heading">
      <div><p className="eyebrow">Replenishment traceability</p><h2 id="replenishment-review-history-title">Review history</h2></div>
      <div className="topbar-actions">
        <label className="compact" htmlFor="replenishment-review-outcome">Outcome</label>
        <select id="replenishment-review-outcome" value={outcome} onChange={event => setOutcome(event.target.value as ReplenishmentReviewOutcome | "all")}>
          <option value="all">All</option>{outcomes.slice(1).map(value => <option key={value} value={value}>{value}</option>)}
        </select>
        <button className="secondary-button compact-button" type="button" onClick={() => void load()} disabled={state === "loading"} aria-label={state === "loading" ? "Refreshing review history" : "Refresh review history"}>Refresh</button>
      </div>
    </div>
    <p className="muted">Immutable snapshots of replenishment recommendations recorded during review. Historical values do not change when current stock or catalog data changes.</p>
    <div className="status" role="status" aria-live="polite">{state === "loading" ? "Loading review history…" : state === "ready" ? `${items.length} review${items.length === 1 ? "" : "s"} shown.` : "Review history could not be loaded."}</div>
    {state === "error" && <div className="notice error" role="alert"><span>{error}</span><button type="button" onClick={() => void load()}>Try again</button></div>}
    <div className="table-wrap">
      <table>
        <caption>Replenishment review history filtered by outcome</caption>
        <thead><tr><th scope="col">Reviewed</th><th scope="col">Product</th><th scope="col">Outcome</th><th scope="col" className="numeric">On hand</th><th scope="col" className="numeric">Suggested</th><th scope="col">Purchase order</th><th scope="col">Reviewer</th></tr></thead>
        <tbody>
          {items.map(item => <tr key={item.id}>
            <td><time dateTime={item.reviewedAt}>{formatDate(item.reviewedAt)}</time></td>
            <td><strong>{item.productName}</strong><br/><span className="compact mono">{item.sku}</span></td>
            <td><span className={`compact review-${item.outcome}`}>{item.outcome}</span></td>
            <td className="numeric">{item.onHand.toLocaleString()}</td>
            <td className="numeric">{item.suggestedQuantity.toLocaleString()} {item.unit}</td>
            <td className="mono">{item.purchaseOrderId ?? "—"}</td>
            <td className="mono">{item.reviewedBy}</td>
          </tr>)}
          {items.length === 0 && state === "ready" && <tr><td colSpan={7}>No replenishment reviews match this filter.</td></tr>}
          {state === "loading" && items.length === 0 && <tr><td colSpan={7}>Loading review history…</td></tr>}
        </tbody>
      </table>
    </div>
  </article>;
}

function formatDate(value: string) {
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return value;
  return new Intl.DateTimeFormat("en-IN", { dateStyle: "medium", timeStyle: "short" }).format(date);
}
