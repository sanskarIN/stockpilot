const EXTENSION_SCHEMA_VERSION = 1;

chrome.runtime.onInstalled.addListener(async () => {
  const current = await chrome.storage.local.get(["extensionSchemaVersion"]);
  if (current.extensionSchemaVersion !== EXTENSION_SCHEMA_VERSION) {
    await chrome.storage.local.set({ extensionSchemaVersion: EXTENSION_SCHEMA_VERSION });
  }
});
