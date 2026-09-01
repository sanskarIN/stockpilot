import { normalizeServerUrl, originPattern } from "./url.js";
import { startScanner, stopScanner } from "./scanner.js";

const form = document.querySelector("#server-form");
const serverInput = document.querySelector("#server-url");
const statusDot = document.querySelector("#status-dot");
const statusTitle = document.querySelector("#status-title");
const statusDetail = document.querySelector("#status-detail");
const refreshButton = document.querySelector("#refresh-button");
const scanButton = document.querySelector("#scan-button");
const openButton = document.querySelector("#open-button");
const scanner = document.querySelector("#scanner");
const scannerVideo = document.querySelector("#scanner-video");
const scannerStatus = document.querySelector("#scanner-status");
const closeScanner = document.querySelector("#close-scanner");
const manualForm = document.querySelector("#manual-form");
const manualCode = document.querySelector("#manual-code");
const resultCard = document.querySelector("#scan-result");
const resultTitle = document.querySelector("#result-title");
const resultDetail = document.querySelector("#result-detail");
const openResult = document.querySelector("#open-result");
const errorBox = document.querySelector("#error");

let serverUrl = "";
let scannedValue = "";

form.addEventListener("submit", async (event) => {
  event.preventDefault();
  clearError();
  try {
    const normalized = normalizeServerUrl(serverInput.value);
    const granted = await chrome.permissions.request({ origins: [originPattern(normalized)] });
    if (!granted) throw new Error("Server access permission was not granted.");
    serverUrl = normalized;
    serverInput.value = normalized;
    await chrome.storage.local.set({ serverUrl });
    openButton.disabled = false;
    await checkStatus();
  } catch (error) { showError(error); }
});

refreshButton.addEventListener("click", async () => {
  clearError();
  try {
    if (!serverUrl) throw new Error("Save a StockPilot server URL first.");
    await checkStatus();
  } catch (error) { showError(error); }
});

scanButton.addEventListener("click", async () => {
  clearError();
  if (!serverUrl) { showError(new Error("Save a StockPilot server URL first.")); return; }
  resultCard.hidden = true;
  scanner.hidden = false;
  scanButton.disabled = true;
  try {
    await startScanner(scannerVideo, handleScan, (status) => { scannerStatus.textContent = status; });
  } catch (error) {
    scannerStatus.textContent = error?.message || "Scanner unavailable.";
    manualCode.focus();
  }
});

closeScanner.addEventListener("click", closeScanView);

manualForm.addEventListener("submit", (event) => {
  event.preventDefault();
  const value = manualCode.value.trim();
  if (value) handleScan(value);
});

openButton.addEventListener("click", () => openInStockPilot());
openResult.addEventListener("click", () => openInStockPilot(scannedValue));

function closeScanView() {
  stopScanner();
  scanner.hidden = true;
  scanButton.disabled = false;
}

function handleScan(value) {
  scannedValue = String(value).trim();
  if (!scannedValue) return;
  closeScanView();
  resultTitle.textContent = "Code captured";
  resultDetail.textContent = `${scannedValue} — open StockPilot to resolve the product using your existing signed-in session.`;
  resultCard.hidden = false;
}

async function openInStockPilot(barcode = scannedValue) {
  if (!serverUrl) return;
  const url = new URL("/", serverUrl);
  url.searchParams.set("barcode", barcode);
  await chrome.tabs.create({ url: url.toString() });
}

async function initialize() {
  const stored = await chrome.storage.local.get(["serverUrl"]);
  if (!stored.serverUrl) return;
  serverUrl = normalizeServerUrl(stored.serverUrl);
  serverInput.value = serverUrl;
  openButton.disabled = false;
  const permission = await chrome.permissions.contains({ origins: [originPattern(serverUrl)] });
  if (permission) await checkStatus();
  else setStatus("idle", "Permission needed", "Save & connect to grant access to this server.");
}

async function checkStatus() {
  setBusy(true); setStatus("idle", "Checking…", "Contacting the configured StockPilot server.");
  try {
    const controller = new AbortController();
    const timeout = setTimeout(() => controller.abort(), 8000);
    try {
      const [metaResponse, readyResponse] = await Promise.all([
        fetch(new URL("/api/v1/meta", serverUrl), { headers: { Accept: "application/json" }, cache: "no-store", credentials: "omit", signal: controller.signal }),
        fetch(new URL("/readyz", serverUrl), { headers: { Accept: "application/json" }, cache: "no-store", credentials: "omit", signal: controller.signal }),
      ]);
      if (!metaResponse.ok) throw new Error(`Server metadata returned HTTP ${metaResponse.status}.`);
      const meta = await metaResponse.json();
      const readiness = await readyResponse.json().catch(() => ({}));
      if (readyResponse.ok) setStatus("online", `${meta.name || "StockPilot"} is ready`, `${meta.version || "Unknown version"} · ${readiness.status || "ready"}`);
      else setStatus("offline", "Server reachable, not ready", readiness.status || `Readiness returned HTTP ${readyResponse.status}.`);
    } finally { clearTimeout(timeout); }
  } catch (error) {
    const message = error?.name === "AbortError" ? "The server did not respond within 8 seconds." : (error?.message || "Unable to reach the StockPilot server.");
    setStatus("offline", "Connection failed", message); throw new Error(message);
  } finally { setBusy(false); }
}

function setBusy(busy) {
  serverInput.disabled = busy; refreshButton.disabled = busy; scanButton.disabled = busy || !serverUrl; openButton.disabled = busy || !serverUrl; form.querySelector("button[type='submit']").disabled = busy;
}
function setStatus(kind, title, detail) { statusDot.className = `status-dot${kind === "idle" ? "" : ` ${kind}`}`; statusTitle.textContent = title; statusDetail.textContent = detail; }
function showError(error) { errorBox.textContent = error?.message || String(error); errorBox.hidden = false; }
function clearError() { errorBox.textContent = ""; errorBox.hidden = true; }
initialize().catch(showError);
