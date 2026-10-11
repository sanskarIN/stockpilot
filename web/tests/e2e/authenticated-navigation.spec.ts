import { expect, test, type Page } from "@playwright/test";

async function mockAuthenticatedWorkspace(page: Page) {
  await page.route("**/api/v1/**", async (route) => {
    const pathname = new URL(route.request().url()).pathname;

    if (pathname === "/api/v1/auth/me") {
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
      return;
    }

    if (pathname === "/api/v1/reports/inventory-valuation") {
      await route.fulfill({
        status: 200,
        contentType: "application/json",
        body: JSON.stringify({ items: [], totals: [] }),
      });
      return;
    }

    await route.fulfill({
      status: 200,
      contentType: "application/json",
      body: JSON.stringify({ items: [] }),
    });
  });
}

test("authenticated users can enter every primary operational workspace", async ({ page }) => {
  await mockAuthenticatedWorkspace(page);
  await page.goto("/");

  await expect(page.getByRole("heading", { name: "Inventory overview", exact: true })).toBeVisible();

  const workspaces = [
    { nav: "Products", heading: "Products" },
    { nav: "Inventory", heading: "Move stock safely" },
    { nav: "Purchase orders", heading: "Purchase orders" },
    { nav: "Warehouses", heading: "Warehouses & locations" },
    { nav: "Lot inventory", heading: "Lot inventory" },
    { nav: "Audit", heading: "Audit history" },
  ];

  for (const workspace of workspaces) {
    await page.getByRole("button", { name: workspace.nav, exact: true }).click();
    await expect(page.getByRole("heading", { name: workspace.heading, exact: true })).toBeVisible();
    await page.getByRole("button", { name: /Overview/ }).click();
    await expect(page.getByRole("heading", { name: "Inventory overview", exact: true })).toBeVisible();
  }
});
