export function normalizeServerUrl(value) {
  const candidate = String(value ?? "").trim();
  if (!candidate) {
    throw new Error("Enter your StockPilot server URL.");
  }

  let url;
  try {
    url = new URL(candidate);
  } catch {
    throw new Error("Enter a valid http:// or https:// server URL.");
  }

  if (url.protocol !== "https:" && url.protocol !== "http:") {
    throw new Error("Server URL must use http:// or https://.");
  }
  if (!url.hostname) {
    throw new Error("Server URL must include a host.");
  }
  if (url.username || url.password) {
    throw new Error("Do not put credentials in the server URL.");
  }
  if ((url.pathname && url.pathname !== "/") || url.search || url.hash) {
    throw new Error("Use the StockPilot server origin without a path, query, or fragment.");
  }

  return url.origin;
}

export function originPattern(value) {
  const url = new URL(normalizeServerUrl(value));
  // Chrome match patterns do not scope host permissions by port, so request
  // only the configured scheme + host instead of a broad all-host permission.
  return `${url.protocol}//${url.hostname}/*`;
}

export function buildInventoryHandoffUrl(serverUrl, barcode, operation = "stock_in") {
  const origin = normalizeServerUrl(serverUrl);
  const code = String(barcode ?? "").trim();
  if (!code) {
    throw new Error("A barcode is required for inventory handoff.");
  }
  const allowedOperations = new Set(["stock_in", "stock_out", "adjustment", "transfer"]);
  if (!allowedOperations.has(operation)) {
    throw new Error("Unsupported inventory operation.");
  }

  const url = new URL("/", origin);
  url.searchParams.set("barcode", code);
  url.searchParams.set("inventoryOperation", operation);
  return url.toString();
}
