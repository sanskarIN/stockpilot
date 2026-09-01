import { useCallback, useEffect, useMemo, useState, type FormEvent } from "react";
import { APIError, stockpilotAPI } from "./api";
import type { Location, Product, PurchaseOrder, Supplier, User, Warehouse } from "./types";

type Props = { user: User; onBack: () => void; onSessionExpired: () => void };
type OrderDraft = { number: string; supplierId: string; warehouseId: string; currency: string; notes: string; productId: string; quantity: number; unitCostMinor: number };
const blank: OrderDraft = { number: "", supplierId: "", warehouseId: "", currency: "INR", notes: "", productId: "", quantity: 1, unitCostMinor: 0 };

export function PurchaseOrdersScreen({ user, onBack, onSessionExpired }: Props) {
  const [orders, setOrders] = useState<PurchaseOrder[]>([]);
  const [suppliers, setSuppliers] = useState<Supplier[]>([]);
  const [warehouses, setWarehouses] = useState<Warehouse[]>([]);
  const [locations, setLocations] = useState<Location[]>([]);
  const [products, setProducts] = useState<Product[]>([]);
  const [selectedOrder, setSelectedOrder] = useState<PurchaseOrder | null>(null);
  const [draft, setDraft] = useState<OrderDraft>(blank);
  const [receiveQuantity, setReceiveQuantity] = useState(1);
  const [receiveLocationId, setReceiveLocationId] = useState("");
  const [status, setStatus] = useState<"loading" | "ready" | "error">("loading");
  const [message, setMessage] = useState("");
  const [saving, setSaving] = useState(false);
  const [showCreate, setShowCreate] = useState(false);
  const canWrite = user.role === "admin" || user.role === "manager" || user.role === "operator";

  const load = useCallback(async () => {
    setStatus("loading");
    try {
      const [orderData, supplierData, warehouseData, locationData, productData] = await Promise.all([
        stockpilotAPI.listOrders(), stockpilotAPI.listSuppliers(), stockpilotAPI.listWarehouses(), stockpilotAPI.listLocations(), stockpilotAPI.listProducts("")
      ]);
      setOrders(orderData.items ?? []); setSuppliers(supplierData.items ?? []); setWarehouses(warehouseData.items ?? []); setLocations(locationData.items ?? []); setProducts(productData.items ?? []);
      setDraft((current) => ({ ...current, supplierId: current.supplierId || supplierData.items?.[0]?.id || "", warehouseId: current.warehouseId || warehouseData.items?.[0]?.id || "", productId: current.productId || productData.items?.[0]?.id || "" }));
      if (!receiveLocationId && locationData.items?.[0]) setReceiveLocationId(locationData.items[0].id);
      setStatus("ready"); setMessage("");
    } catch (error: unknown) {
      if (error instanceof APIError && error.status === 401) { onSessionExpired(); return; }
      setStatus("error"); setMessage(error instanceof Error ? error.message : "Could not load purchasing data.");
    }
  }, [onSessionExpired, receiveLocationId]);

  useEffect(() => { void load(); }, [load]);

  const supplierNames = useMemo(() => new Map(suppliers.map((item) => [item.id, item.name])), [suppliers]);
  const warehouseNames = useMemo(() => new Map(warehouses.map((item) => [item.id, item.name])), [warehouses]);
  const selectedProduct = products.find((item) => item.id === draft.productId);

  async function createOrder(event: FormEvent) {
    event.preventDefault(); setMessage("");
    if (!canWrite) { setMessage("Your role cannot create purchase orders."); return; }
    if (draft.number.trim().length < 2 || draft.number.trim().length > 64) { setMessage("Order number must be 2-64 characters."); return; }
    if (!draft.supplierId || !draft.warehouseId || !draft.productId) { setMessage("Supplier, warehouse, and product are required."); return; }
    if (!Number.isSafeInteger(draft.quantity) || draft.quantity <= 0) { setMessage("Quantity must be a positive whole number."); return; }
    if (!Number.isSafeInteger(draft.unitCostMinor) || draft.unitCostMinor < 0) { setMessage("Unit cost must be a non-negative whole number."); return; }
    if (!/^[A-Z]{3}$/.test(draft.currency.trim().toUpperCase())) { setMessage("Currency must be a three-letter code."); return; }
    setSaving(true);
    try {
      const created = await stockpilotAPI.createOrder({ number: draft.number.trim(), supplierId: draft.supplierId, warehouseId: draft.warehouseId, status: "draft", currency: draft.currency.trim().toUpperCase(), notes: draft.notes.trim() || undefined, lines: [{ productId: draft.productId, quantity: draft.quantity, received: 0, unitCostMinor: draft.unitCostMinor }] });
      setSelectedOrder(created); setDraft(blank); setShowCreate(false); setMessage("Purchase order created."); await load();
    } catch (error: unknown) {
      if (error instanceof APIError && error.status === 401) { onSessionExpired(); return; }
      setMessage(error instanceof Error ? error.message : "Could not create the purchase order.");
    } finally { setSaving(false); }
  }

  async function receiveLine() {
    if (!selectedOrder || !canWrite) { setMessage("Select an order and use a role with purchasing access."); return; }
    const line = selectedOrder.lines?.[0];
    if (!line) { setMessage("This order has no receivable lines."); return; }
    const remaining = Math.max(line.quantity - line.received, 0);
    if (!Number.isSafeInteger(receiveQuantity) || receiveQuantity <= 0 || receiveQuantity > remaining) { setMessage(`Receive between 1 and ${remaining.toLocaleString()} units.`); return; }
    if (!receiveLocationId) { setMessage("Choose a receiving location."); return; }
    setSaving(true); setMessage("");
    try {
      const updated = await stockpilotAPI.receiveOrderLine(selectedOrder.id, line.id, { quantity: receiveQuantity, locationId: receiveLocationId });
      setSelectedOrder(updated); setReceiveQuantity(1); setMessage("Receipt recorded and inventory was updated."); await load();
    } catch (error: unknown) {
      if (error instanceof APIError && error.status === 401) { onSessionExpired(); return; }
      setMessage(error instanceof Error ? error.message : "Could not receive the order line.");
    } finally { setSaving(false); }
  }

  return <main className="catalog-page" aria-labelledby="purchase-orders-title">
    <header className="catalog-header"><div><button className="secondary-button" type="button" onClick={onBack}>← Overview</button><p className="eyebrow">Purchasing</p><h1 id="purchase-orders-title">Purchase orders</h1><p className="muted">Create draft purchase orders and receive ordered quantities into a controlled inventory location.</p></div><div className="topbar-actions"><button className="secondary-button" type="button" onClick={() => void load()} disabled={status === "loading"}>Refresh</button>{canWrite && <button className="primary-button" type="button" onClick={() => { setDraft(blank); setShowCreate(true); setMessage(""); }}>New order</button>}</div></header>
    {message && <p className="notice" role="status">{message}</p>}
    <section className="dashboard-grid">
      <article className="panel wide"><div className="panel-heading"><div><p className="eyebrow">Order register</p><h2>Recent purchase orders</h2></div><span className="count-pill">{orders.length}</span></div>{orders.length === 0 && status !== "loading" ? <div className="empty-state"><strong>No purchase orders</strong><p>Create a draft order to begin the purchasing workflow.</p></div> : <div className="table-wrap"><table><thead><tr><th>Order</th><th>Supplier</th><th>Warehouse</th><th>Status</th><th>Created</th><th>Open</th></tr></thead><tbody>{orders.map((order) => <tr key={order.id}><td><strong>{order.number}</strong><small className="table-secondary">{order.currency}</small></td><td>{supplierNames.get(order.supplierId) ?? "Unknown"}</td><td>{warehouseNames.get(order.warehouseId) ?? "Unknown"}</td><td><span className="status-badge">{order.status.replaceAll("_", " ")}</span></td><td>{formatDate(order.createdAt)}</td><td><button className="secondary-button compact-button" type="button" onClick={() => { setSelectedOrder(order); setMessage(""); }}>Open</button></td></tr>)}</tbody></table></div>}</article>
      {selectedOrder && <article className="panel"><div className="panel-heading"><div><p className="eyebrow">Order detail</p><h2>{selectedOrder.number}</h2></div><span className="status-badge">{selectedOrder.status.replaceAll("_", " ")}</span></div><p className="muted">Supplier: {supplierNames.get(selectedOrder.supplierId) ?? "Unknown"}<br />Warehouse: {warehouseNames.get(selectedOrder.warehouseId) ?? "Unknown"}</p><div className="order-lines">{(selectedOrder.lines ?? []).map((line) => <div className="order-line" key={line.id}><strong>{products.find((item) => item.id === line.productId)?.name ?? line.productId}</strong><span>{line.received.toLocaleString()} / {line.quantity.toLocaleString()} received</span></div>)}</div>{selectedOrder.lines?.[0] && selectedOrder.lines[0].received < selectedOrder.lines[0].quantity && canWrite && <div className="receive-box"><p className="eyebrow">Receive</p><div className="form-grid"><label><span>Quantity</span><input type="number" min="1" max={selectedOrder.lines[0].quantity - selectedOrder.lines[0].received} step="1" value={receiveQuantity} onChange={(event) => setReceiveQuantity(Number(event.target.value))} /></label><label><span>Location</span><select value={receiveLocationId} onChange={(event) => setReceiveLocationId(event.target.value)}>{locations.map((location) => <option value={location.id} key={location.id}>{location.name} ({location.code})</option>)}</select></label></div><button className="primary-button" type="button" onClick={() => void receiveLine()} disabled={saving}>Receive into inventory</button></div>}</article>}
    </section>
    {showCreate && <div className="modal-backdrop" role="presentation"><section className="modal-card" role="dialog" aria-modal="true" aria-labelledby="new-order-title"><div className="panel-heading"><div><p className="eyebrow">New order</p><h2 id="new-order-title">Create purchase order</h2></div><button className="secondary-button" type="button" onClick={() => setShowCreate(false)} disabled={saving}>Close</button></div><form className="catalog-form" onSubmit={createOrder}><div className="form-grid"><label><span>Order number</span><input value={draft.number} onChange={(event) => setDraft((x) => ({ ...x, number: event.target.value }))} maxLength={64} required /></label><label><span>Supplier</span><select value={draft.supplierId} onChange={(event) => setDraft((x) => ({ ...x, supplierId: event.target.value }))} required>{suppliers.map((item) => <option value={item.id} key={item.id}>{item.name} ({item.code})</option>)}</select></label><label><span>Warehouse</span><select value={draft.warehouseId} onChange={(event) => setDraft((x) => ({ ...x, warehouseId: event.target.value }))} required>{warehouses.map((item) => <option value={item.id} key={item.id}>{item.name}</option>)}</select></label><label><span>Product</span><select value={draft.productId} onChange={(event) => { const id = event.target.value; const product = products.find((item) => item.id === id); setDraft((x) => ({ ...x, productId: id, unitCostMinor: product?.unitCostMinor ?? 0, currency: product?.currency ?? x.currency })); }} required>{products.map((item) => <option value={item.id} key={item.id}>{item.name} ({item.sku})</option>)}</select></label><label><span>Quantity</span><input type="number" min="1" step="1" value={draft.quantity} onChange={(event) => setDraft((x) => ({ ...x, quantity: Number(event.target.value) }))} required /></label><label><span>Unit cost (minor units)</span><input type="number" min="0" step="1" value={draft.unitCostMinor} onChange={(event) => setDraft((x) => ({ ...x, unitCostMinor: Number(event.target.value) }))} required /></label><label><span>Currency</span><input value={draft.currency} onChange={(event) => setDraft((x) => ({ ...x, currency: event.target.value }))} maxLength={3} required /></label></div><label><span>Notes</span><textarea rows={3} maxLength={1000} value={draft.notes} onChange={(event) => setDraft((x) => ({ ...x, notes: event.target.value }))} /></label><p className="muted">Estimated total: {formatMoney(draft.quantity * draft.unitCostMinor, draft.currency)}</p><div className="form-actions"><button className="secondary-button" type="button" onClick={() => setShowCreate(false)} disabled={saving}>Cancel</button><button className="primary-button" type="submit" disabled={saving || !selectedProduct}>{saving ? "Creating…" : "Create draft order"}</button></div></form></section></div>}
  </main>;
}

function formatDate(value: string) { const date = new Date(value); return Number.isNaN(date.getTime()) ? "—" : new Intl.DateTimeFormat("en-IN", { dateStyle: "medium" }).format(date); }
function formatMoney(minor: number, currency: string) { try { return new Intl.NumberFormat("en-IN", { style: "currency", currency: currency || "INR", maximumFractionDigits: 2 }).format(minor / 100); } catch { return `${currency || "INR"} ${(minor / 100).toFixed(2)}`; } }
