import { normalizeServerUrl, originPattern } from "./url.js";

const form = document.querySelector("#server-form");
const serverInput = document.querySelector("#server-url");
const statusDot = document.querySelector("#status-dot");
const statusTitle = document.querySelector("#status-title");
const statusDetail = document.querySelector("#status-detail");
const refreshButton = document.querySelector("#refresh-button");
const openButton = document.querySelector("#open-button");
const errorBox = document.querySelector("#error");

let serverUrl = "";

form.addEventListener("submit", async (event) => {
  event.preventDefault();
  clearError();

  try {
    const normalized = normalizeServerUrl(serverInput.value);
    const pattern = originPattern(normalized);
    const alreadyGranted = await chrome.permissions.contains({ origins: [pattern] });
    const granted = alreadyGranted || await chrome.permissions.request({ origins: [pattern] });
    if (!granted) {
      throw new Error("Server access permission was not granted.");
    }

    serverUrl = normalized;
    serverInput.value = normalized;
    await chrome.storage.local.set({ serverUrl });
    openButton.disabled = false;
    await checkStatus();
  } catch (error) {
    showError(error);
  }
});

refreshButton.addEventListener("click", async () => {
  clearError();
  try {
    if (!serverUrl) {
      throw new Error("Save a StockPilot server URL first.");
    }
    await checkStatus();
  } catch (error) {
    showError(error);
  }
});

openButton.addEventListener("click", async () => {
  if (serverUrl) {
    await chrome.tabs.create({ url: serverUrl });
  }
});

async function initialize() {
  const stored = await chrome.storage.local.get(["serverUrl"]);
  if (!stored.serverUrl) {
    return;
  }

  try {
    serverUrl = normalizeServerUrl(stored.serverUrl);
    serverInput.value = serverUrl;
    openButton.disabled = false;
    const permission = await chrome.permissions.contains({ origins: [originPattern(serverUrl)] });
    if (permission) {
      await checkStatus();
    } else {
      setStatus("idle", "Permission needed", "Save & connect to grant access to this server.");
    }
  } catch (error) {
    showError(error);
  }
}

async function checkStatus() {
  setBusy(true);
  setStatus("idle", "Checking…", "Contacting the configured StockPilot server.");
  try {
    const controller = new AbortController();
    const timeout = setTimeout(() => controller.abort(), 8000);
    try {
      const [metaResponse, readyResponse] = await Promise.all([
        fetch(new URL("/api/v1/meta", serverUrl), {
          headers: { Accept: "application/json" },
          cache: "no-store",
          credentials: "omit",
          signal: controller.signal,
        }),
        fetch(new URL("/readyz", serverUrl), {
          headers: { Accept: "application/json" },
          cache: "no-store",
          credentials: "omit",
          signal: controller.signal,
        }),
      ]);

      if (!metaResponse.ok) {
        throw new Error(`Server metadata returned HTTP ${metaResponse.status}.`);
      }
      const meta = await metaResponse.json();
      const readiness = await readyResponse.json().catch(() => ({}));
      if (readyResponse.ok) {
        setStatus(
          "online",
          `${meta.name || "StockPilot"} is ready`,
          `${meta.version || "Unknown version"} · ${readiness.status || "ready"}`,
        );
      } else {
        setStatus(
          "offline",
          "Server reachable, not ready",
          readiness.status || `Readiness returned HTTP ${readyResponse.status}.`,
        );
      }
    } finally {
      clearTimeout(timeout);
    }
  } catch (error) {
    const message = error?.name === "AbortError"
      ? "The server did not respond within 8 seconds."
      : (error?.message || "Unable to reach the StockPilot server.");
    setStatus("offline", "Connection failed", message);
    throw new Error(message);
  } finally {
    setBusy(false);
  }
}

function setBusy(busy) {
  serverInput.disabled = busy;
  refreshButton.disabled = busy;
  openButton.disabled = busy || !serverUrl;
  form.querySelector("button[type='submit']").disabled = busy;
}

function setStatus(kind, title, detail) {
  statusDot.className = `status-dot${kind === "idle" ? "" : ` ${kind}`}`;
  statusTitle.textContent = title;
  statusDetail.textContent = detail;
}

function showError(error) {
  const message = error?.message || String(error);
  errorBox.textContent = message;
  errorBox.hidden = false;
}

function clearError() {
  errorBox.textContent = "";
  errorBox.hidden = true;
}

initialize();
