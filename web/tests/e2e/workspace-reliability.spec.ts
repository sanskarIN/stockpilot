import { test, expect } from "@playwright/test";

async function mockAuthenticatedWorkspace(page: Parameters<Parameters<typeof test>[1]>[0]["page"]) {
  await page.route("**/api/v1/auth/me", async (route) => {
    await route.fulfill({
      status: 200,
      contentType: "application/json",
      body: JSON.stringify({
        id: "e2e-user",
        email: "e2e@example.test",
        displayName: "E2E Operator",
        role: "admin",
        active: true,
        createdAt: "2026-01-01T00:00:00Z",
        updatedAt: "2026-01-01T00:00:00Z",
      }),
    });
  });
  await page.route("**/api/v1/reports/inventory-valuation**", async (route) => {
    await route.fulfill({ status: 200, contentType: "application/json", body: JSON.stringify({ items: [], totals: [] }) });
  });
  await page.route("**/api/v1/products**", async (route) => {
    await route.fulfill({ status: 200, contentType: "application/json", body: JSON.stringify({ items: [] }) });
  });
  await page.route("**/api/v1/inventory/reorder-suggestions**", async (route) => {
    await route.fulfill({ status: 200, contentType: "application/json", body: JSON.stringify({ items: [] }) });
  });
  await page.route("**/api/v1/orders**", async (route) => {
    await route.fulfill({ status: 200, contentType: "application/json", body: JSON.stringify({ items: [] }) });
  });
}

test.describe("workspace reliability", () => {
  test("renders the signed-out login contract", async ({ page }) => {
    await page.route("**/api/v1/auth/me", async (route) => {
      await route.fulfill({ status: 401, contentType: "application/json", body: JSON.stringify({ error: "unauthorized" }) });
    });

    await page.goto("/");

    await expect(page.getByRole("heading", { name: "Sign in to StockPilot" })).toBeVisible();
    await expect(page.getByLabel("Email")).toHaveAttribute("autocomplete", "username");
    await expect(page.getByLabel("Password")).toHaveAttribute("autocomplete", "current-password");
    await expect(page.getByRole("button", { name: "Sign in" })).toBeDisabled();
  });

  test("loads the authenticated dashboard shell with API fixtures", async ({ page }) => {
    await mockAuthenticatedWorkspace(page);

    await page.goto("/");

    await expect(page.getByRole("heading", { name: "Inventory overview" })).toBeVisible();
    await expect(page.getByText("E2E Operator")).toBeVisible();
    await expect(page.getByRole("button", { name: "Sign out" })).toBeVisible();
  });
});
