import React, { useCallback, useEffect, useMemo, useState } from "react";
import { createRoot } from "react-dom/client";
import "./styles.css";

type Product = {
  id: string;
  sku: string;
  name: string;
  unit: string;
  unitCostMinor: number;
  currency: string;
  reorderPoint: number;
  active: boolean;
};

type StockBalance = {
  productId: string;
  locationId: string;
  lotId?: string;
  quantity: number;
  updatedAt: string;
};

type PurchaseOrder = {
  id: string;
  number: string;
  status: string;
  currency: string;
  createdAt: string;
};

type ListResponse<T> = { items: T[] };

type LoadState =
  | { kind: "loading" }
  | { kind: "ready" }
  | { kind: "error"; message: string };

const navigation = ["Overview", "Products", "Inventory", "Purchase orders", "Warehouses", "Reports"];

function App() {
  const [products, setProducts] = useState<Product[]>([]);
  const [lowStock, setLowStock] = useState<StockBalance[]>([]);
  const [orders, setOrders] = useState<PurchaseOrder[]>([]);
  const [state, setState] = useState<LoadState>({ kind: "loading" });
  const [theme, setTheme] = useState<"light" | "dark">(() =>
    window.matchMedia?.("(prefers-color-scheme: dark)").matches ? "dark" : "light"
  );

  useEffect(() => {
    document.documentElement.dataset.theme = theme;
  }, [theme]);

  const load = useCallback(async () => {
    setState({ kind: "loading" });
    try {
      const [productsResponse, lowStockResponse, ordersResponse] = await Promise.all([
        fetch("/api/v1/products?active=true&limit=100", { headers: { Accept: "application/json" } }),
        fetch("/api/v1/inventory/low-stock?limit=100", { headers: { Accept: "application/json" } }),
        fetch("/api/v1/orders?limit=100", { headers: { Accept: "application/json" } })
      ]);
      if (!productsResponse.ok || !lowStockResponse.ok || !ordersResponse.ok) {
        throw new Error("StockPilot could not load the dashboard data.");
      }
      const [productData, lowStockData, orderData] = await Promise.all([
        productsResponse.json() as Promise<ListResponse<Product>>,
        lowStockResponse.json() as Promise<ListResponse<StockBalance>>,
        ordersResponse.json() as Promise<ListResponse<PurchaseOrder>>
      ]);
      setProducts(productData.items ?? []);
      setLowStock(lowStockData.items ?? []);
      setOrders(orderData.items ?? []);
      setState({ kind: "ready" });
    } catch (error) {
      setState({
        kind: "error",
        message: error instanceof Error ? error.message : "StockPilot could not load dashboard data."
      });
    }
  }, []);

  useEffect(() => {
    void load();
  }, [load]);

  const productsByID = useMemo(() => new Map(products.map((product) => [product.id, product])), [products]);
  const openOrders = orders.filter((order) => !["received", "cancelled"].includes(order.status));
  const healthyProducts = Math.max(products.length - new Set(lowStock.map((item) => item.productId)).size, 0);
  const healthPercent = products.length === 0 ? 100 : Math.round((healthyProducts / products.length) * 100);

  return (
    <div className="app-shell">
      <aside className="sidebar" aria-label="Primary navigation">
        <a className="brand" href="#main" aria-label="StockPilot home">
          <span className="brand-mark" aria-hidden="true">SP</span>
          <span>
            <strong>StockPilot</strong>
            <small>Inventory control</small>
          </span>
        </a>
        <nav>
          {navigation.map((item, index) => (
            <a className={index === 0 ? "nav-link active" : "nav-link"} href={index === 0 ? "#main" : `#${item.toLowerCase().replaceAll(" ", "-")}`} key={item}>
              <span aria-hidden="true" className="nav-dot" />
              {item}
            </a>
          ))}
        </nav>
        <div className="sidebar-footer">
          <span>Made by the Sanskar</span>
          <a href="https://github.com/sanskarIN" rel="noreferrer">GitHub</a>
        </div>
      </aside>

      <main id="main" className="content" tabIndex={-1}>
        <header className="topbar">
          <div>
            <p className="eyebrow">Operations control center</p>
            <h1>Inventory overview</h1>
            <p className="muted">Live stock health, replenishment signals, and purchasing activity.</p>
          </div>
          <div className="topbar-actions">
            <button className="secondary-button" type="button" onClick={() => setTheme((value) => value === "light" ? "dark" : "light")}>
              {theme === "light" ? "Dark mode" : "Light mode"}
            </button>
            <button className="primary-button" type="button" onClick={() => void load()} disabled={state.kind === "loading"}>
              Refresh
            </button>
          </div>
        </header>

        <div className="status-region" aria-live="polite">
          {state.kind === "loading" && <p className="notice">Refreshing inventory data…</p>}
          {state.kind === "error" && (
            <div className="notice error" role="alert">
              <span>{state.message}</span>
              <button type="button" onClick={() => void load()}>Try again</button>
            </div>
          )}
        </div>

        <section className="metrics" aria-label="Inventory summary">
          <Metric label="Active products" value={products.length.toLocaleString()} detail="Tracked catalog items" />
          <Metric label="Low-stock balances" value={lowStock.length.toLocaleString()} detail="At or below reorder point" tone={lowStock.length > 0 ? "warning" : "good"} />
          <Metric label="Open purchase orders" value={openOrders.length.toLocaleString()} detail="Draft, ordered, or partial" />
          <Metric label="Stock health" value={`${healthPercent}%`} detail={`${healthyProducts} products above threshold`} tone={healthPercent >= 90 ? "good" : "warning"} />
        </section>

        <section className="dashboard-grid">
          <article className="panel wide" aria-labelledby="low-stock-title">
            <div className="panel-heading">
              <div>
                <p className="eyebrow">Replenishment</p>
                <h2 id="low-stock-title">Low-stock attention</h2>
              </div>
              <span className="count-pill">{lowStock.length}</span>
            </div>
            {lowStock.length === 0 && state.kind !== "loading" ? (
              <EmptyState title="No low-stock balances" body="Tracked balances are currently above their product reorder points." />
            ) : (
              <div className="table-wrap">
                <table>
                  <thead>
                    <tr><th>Product</th><th>SKU</th><th>Location</th><th className="numeric">On hand</th><th>Updated</th></tr>
                  </thead>
                  <tbody>
                    {lowStock.slice(0, 8).map((balance) => {
                      const product = productsByID.get(balance.productId);
                      return (
                        <tr key={`${balance.productId}:${balance.locationId}:${balance.lotId ?? ""}`}>
                          <td><strong>{product?.name ?? "Unknown product"}</strong></td>
                          <td className="mono">{product?.sku ?? "—"}</td>
                          <td className="mono compact">{shortID(balance.locationId)}</td>
                          <td className="numeric"><span className="status-badge warning">{balance.quantity.toLocaleString()}</span></td>
                          <td>{formatDate(balance.updatedAt)}</td>
                        </tr>
                      );
                    })}
                  </tbody>
                </table>
              </div>
            )}
          </article>

          <article className="panel" aria-labelledby="orders-title">
            <div className="panel-heading">
              <div>
                <p className="eyebrow">Purchasing</p>
                <h2 id="orders-title">Open orders</h2>
              </div>
            </div>
            {openOrders.length === 0 && state.kind !== "loading" ? (
              <EmptyState title="No open purchase orders" body="New orders will appear here as soon as they are created." />
            ) : (
              <ul className="activity-list">
                {openOrders.slice(0, 6).map((order) => (
                  <li key={order.id}>
                    <span><strong>{order.number}</strong><small>{formatDate(order.createdAt)}</small></span>
                    <span className="status-badge">{order.status.replaceAll("_", " ")}</span>
                  </li>
                ))}
              </ul>
            )}
          </article>

          <article className="panel" aria-labelledby="catalog-title">
            <div className="panel-heading">
              <div>
                <p className="eyebrow">Catalog</p>
                <h2 id="catalog-title">Recent products</h2>
              </div>
            </div>
            {products.length === 0 && state.kind !== "loading" ? (
              <EmptyState title="Catalog is empty" body="Create the first product through the API to begin tracking inventory." />
            ) : (
              <ul className="activity-list">
                {products.slice(0, 6).map((product) => (
                  <li key={product.id}>
                    <span><strong>{product.name}</strong><small className="mono">{product.sku}</small></span>
                    <span>{formatMoney(product.unitCostMinor, product.currency)}</span>
                  </li>
                ))}
              </ul>
            )}
          </article>
        </section>
      </main>
    </div>
  );
}

function Metric({ label, value, detail, tone = "neutral" }: { label: string; value: string; detail: string; tone?: "neutral" | "good" | "warning" }) {
  return (
    <article className={`metric ${tone}`}>
      <p>{label}</p>
      <strong>{value}</strong>
      <span>{detail}</span>
    </article>
  );
}

function EmptyState({ title, body }: { title: string; body: string }) {
  return <div className="empty-state"><strong>{title}</strong><p>{body}</p></div>;
}

function shortID(value: string) {
  if (value.length <= 14) return value;
  return `${value.slice(0, 7)}…${value.slice(-5)}`;
}

function formatDate(value: string) {
  const date = new Date(value);
  return Number.isNaN(date.getTime()) ? "—" : new Intl.DateTimeFormat("en-IN", { dateStyle: "medium" }).format(date);
}

function formatMoney(minor: number, currency: string) {
  try {
    return new Intl.NumberFormat("en-IN", { style: "currency", currency: currency || "INR", maximumFractionDigits: 2 }).format(minor / 100);
  } catch {
    return `${currency || "INR"} ${(minor / 100).toFixed(2)}`;
  }
}

const root = document.getElementById("root");
if (!root) throw new Error("StockPilot root element was not found.");
createRoot(root).render(<React.StrictMode><App /></React.StrictMode>);
