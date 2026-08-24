import test from "node:test";
import assert from "node:assert/strict";
import { normalizeServerUrl, originPattern } from "../src/url.js";

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
