export type Role = "admin" | "manager" | "operator" | "viewer";

export type User = {
  id: string;
  email: string;
  displayName: string;
  role: Role;
  active: boolean;
  createdAt: string;
  updatedAt: string;
  lastLoginAt?: string;
};

export type Category = {
  id: string;
  name: string;
  description: string;
  createdAt: string;
  updatedAt: string;
};

export type Supplier = {
  id: string;
  code: string;
  name: string;
  email?: string;
  phone?: string;
  notes: string;
  active: boolean;
  createdAt: string;
  updatedAt: string;
};

export type Product = {
  id: string;
  sku: string;
  name: string;
  description?: string;
  categoryId?: string;
  supplierId?: string;
  barcode?: string;
  unit: string;
  unitCostMinor: number;
  currency: string;
  reorderPoint: number;
  reorderQuantity: number;
  trackLots: boolean;
  trackExpiry: boolean;
  active: boolean;
  createdAt: string;
  updatedAt: string;
};

export type StockBalance = {
  productId: string;
  locationId: string;
  lotId?: string;
  quantity: number;
  updatedAt: string;
};

export type ReorderSuggestion = {
  productId: string;
  sku: string;
  name: string;
  supplierId?: string;
  unit: string;
  onHand: number;
  reorderPoint: number;
  reorderQuantity: number;
  targetStock: number;
  suggestedQuantity: number;
};

export type InventoryValuationItem = {
  productId: string;
  sku: string;
  name: string;
  unit: string;
  onHand: number;
  unitCostMinor: number;
  currency: string;
  valueMinor: number;
};

export type InventoryValuationTotal = {
  currency: string;
  valueMinor: number;
};

export type InventoryValuationReport = {
  items: InventoryValuationItem[];
  totals: InventoryValuationTotal[];
};

export type PurchaseOrder = {
  id: string;
  number: string;
  status: string;
  currency: string;
  createdAt: string;
};

export type ListResponse<T> = { items: T[] };
