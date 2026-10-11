import { useCallback, useEffect, useMemo, useState, type FormEvent } from "react";
import { APIError, stockpilotAPI } from "./api";
import { BarcodeScanner } from "./BarcodeScanner";
import type { Category, Product, Supplier, User } from "./types";
import "./catalog.css";
import "./barcodeScanner.css";

type ProductDraft = Omit<Product, "id" | "createdAt" | "updatedAt">;

type LoadState =
  | { kind: "loading" }
  | { kind: "ready" }
  | { kind: "error"; message: string };

const blankProduct: ProductDraft = {
  sku: "",
  name: "",
  description: "",
  categoryId: "",
  supplierId: "",
  barcode: "",
  unit: "pcs",
  unitCostMinor: 0,
  currency: "INR",
  reorderPoint: 0,
  reorderQuantity: 0,
  trackLots: false,
  trackExpiry: false,
  active: true
};

export function ProductsScreen({ user, onBack, onSessionExpired }: { user: User; onBack: () => void; onSessionExpired: () => void }) {
  const [products, setProducts] = useState<Product[]>([]);
  const [categories, setCategories] = useState<Category[]>([]);
  const [suppliers, setSuppliers] = useState<Supplier[]>([]);
  const [state, setState] = useState<LoadState>({ kind: "loading" });
  const [query, setQuery] = useState("");
  const [editing, setEditing] = useState<Product | null>(null);
  const [formOpen, setFormOpen] = useState(false);
  const [scannerOpen, setScannerOpen] = useState(false);
  const [draft, setDraft] = useState<ProductDraft>(blankProduct);
  const [saving, setSaving] = useState(false);
  const [formError, setFormError] = useState("");
  const canWrite = user.role === "admin" || user.role === "manager";

  const load = useCallback(async () => {
    setState({ kind: "loading" });
    try {
      const [productResponse, categoryResponse, supplierResponse] = await Promise.all([
        stockpilotAPI.listProducts(query),
        stockpilotAPI.listCategories(),
        stockpilotAPI.listSuppliers()
      ]);
      setProducts(productResponse.items ?? []);
      setCategories(categoryResponse.items ?? []);
      setSuppliers(supplierResponse.items ?? []);
      setState({ kind: "ready" });
    } catch (error) {
      if (error instanceof APIError && error.status === 401) {
        onSessionExpired();
        return;
      }
      setState({ kind: "error", message: error instanceof Error ? error.message : "Could not load the catalog." });
    }
  }, [onSessionExpired, query]);

  useEffect(() => {
    const timer = window.setTimeout(() => void load(), 180);
    return () => window.clearTimeout(timer);
  }, [load]);

  const categoryNames = useMemo(() => new Map(categories.map((category) => [category.id, category.name])), [categories]);
  const supplierNames = useMemo(() => new Map(suppliers.map((supplier) => [supplier.id, supplier.name])), [suppliers]);

  function startCreate() {
    setEditing(null);
    setDraft(blankProduct);
    setFormError("");
    setFormOpen(true);
  }

  function startEdit(product: Product) {
    setEditing(product);
    setDraft({
      sku: product.sku,
      name: product.name,
      description: product.description ?? "",
      categoryId: product.categoryId ?? "",
      supplierId: product.supplierId ?? "",
      barcode: product.barcode ?? "",
      unit: product.unit,
      unitCostMinor: product.unitCostMinor,
      currency: product.currency,
      reorderPoint: product.reorderPoint,
      reorderQuantity: product.reorderQuantity,
      trackLots: product.trackLots,
      trackExpiry: product.trackExpiry,
      active: product.active
    });
    setFormError("");
    setFormOpen(true);
  }

  function closeForm() {
    if (!saving) {
      setFormOpen(false);
      setEditing(null);
      setScannerOpen(false);
      setFormError("");
    }
  }

  function updateDraft<K extends keyof ProductDraft>(key: K, value: ProductDraft[K]) {
    setDraft((current) => ({ ...current, [key]: value }));
  }

  async function submit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setFormError("");
    const normalized: ProductDraft = {
      ...draft,
      sku: draft.sku.trim(),
      name: draft.name.trim(),
      description: draft.description?.trim() ?? "",
      categoryId: draft.categoryId?.trim() ?? "",
      supplierId: draft.supplierId?.trim() ?? "",
      barcode: draft.barcode?.trim() ?? "",
      unit: draft.unit.trim(),
      currency: draft.currency.trim().toUpperCase()
    };

    if (normalized.sku.length < 2 || normalized.sku.length > 64) {
      setFormError("SKU must be 2-64 characters.");
      return;
    }
    if (normalized.name.length < 2 || normalized.name.length > 200) {
      setFormError("Product name must be 2-200 characters.");
      return;
    }
    if (!normalized.unit || normalized.unit.length > 32) {
      setFormError("Unit is required and must be at most 32 characters.");
      return;
    }
    if (!/^[A-Z]{3}$/.test(normalized.currency)) {
      setFormError("Currency must be a three-letter code such as INR.");
      return;
    }
    if (!Number.isSafeInteger(normalized.unitCostMinor) || !Number.isSafeInteger(normalized.reorderPoint) || !Number.isSafeInteger(normalized.reorderQuantity)) {
      setFormError("Cost and reorder values must be whole numbers within the safe numeric range.");
      return;
    }
    if (normalized.unitCostMinor < 0 || normalized.reorderPoint < 0 || normalized.reorderQuantity < 0) {
      setFormError("Cost and reorder values cannot be negative.");
      return;
    }
    if (normalized.trackExpiry && !normalized.trackLots) {
      setFormError("Expiry tracking requires lot tracking.");
      return;
    }

    setSaving(true);
    try {
      if (editing) {
        await stockpilotAPI.updateProduct(editing.id, normalized);
      } else {
        await stockpilotAPI.createProduct(normalized);
      }
      setFormOpen(false);
      setEditing(null);
      await load();
    } catch (error) {
      if (error instanceof APIError && error.status === 401) {
        onSessionExpired();
        return;
      }
      setFormError(error instanceof Error ? error.message : "Could not save the product.");
    } finally {
      setSaving(false);
    }
  }

  return (
    <main className="catalog-page" aria-labelledby="products-title">
      <header className="catalog-header">
        <div>
          <button className="secondary-button" type="button" onClick={onBack}>← Overview</button>
          <p className="eyebrow">Catalog management</p>
          <h1 id="products-title">Products</h1>
          <p className="muted">Create, update, and search the product master without bypassing server-side validation or permissions.</p>
        </div>
        {canWrite && <button className="primary-button" type="button" onClick={startCreate}>Add product</button>}
      </header>

      <section className="catalog-toolbar" aria-label="Catalog filters">
        <label className="search-field">
          <span>Search products</span>
          <input value={query} onChange={(event) => setQuery(event.target.value)} placeholder="Search by SKU or name" maxLength={100} />
        </label>
        <span className="count-pill">{products.length} shown</span>
      </section>

      <div className="status-region" aria-live="polite">
        {state.kind === "loading" && <p className="notice">Loading catalog…</p>}
        {state.kind === "error" && (
          <div className="notice error" role="alert">
            <span>{state.message}</span>
            <button type="button" onClick={() => void load()}>Try again</button>
          </div>
        )}
      </div>

      <section className="panel" aria-labelledby="product-table-title">
        <div className="panel-heading">
          <div><p className="eyebrow">Product master</p><h2 id="product-table-title">Catalog items</h2></div>
        </div>
        {products.length === 0 && state.kind !== "loading" ? (
          <div className="empty-state"><strong>No products found</strong><p>{query ? "Try a different search term." : "Create your first product to start tracking inventory."}</p></div>
        ) : (
          <div className="table-wrap">
            <table>
              <thead>
                <tr><th>Product</th><th>SKU</th><th>Category</th><th>Supplier</th><th>Unit cost</th><th>Reorder</th><th>Status</th>{canWrite && <th>Action</th>}</tr>
              </thead>
              <tbody>
                {products.map((product) => (
                  <tr key={product.id}>
                    <td><strong>{product.name}</strong><small className="table-secondary">{product.barcode ? `Barcode ${product.barcode}` : "No barcode"}</small></td>
                    <td className="mono">{product.sku}</td>
                    <td>{product.categoryId ? categoryNames.get(product.categoryId) ?? "Unknown" : "—"}</td>
                    <td>{product.supplierId ? supplierNames.get(product.supplierId) ?? "Unknown" : "—"}</td>
                    <td>{formatMoney(product.unitCostMinor, product.currency)}</td>
                    <td>{product.reorderPoint.toLocaleString()} / {product.reorderQuantity.toLocaleString()}</td>
                    <td><span className={`status-badge ${product.active ? "good" : ""}`}>{product.active ? "Active" : "Inactive"}</span></td>
                    {canWrite && <td><button className="secondary-button compact-button" type="button" onClick={() => startEdit(product)}>Edit</button></td>}
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </section>

      {canWrite && formOpen && (
        <div className="modal-backdrop" role="presentation" onMouseDown={(event) => { if (event.target === event.currentTarget) closeForm(); }}>
          <section className="modal-card" role="dialog" aria-modal="true" aria-labelledby="product-form-title">
            <div className="panel-heading">
              <div><p className="eyebrow">{editing ? "Edit product" : "New product"}</p><h2 id="product-form-title">{editing ? editing.name : "Add a product"}</h2></div>
              <button className="secondary-button" type="button" onClick={closeForm} disabled={saving}>Close</button>
            </div>

            <form className="catalog-form" onSubmit={submit}>
              <div className="form-grid">
                <label><span>SKU</span><input value={draft.sku} onChange={(event) => updateDraft("sku", event.target.value)} required maxLength={64} /></label>
                <label><span>Product name</span><input value={draft.name} onChange={(event) => updateDraft("name", event.target.value)} required maxLength={200} /></label>
                <label><span>Unit</span><input value={draft.unit} onChange={(event) => updateDraft("unit", event.target.value)} required maxLength={32} /></label>
                <label className="scanner-field"><span>Barcode</span><div className="scanner-input-row"><input value={draft.barcode ?? ""} onChange={(event) => updateDraft("barcode", event.target.value)} maxLength={128} /><button className="secondary-button compact-button" type="button" onClick={() => setScannerOpen(true)} disabled={saving}>Scan</button></div></label>
                <label><span>Unit cost (minor units)</span><input type="number" min="0" step="1" value={draft.unitCostMinor} onChange={(event) => updateDraft("unitCostMinor", Number(event.target.value))} required /></label>
                <label><span>Currency</span><input value={draft.currency} onChange={(event) => updateDraft("currency", event.target.value)} maxLength={3} required /></label>
                <label><span>Reorder point</span><input type="number" min="0" step="1" value={draft.reorderPoint} onChange={(event) => updateDraft("reorderPoint", Number(event.target.value))} required /></label>
                <label><span>Reorder quantity</span><input type="number" min="0" step="1" value={draft.reorderQuantity} onChange={(event) => updateDraft("reorderQuantity", Number(event.target.value))} required /></label>
                <label><span>Category</span><select value={draft.categoryId ?? ""} onChange={(event) => updateDraft("categoryId", event.target.value)}><option value="">No category</option>{categories.map((category) => <option value={category.id} key={category.id}>{category.name}</option>)}</select></label>
                <label><span>Supplier</span><select value={draft.supplierId ?? ""} onChange={(event) => updateDraft("supplierId", event.target.value)}><option value="">No supplier</option>{suppliers.map((supplier) => <option value={supplier.id} key={supplier.id}>{supplier.name} ({supplier.code})</option>)}</select></label>
              </div>

              <label><span>Description</span><textarea value={draft.description ?? ""} onChange={(event) => updateDraft("description", event.target.value)} rows={3} maxLength={2000} /></label>

              <div className="check-grid">
                <label className="check-control"><input type="checkbox" checked={draft.trackLots} onChange={(event) => updateDraft("trackLots", event.target.checked)} /><span>Track lots</span></label>
                <label className="check-control"><input type="checkbox" checked={draft.trackExpiry} onChange={(event) => updateDraft("trackExpiry", event.target.checked)} /><span>Track expiry</span></label>
                <label className="check-control"><input type="checkbox" checked={draft.active} onChange={(event) => updateDraft("active", event.target.checked)} /><span>Active product</span></label>
              </div>

              {formError && <p className="form-error" role="alert">{formError}</p>}
              <div className="form-actions">
                <button className="secondary-button" type="button" onClick={closeForm} disabled={saving}>Cancel</button>
                <button className="primary-button" type="submit" disabled={saving}>{saving ? "Saving…" : editing ? "Save changes" : "Create product"}</button>
              </div>
            </form>
          </section>
          {scannerOpen && <BarcodeScanner onDetected={(value) => { updateDraft("barcode", value); setScannerOpen(false); }} onClose={() => setScannerOpen(false)} />}
        </div>
      )}
    </main>
  );
}

function formatMoney(minor: number, currency: string) {
  try {
    return new Intl.NumberFormat("en-IN", { style: "currency", currency: currency || "INR", maximumFractionDigits: 2 }).format(minor / 100);
  } catch {
    return `${currency || "INR"} ${(minor / 100).toFixed(2)}`;
  }
}
