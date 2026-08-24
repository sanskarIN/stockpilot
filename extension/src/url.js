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
  return `${normalizeServerUrl(value)}/*`;
}
