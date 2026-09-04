import { stockpilotAPI } from "./api";
import type { StockMovementHistoryReport } from "./types";

export function stockMovementHistory(days = 30, limit = 500) {
  const safeDays = Math.min(Math.max(days, 1), 365);
  const safeLimit = Math.min(Math.max(limit, 1), 5000);
  return stockpilotAPI.requestStockMovementHistory(safeDays, safeLimit);
}

export function stockMovementHistoryCSV(days = 30, limit = 5000) {
  const safeDays = Math.min(Math.max(days, 1), 365);
  const safeLimit = Math.min(Math.max(limit, 1), 5000);
  window.location.assign(`/api/v1/reports/stock-movement-history?format=csv&days=${safeDays}&limit=${safeLimit}`);
}

export type { StockMovementHistoryReport };
