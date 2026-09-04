import { useCallback, useEffect, useState } from "react";
import { APIError } from "./api";
import { replenishmentReadiness, replenishmentReadinessCSV, type ReplenishmentReadinessReport, type ReplenishmentRisk } from "./replenishmentReadinessApi";

type Props = { canExport: boolean; onSessionExpired: () => void };
const emptyReport: ReplenishmentReadinessReport = { asOf: "", windowDays: 30, items: [] };

export function ReplenishmentReadinessPanel({ canExport, onSessionExpired }: Props) {
  const [report, setReport] = useState(emptyReport);
  const [state, setState] = useState<"loading" | "ready" | "error">("loading");
  const [error, setError] = useState("");

  const load = useCallback(async () => {
    setState("loading");
    try {
      setReport(await replenishmentReadiness(30, 500));
      setError("");
      setState("ready");
    } catch (err) {
      if (err instanceof APIError && err.status === 401) { onSessionExpired(); return; }
      setError(err instanceof Error ? err.message : "Replenishment readiness could not be loaded.");
      setState("error");
    }
  }, [onSessionExpired]);

  useEffect(() => { void load(); }, [load]);

  return <article className="panel wide">
    <div className="panel-heading">
      <div><p className="eyebrow">Replenishment readiness</p><h2>Review stock risk before ordering</h2></div>
      <div className="topbar-actions">
        {canExport && <button className="secondary-button compact-button" type="button" onClick={() => replenishmentReadinessCSV(report.windowDays)}>Export CSV</button>}
        <button className="secondary-button compact-button" type="button" onClick={() => void load()} disabled={state === "loading"}>Refresh</button>
      </div>
    </div>
    <p className="muted">Advisory view using the last {report.windowDays} days of outbound velocity. It does not create orders or mutate inventory.</p>
    {state === "error" && <div className="notice error" role="alert"><span>{error}</span><button type="button" onClick={() => void load()}>Try again</button></div>}
    <div className="table-wrap">
      <table>
        <thead><tr><th>Product</th><th>Risk</th><th className="numeric">On hand</th><th className="numeric">Avg/day</th><th className="numeric">Days cover</th><th className="numeric">Suggested</th></tr></thead>
        <tbody>
          {report.items.slice(0, 25).map(item => <tr key={item.productId}>
            <td><strong>{item.name}</strong><br/><span className="compact mono">{item.sku}</span></td>
            <td><RiskBadge risk={item.risk} /></td>
            <td className="numeric">{item.onHand.toLocaleString()}</td>
            <td className="numeric">{item.averageDailyOutbound.toFixed(2)}</td>
            <td className="numeric">{item.daysOfCover === undefined ? "—" : `${item.daysOfCover.toFixed(1)}d`}</td>
            <td className="numeric"><strong>{item.suggestedQuantity.toLocaleString()}</strong></td>
          </tr>)}
          {report.items.length === 0 && state === "ready" && <tr><td colSpan={6}>No reorder suggestions require review.</td></tr>}
          {state === "loading" && <tr><td colSpan={6}>Loading replenishment readiness…</td></tr>}
        </tbody>
      </table>
    </div>
  </article>;
}

function RiskBadge({ risk }: { risk: ReplenishmentRisk }) {
  const label = risk.replaceAll("_", " ");
  return <span className={`compact risk-${risk}`} aria-label={`Risk: ${label}`}>{label}</span>;
}
