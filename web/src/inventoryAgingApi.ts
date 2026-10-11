import type { InventoryAgingReport } from "./types";

export async function inventoryAging(limit = 500): Promise<InventoryAgingReport> {
  const bounded = Math.min(Math.max(limit, 1), 5000);
  const response = await fetch(`/api/v1/reports/inventory-aging?limit=${bounded}`, { credentials: "include", headers: { Accept: "application/json" } });
  if (!response.ok) throw new Error(`Inventory aging request failed with status ${response.status}.`);
  return await response.json() as InventoryAgingReport;
}
