export type ReplenishmentRisk = "out_of_stock" | "critical" | "reorder" | "watch" | "healthy";

export type ReplenishmentReadinessItem = {
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
  outboundUnits: number;
  averageDailyOutbound: number;
  daysOfCover?: number;
  risk: ReplenishmentRisk;
};

export type ReplenishmentReadinessReport = {
  asOf: string;
  windowDays: number;
  items: ReplenishmentReadinessItem[];
  nextCursor?: string;
};

export async function replenishmentReadiness(days = 30, limit = 500, cursor?: string): Promise<ReplenishmentReadinessReport> {
  const safeDays = Math.min(Math.max(days, 1), 365);
  const safeLimit = Math.min(Math.max(limit, 1), 5000);
  const params = new URLSearchParams({ days: String(safeDays), limit: String(safeLimit) });
  if (cursor) {
    params.set("cursor", cursor);
  }
  const response = await fetch(`/api/v1/reports/replenishment-readiness?${params.toString()}`, { credentials: "include", headers: { Accept: "application/json" } });
  if (!response.ok) {
    throw new Error(`Replenishment readiness request failed (${response.status}).`);
  }
  return await response.json() as ReplenishmentReadinessReport;
}

export function replenishmentReadinessCSV(days = 30, limit = 5000) {
  const safeDays = Math.min(Math.max(days, 1), 365);
  const safeLimit = Math.min(Math.max(limit, 1), 5000);
  const params = new URLSearchParams({ format: "csv", days: String(safeDays), limit: String(safeLimit) });
  window.location.assign(`/api/v1/reports/replenishment-readiness?${params.toString()}`);
}
