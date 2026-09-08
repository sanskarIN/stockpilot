import { expect, test, type Page } from "@playwright/test";

const user = {
  id: "e2e-user",
  email: "e2e@example.test",
  displayName: "E2E Operator",
  role: "admin",
  active: true,
  createdAt: "2026-01-01T00:00:00Z",
  updatedAt: "2026-01-01T00:00:00Z",
};

const category = { id: "cat-1", name: "General", description: "", createdAt: "2026-01-01T00:00:00Z", updatedAt: "2026-01-01T00:00:00Z" };
const supplier = { id: "sup-1", code: "SUP-001", name: "Synthetic Supplier", email: "supplier@example.test", phone: "", active: true };
const warehouse = { id: "wh-1", code: "WH-001", name: "Synthetic Warehouse", timezone: "Asia/Kolkata", active: true };
const location = { id: "loc-1", warehouseId: "wh-1", code: "BIN-001", name: "Primary Bin", active: true };
const product = {
  id: "product-1",
  sku: "E2E-001",
  name: "Synthetic Product",
  description: "",
  categoryId: "cat-1",
  supplierId: "sup-1",
  barcode: "890000000001",
  unit: "pcs",
  unitCostMinor: 12500,
  currency: "INR",
  reorderPoint: 5,
  reorderQuantity: 10,
  trackLots: false,
  trackExpiry: false,
  active: true,
  createdAt: "2026-01-01T00:00:00Z",
  updatedAt: "2026-01-01T00:00:00Z",
};

async function mockAuthenticatedData(page: Page) {
  await page.route("**/api/v1/auth/me", async (route) => {
    await route.fulfill({ status: 200, contentType: "application/json", body: JSON.stringify(user) });
  });
  await page.route("**/api/v1/reports/inventory-valuation**", async (route) => {
    await route.fulfill({ status: 200, contentType: "application/json", body: JSON.stringify({ items: [], totals: [] }) });
  });
  await page.route("**/api/v1/inventory/reorder-suggestions**", async (route) => {
    await route.fulfill({ status: 200, contentType: "application/json", body: JSON.stringify({ items: [] }) });
  });
  await page.route("**/api/v1/orders**", async (route) => {
    if (route.request().method() === "GET") {
      await route.fulfill({ status: 200, contentType: "application/json", body: JSON.stringify({ items: [] }) });
      return;
    }
    await route.fallback();
  });
}

async function mockCatalog(page: Page) {
  let products = [] as typeof product[];
  await page.route("**/api/v1/categories**", async (route) => {
    await route.fulfill({ status: 200, contentType: "application/json", body: JSON.stringify({ items: [category] }) });
  });
  await page.route("**/api/v1/suppliers**", async (route) => {
    await route.fulfill({ status: 200, contentType: "application/json", body: JSON.stringify({ items: [supplier] }) });
  });
  await page.route("**/api/v1/products**", async (route) => {
    if (route.request().method() === "GET") {
      await route.fulfill({ status: 200, contentType: "application/json", body: JSON.stringify({ items: products }) });
      return;
    }
    if (route.request().method() === "POST") {
      products = [product];
      await route.fulfill({ status: 201, contentType: "application/json", body: JSON.stringify(product) });
      return;
    }
    await route.fallback();
  });
}

async function mockInventory(page: Page) {
  await page.route("**/api/v1/products**", async (route) => {
    await route.fulfill({ status: 200, contentType: "application/json", body: JSON.stringify({ items: [product] }) });
  });
  await page.route("**/api/v1/warehouses**", async (route) => {
    await route.fulfill({ status: 200, contentType: "application/json", body: JSON.stringify({ items: [warehouse] }) });
  });
  await page.route("**/api/v1/locations**", async (route) => {
    await route.fulfill({ status: 200, contentType: "application/json", body: JSON.stringify({ items: [location] }) });
  });
  await page.route("**/api/v1/inventory/movements", async (route) => {
    expect(route.request().method()).toBe("POST");
    const body = route.request().postDataJSON();
    expect(body).toMatchObject({ productId: "product-1", locationId: "loc-1", type: "stock_in", quantityDelta: 7 });
    await route.fulfill({
      status: 200,
      contentType: "application/json",
      body: JSON.stringify({ movementId: "movement-1", balance: { productId: "product-1", locationId: "loc-1", quantity: 7 } }),
    });
  });
}

test.describe("authenticated mutation workflows", () => {
  test("creates a catalog product through the authenticated UI", async ({ page }) => {
    await mockAuthenticatedData(page);
    await mockCatalog(page);

    await page.goto("/");
    await expect(page.getByRole("heading", { name: "Inventory overview" })).toBeVisible();
    await page.getByRole("button", { name: "Products", exact: true }).click();
    await expect(page.getByRole("heading", { name: "Products" })).toBeVisible();
    await page.getByRole("button", { name: "Add product" }).click();

    await page.getByLabel("SKU").fill("E2E-001");
    await page.getByLabel("Product name").fill("Synthetic Product");
    await page.getByRole("textbox", { name: "Unit" }).fill("pcs");
    await page.getByLabel("Unit cost (minor units)").fill("12500");
    await page.getByLabel("Currency").fill("INR");
    await page.getByLabel("Reorder point").fill("5");
    await page.getByLabel("Reorder quantity").fill("10");
    await page.getByRole("button", { name: "Create product" }).click();

    await expect(page.getByText("Synthetic Product").first()).toBeVisible();
    await expect(page.getByText("E2E-001").first()).toBeVisible();
  });

  test("records a stock-in mutation through the authenticated inventory UI", async ({ page }) => {
    await mockAuthenticatedData(page);
    await mockInventory(page);

    await page.goto("/");
    await page.getByRole("button", { name: "Inventory", exact: true }).click();
    await expect(page.getByRole("heading", { name: "Move stock safely" })).toBeVisible();

    await page.getByLabel("Quantity").fill("7");
    await page.getByRole("button", { name: "Stock in" }).last().click();

    await expect(page.getByRole("status")).toContainText("Movement movement-1 recorded");
    await expect(page.getByRole("status")).toContainText("New balance: 7");
  });
});
