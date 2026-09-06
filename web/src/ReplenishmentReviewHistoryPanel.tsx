import { useCallback, useEffect, useState } from "react";
import { APIError, stockpilotAPI, type ReplenishmentReview, type ReplenishmentReviewOutcome } from "./api";

type Props = { onSessionExpired: () => void };
const outcomes: Array<ReplenishmentReviewOutcome | "all"> = ["all", "accepted", "modified", "dismissed", "expired"];

export function ReplenishmentReviewHistoryPanel({ onSessionExpired }: Props) {
  const [outcome, setOutcome] = useState<ReplenishmentReviewOutcome | "all">("all");
  const [items, setItems] = useState<ReplenishmentReview[]>([]);
  const [state, setState] = useState<"loading" | "ready" | "error">("loading");
  const [error, setError] = useState("");

  const load = useCallback(async () => {
    setState("loading");
    try {
      const response = await stockpilotAPI.listReplenishmentReviews(outcome === "all" ? undefined : outcome, 50, 0);
      setItems(response.items ?? []);
      setError("");
      setState("ready");
    } catch (err) {
      if (err instanceof APIError && err.status === 401) { onSessionExpired(); return; }
      setError(err instanceof Error ? err.message : "Review history could not be loaded.");
      setState("error");
    }
  }, [onSessionExpired, outcome]);

  useEffect(() => { void load(); }, [load]);

  return <article className="panel wide">
    <div className="panel-heading">
      <div><p className="eyebrow">Replenishment traceability</p><h2>Review history</h2></div>
      <div className="topbar-actions">
        <label className="compact">Outcome <select value={outcome} onChange={event => setOutcome(event.target.value as ReplenishmentReviewOutcome | "all")}><option value="all">All</option>{outcomes.slice(1).map(value => <option key={value} value={value}>{value}</option>)}</select></label>
        <button className="secondary-button compact-button" type="button" onClick={() => void load()} disabled={state === "loading"}>Refresh</button>
      </div>
    </div>
    <p className="muted">Immutable snapshots of replenishment recommendations recorded during review. Historical values do not change when current stock or catalog data changes.</p>
    {state === "error" && <div className="notice error" role="alert"><span>{error}</span><button type="button" onClick={() => void load()}>Try again</button></div>}
    <div className="table-wrap">
      <table>
        <thead><tr><th>Reviewed</th><th>Product</th><th>Outcome</th><th className="numeric">On hand</th><th className="numeric">Suggested</th><th>Purchase order</th><th>Reviewer</th></tr></thead>
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
          {state === "loading" && <tr><td colSpan={7}>Loading review history…</td></tr>}
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
