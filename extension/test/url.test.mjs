import test from "node:test";
import assert from "node:assert/strict";
import { buildInventoryHandoffUrl, normalizeServerUrl, originPattern } from "../src/url.js";

test("normalizes a StockPilot origin", () => {
  assert.equal(
    normalizeServerUrl("  https://inventory.example.com/  "),
    "https://inventory.example.com",
  );
});

test("preserves a configured server port", () => {
  assert.equal(
    normalizeServerUrl("http://localhost:8080"),
    "http://localhost:8080",
  );
});

test("host permission pattern remains host-scoped", () => {
  assert.equal(
    originPattern("http://localhost:8080"),
    "http://localhost/*",
  );
});

test("rejects paths because APIs are rooted at the deployment origin", () => {
  assert.throws(
    () => normalizeServerUrl("https://inventory.example.com/stockpilot"),
    /without a path/,
  );
});

test("rejects embedded credentials", () => {
  assert.throws(
    () => normalizeServerUrl("https://user:password@inventory.example.com"),
    /credentials/,
  );
});

test("rejects unsupported protocols", () => {
  assert.throws(
    () => normalizeServerUrl("ftp://inventory.example.com"),
    /http:\/\//,
  );
});

test("builds a safe stock-out handoff", () => {
  const url = new URL(buildInventoryHandoffUrl("https://inventory.example.com", "890123", "stock_out"));
  assert.equal(url.origin, "https://inventory.example.com");
  assert.equal(url.pathname, "/");
  assert.equal(url.searchParams.get("barcode"), "890123");
  assert.equal(url.searchParams.get("inventoryOperation"), "stock_out");
});

test("rejects an empty barcode", () => {
  assert.throws(
    () => buildInventoryHandoffUrl("https://inventory.example.com", "", "stock_in"),
    /barcode is required/,
  );
});

test("rejects unsupported inventory operations", () => {
  assert.throws(
    () => buildInventoryHandoffUrl("https://inventory.example.com", "890123", "delete"),
    /Unsupported inventory operation/,
  );
});
