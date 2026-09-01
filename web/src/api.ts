import type {
  Category,
  InventoryValuationReport,
  ListResponse,
  Product,
  PurchaseOrder,
  ReorderSuggestion,
  Supplier,
  User
} from "./types";

export class APIError extends Error {
  readonly status: number;

  constructor(status: number, message: string) {
    super(message);
    this.name = "APIError";
    this.status = status;
  }
}

type ErrorBody = { error?: string };

async function request<T>(path: string, init: RequestInit = {}): Promise<T> {
  const method = (init.method ?? "GET").toUpperCase();
  const headers = new Headers(init.headers);
  headers.set("Accept", "application/json");
  if (init.body && !headers.has("Content-Type")) {
    headers.set("Content-Type", "application/json");
  }
  if (!["GET", "HEAD", "OPTIONS"].includes(method) && path !== "/api/v1/auth/login") {
    headers.set("X-StockPilot-CSRF", "1");
  }

  const response = await fetch(path, {
    ...init,
    method,
    headers,
    credentials: "include"
  });
  if (!response.ok) {
    let message = `Request failed with status ${response.status}.`;
    try {
      const body = (await response.json()) as ErrorBody;
      if (typeof body.error === "string" && body.error.trim()) message = body.error;
    } catch {
      // Preserve the safe fallback when the response is not JSON.
    }
    throw new APIError(response.status, message);
  }
  if (response.status === 204) return undefined as T;
  return (await response.json()) as T;
}

const emptyValuation: InventoryValuationReport = { items: [], totals: [] };

export const stockpilotAPI = {
  me: () => request<User>("/api/v1/auth/me"),
  login: (email: string, password: string) =>
    request<{ user: User; expiresAt: string }>("/api/v1/auth/login", {
      method: "POST",
      body: JSON.stringify({ email, password })
    }),
  logout: () => request<void>("/api/v1/auth/logout", { method: "POST" }),
  listCategories: () => request<ListResponse<Category>>("/api/v1/categories"),
  listSuppliers: () => request<ListResponse<Supplier>>("/api/v1/suppliers?active=true"),
  listProducts: (query = "") => request<ListResponse<Product>>(`/api/v1/products?active=false&limit=200&q=${encodeURIComponent(query.trim())}`),
  createProduct: (product: Omit<Product, "id" | "createdAt" | "updatedAt">) =>
    request<Product>("/api/v1/products", { method: "POST", body: JSON.stringify(product) }),
  updateProduct: (id: string, product: Omit<Product, "id" | "createdAt" | "updatedAt">) =>
    request<Product>(`/api/v1/products/${encodeURIComponent(id)}`, { method: "PUT", body: JSON.stringify(product) }),
  productByBarcode: (barcode: string) => request<Product>(`/api/v1/products/by-barcode/${encodeURIComponent(barcode.trim())}`),
  dashboard: async (includeValuation = true) => {
    const valuationRequest = includeValuation
      ? request<InventoryValuationReport>("/api/v1/reports/inventory-valuation?limit=100")
      : Promise.resolve(emptyValuation);
    const [products, reorderSuggestions, valuation, orders] = await Promise.all([
      request<ListResponse<Product>>("/api/v1/products?active=true&limit=100"),
      request<ListResponse<ReorderSuggestion>>("/api/v1/inventory/reorder-suggestions?limit=100"),
      valuationRequest,
      request<ListResponse<PurchaseOrder>>("/api/v1/orders?limit=100")
    ]);
    return {
      products: products.items ?? [],
      reorderSuggestions: reorderSuggestions.items ?? [],
      valuation: {
        items: valuation.items ?? [],
        totals: valuation.totals ?? []
      },
      orders: orders.items ?? []
    };
  }
};
