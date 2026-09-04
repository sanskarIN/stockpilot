import { useCallback, useEffect, useState } from "react";
import { APIError, stockpilotAPI } from "./api";
import { inventoryAging } from "./inventoryAgingApi";
import type { InventoryValuationReport, InventorySummary, PurchasingSummary, User, InventoryAgingReport } from "./types";

type Props = { user: User; onBack: () => void; onSessionExpired: () => void };
type LoadState = { kind: "loading" } | { kind: "ready" } | { kind: "error"; message: string };
const emptyValuation: InventoryValuationReport = { items: [], totals: [] };
const emptyAging: InventoryAgingReport = { items: [] };

export function ReportsScreen({ user, onBack, onSessionExpired }: Props) {
  const [inventory, setInventory] = useState<InventorySummary | null>(null);
  const [purchasing, setPurchasing] = useState<PurchasingSummary | null>(null);
  const [valuation, setValuation] = useState<InventoryValuationReport>(emptyValuation);
  const [aging, setAging] = useState<InventoryAgingReport>(emptyAging);
  const [state, setState] = useState<LoadState>({ kind: "loading" });
  const canExport = user.role !== "operator";

  const load = useCallback(async () => {
    setState({ kind: "loading" });
    try {
      const [overview, valuationReport, agingReport] = await Promise.all([
        stockpilotAPI.reportOverview(), stockpilotAPI.inventoryValuation(), inventoryAging(500),
      ]);
      setInventory(overview.inventory); setPurchasing(overview.purchasing); setValuation(valuationReport); setAging(agingReport); setState({ kind: "ready" });
    } catch (error) {
      if (error instanceof APIError && error.status === 401) { onSessionExpired(); return; }
      setState({ kind: "error", message: error instanceof Error ? error.message : "Reports could not be loaded." });
    }
  }, [onSessionExpired]);

  useEffect(() => { void load(); }, [load]);
  function exportValuation() { window.location.assign("/api/v1/reports/inventory-valuation/export.csv?limit=5000"); }
  function exportAging() { window.location.assign("/api/v1/reports/inventory-aging?format=csv&limit=5000"); }

  return <div className="app-shell"><main className="content report-page" tabIndex={-1}>
    <header className="topbar"><div><p className="eyebrow">Operational reporting</p><h1>Reports & analytics</h1><p className="muted">Read-only reporting built from authoritative catalog, inventory, purchasing, valuation, and movement data.</p></div><div className="topbar-actions"><button className="secondary-button" type="button" onClick={onBack}>Back to dashboard</button><button className="primary-button" type="button" onClick={() => void load()} disabled={state.kind === "loading"}>Refresh</button></div></header>
    <div className="status-region" aria-live="polite">{state.kind === "loading" && <p className="notice">Refreshing reports…</p>}{state.kind === "error" && <div className="notice error" role="alert"><span>{state.message}</span><button type="button" onClick={() => void load()}>Try again</button></div>}</div>
    <section className="metrics"><Metric label="Active products" value={(inventory?.activeProductCount ?? 0).toLocaleString()} detail="Current active catalog"/><Metric label="Units on hand" value={(inventory?.totalUnits ?? 0).toLocaleString()} detail="Across inventory balances"/><Metric label="Low-stock balances" value={(inventory?.lowStockBalanceCount ?? 0).toLocaleString()} detail="Balances at or below threshold"/><Metric label="Outstanding PO units" value={(purchasing?.outstandingUnits ?? 0).toLocaleString()} detail="Not yet fully received" tone={(purchasing?.outstandingUnits ?? 0) > 0 ? "warning" : "good"}/></section>
    <section className="dashboard-grid report-grid">
      <article className="panel wide"><div className="panel-heading"><div><p className="eyebrow">Inventory valuation</p><h2>Current on-hand value</h2></div>{canExport && <button className="secondary-button compact-button" type="button" onClick={exportValuation}>Export CSV</button>}</div><p className="muted">Values are reported in minor currency units by product and grouped by currency.</p><ul className="activity-list">{valuation.totals.map(total => <li key={total.currency}><span><strong>{total.currency}</strong><small>Total on-hand valuation</small></span><strong>{formatMoney(total.valueMinor, total.currency)}</strong></li>)}{valuation.totals.length === 0 && state.kind !== "loading" && <li><span className="muted">No valuation data available.</span></li>}</ul></article>
      <article className="panel"><div className="panel-heading"><div><p className="eyebrow">Purchasing</p><h2>Order pipeline</h2></div></div><ul className="activity-list"><Row label="Draft" value={purchasing?.draftOrders ?? 0}/><Row label="Ordered" value={purchasing?.orderedOrders ?? 0}/><Row label="Partially received" value={purchasing?.partiallyReceivedOrders ?? 0}/><Row label="Received" value={purchasing?.receivedOrders ?? 0}/><Row label="Cancelled" value={purchasing?.cancelledOrders ?? 0}/></ul></article>
      <article className="panel wide"><div className="panel-heading"><div><p className="eyebrow">Inventory aging</p><h2>Time since last movement</h2></div>{canExport && <button className="secondary-button compact-button" type="button" onClick={exportAging}>Export CSV</button>}</div><p className="muted">Positive balances grouped into deterministic 0–30, 31–60, 61–90, 91–180, and 181+ day buckets.</p><div className="table-wrap"><table><thead><tr><th>Product</th><th>Location</th><th className="numeric">Quantity</th><th className="numeric">Age</th><th>Bucket</th></tr></thead><tbody>{aging.items.slice(0, 25).map(item => <tr key={`${item.productId}-${item.locationId}-${item.lotId}`}><td><strong>{item.name}</strong><br/><span className="compact mono">{item.sku}</span></td><td className="mono">{item.locationId}</td><td className="numeric">{item.quantity.toLocaleString()}</td><td className="numeric">{item.ageDays.toLocaleString()}d</td><td><strong>{item.bucket}</strong></td></tr>)}{aging.items.length === 0 && state.kind !== "loading" && <tr><td colSpan={5}>No aging data available.</td></tr>}</tbody></table></div></article>
      <article className="panel wide"><div className="panel-heading"><div><p className="eyebrow">Highest value items</p><h2>Valuation breakdown</h2></div></div><div className="table-wrap"><table><thead><tr><th>Product</th><th>SKU</th><th className="numeric">On hand</th><th className="numeric">Unit cost</th><th className="numeric">Value</th></tr></thead><tbody>{valuation.items.slice(0, 25).map(item => <tr key={item.productId}><td><strong>{item.name}</strong><br/><span className="compact">{item.unit}</span></td><td className="mono">{item.sku}</td><td className="numeric">{item.onHand.toLocaleString()}</td><td className="numeric">{formatMoney(item.unitCostMinor, item.currency)}</td><td className="numeric"><strong>{formatMoney(item.valueMinor, item.currency)}</strong></td></tr>)}{valuation.items.length === 0 && state.kind !== "loading" && <tr><td colSpan={5}>No valuation items available.</td></tr>}</tbody></table></div></article>
    </section>
  </main></div>;
}
function Metric({ label, value, detail, tone = "neutral" }: { label: string; value: string; detail: string; tone?: "neutral" | "good" | "warning" }) { return <article className={`metric ${tone}`}><p>{label}</p><strong>{value}</strong><span>{detail}</span></article>; }
function Row({ label, value }: { label: string; value: number }) { return <li><span><strong>{label}</strong></span><strong>{value.toLocaleString()}</strong></li>; }
function formatMoney(minor: number, currency: string) { try { return new Intl.NumberFormat("en-IN", { style: "currency", currency: currency || "INR", maximumFractionDigits: 2 }).format(minor / 100); } catch { return `${currency || "INR"} ${(minor / 100).toFixed(2)}`; } }
