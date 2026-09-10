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

function makeOrder(status: "draft" | "ordered" | "received", received = 0) {
  return {
    id: "po-1",
    number: "E2E-PO-001",
    supplierId: supplier.id,
    warehouseId: warehouse.id,
    status,
    currency: "INR",
    notes: "Synthetic E2E purchase order",
    lines: [{ id: "pol-1", purchaseOrderId: "po-1", productId: product.id, quantity: 3, received, unitCostMinor: 12500 }],
    createdAt: "2026-01-01T00:00:00Z",
    updatedAt: "2026-01-01T00:00:00Z",
  };
}

async function mockAuthenticatedData(page: Page) {
  await page.route("**/api/v1/auth/me", async (route) => {
    await route.fulfill({ status: 200, contentType: "application/json", body: JSON.stringify(user) });
  });
  await page.route("**/api/v1/products**", async (route) => {
    if (route.request().method() === "GET") {
      await route.fulfill({ status: 200, contentType: "application/json", body: JSON.stringify({ items: [] }) });
      return;
    }
    await route.fallback();
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

async function mockPurchasing(page: Page) {
  let order = makeOrder("draft");
  await page.route("**/api/v1/suppliers**", async (route) => {
    await route.fulfill({ status: 200, contentType: "application/json", body: JSON.stringify({ items: [supplier] }) });
  });
  await page.route("**/api/v1/warehouses**", async (route) => {
    await route.fulfill({ status: 200, contentType: "application/json", body: JSON.stringify({ items: [warehouse] }) });
  });
  await page.route("**/api/v1/locations**", async (route) => {
    await route.fulfill({ status: 200, contentType: "application/json", body: JSON.stringify({ items: [location] }) });
  });
  await page.route("**/api/v1/products**", async (route) => {
    await route.fulfill({ status: 200, contentType: "application/json", body: JSON.stringify({ items: [product] }) });
  });
  await page.route("**/api/v1/orders", async (route) => {
    const method = route.request().method();
    if (method === "GET") {
      await route.fulfill({ status: 200, contentType: "application/json", body: JSON.stringify({ items: [order] }) });
      return;
    }
    if (method === "POST") {
      const body = route.request().postDataJSON();
      expect(body).toMatchObject({ number: "E2E-PO-001", supplierId: supplier.id, warehouseId: warehouse.id, status: "draft", currency: "INR" });
      expect(body.lines).toHaveLength(1);
      order = makeOrder("draft");
      await route.fulfill({ status: 201, contentType: "application/json", body: JSON.stringify(order) });
      return;
    }
    await route.fallback();
  });
  await page.route("**/api/v1/orders/po-1/status", async (route) => {
    expect(route.request().method()).toBe("PATCH");
    const body = route.request().postDataJSON();
    expect(body).toEqual({ status: "ordered" });
    order = makeOrder("ordered");
    await route.fulfill({ status: 200, contentType: "application/json", body: JSON.stringify(order) });
  });
  await page.route("**/api/v1/orders/po-1/lines/pol-1/receive", async (route) => {
    expect(route.request().method()).toBe("POST");
    const body = route.request().postDataJSON();
    expect(body).toMatchObject({ quantity: 3, locationId: location.id });
    order = makeOrder("received", 3);
    await route.fulfill({ status: 200, contentType: "application/json", body: JSON.stringify(order) });
  });
}

async function mockReports(page: Page) {
  await page.route("**/api/v1/reports/overview", async (route) => {
    await route.fulfill({
      status: 200,
      contentType: "application/json",
      body: JSON.stringify({
        inventory: { productCount: 1, activeProductCount: 1, activeWarehouseCount: 1, activeLocationCount: 1, totalUnits: 12, lowStockBalanceCount: 0, outOfStockCount: 0 },
        purchasing: { totalOrders: 2, draftOrders: 1, orderedOrders: 1, partiallyReceivedOrders: 0, receivedOrders: 0, cancelledOrders: 0, outstandingUnits: 3 },
      }),
    });
  });
  await page.route("**/api/v1/reports/inventory-valuation**", async (route) => {
    await route.fulfill({ status: 200, contentType: "application/json", body: JSON.stringify({ items: [
      { productId: product.id, sku: product.sku, name: product.name, unit: product.unit, onHand: 12, unitCostMinor: 12500, currency: "INR", valueMinor: 150000 },
    ], totals: [{ currency: "INR", valueMinor: 150000 }] }) });
  });
  await page.route("**/api/v1/reports/inventory-aging**", async (route) => {
    await route.fulfill({ status: 200, contentType: "application/json", body: JSON.stringify({ items: [
      { productId: product.id, sku: product.sku, name: product.name, locationId: location.id, quantity: 12, ageDays: 12, bucket: "0-30", asOf: "2026-09-01T00:00:00Z", lastMovementAt: "2026-08-20T00:00:00Z" },
    ] }) });
  });
  await page.route("**/api/v1/reports/stock-movement-history**", async (route) => {
    await route.fulfill({ status: 200, contentType: "application/json", body: JSON.stringify({ asOf: "2026-09-01T00:00:00Z", windowDays: 30, items: [
      { productId: product.id, sku: product.sku, name: product.name, locationId: location.id, movementCount: 4, inboundUnits: 12, outboundUnits: 5, netUnits: 7, averageDailyOutbound: 0.17, lastMovementAt: "2026-08-31T00:00:00Z" },
    ] }) });
  });
  await page.route("**/api/v1/reports/supplier-performance**", async (route) => {
    await route.fulfill({ status: 200, contentType: "application/json", body: JSON.stringify({ asOf: "2026-09-01T00:00:00Z", windowDays: 30, items: [
      { supplierId: supplier.id, supplierCode: supplier.code, supplierName: supplier.name, orderCount: 2, orderedUnits: 6, receivedUnits: 3, openUnits: 3, orderedValueMinor: 75000, receivedValueMinor: 37500, averageLeadTimeDays: 4.5, completedOrderCount: 1, onTimeOrderCount: 1 },
    ] }) });
  });
  await page.route("**/api/v1/reports/replenishment-readiness**", async (route) => {
    await route.fulfill({ status: 200, contentType: "application/json", body: JSON.stringify({ asOf: "2026-09-01T00:00:00Z", windowDays: 30, items: [
      { productId: product.id, sku: product.sku, name: product.name, supplierId: supplier.id, unit: product.unit, onHand: 12, reorderPoint: 5, reorderQuantity: 10, targetStock: 15, suggestedQuantity: 3, outboundUnits: 5, averageDailyOutbound: 0.17, daysOfCover: 70.6, risk: "healthy" },
    ] }) });
  });
  await page.route("**/api/v1/replenishment/reviews**", async (route) => {
    await route.fulfill({ status: 200, contentType: "application/json", body: JSON.stringify({ items: [] }) });
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

  test("creates, submits, and receives a purchase order through the authenticated UI", async ({ page }) => {
    await mockAuthenticatedData(page);
    await mockPurchasing(page);

    await page.goto("/");
    await page.getByRole("button", { name: "Purchase orders", exact: true }).click();
    await expect(page.getByRole("heading", { name: "Purchase order workflow" })).toBeVisible();
    await page.getByRole("button", { name: "New order" }).click();

    await page.getByLabel("Order number").fill("E2E-PO-001");
    await page.getByRole("button", { name: "Create order" }).click();
    await expect(page.getByRole("heading", { name: "E2E-PO-001" })).toBeVisible();

    await page.getByRole("button", { name: "Submit order" }).click();
    await expect(page.getByRole("status")).toContainText("Order is now ordered.");

    await page.getByLabel("Quantity").fill("3");
    await page.getByRole("button", { name: "Receive into inventory" }).click();
    await expect(page.getByRole("status")).toContainText("Receipt committed");
    await expect(page.getByText("3 / 3 received")).toBeVisible();
  });

  test("loads the authenticated reports workspace with complete synthetic fixtures", async ({ page }) => {
    await mockAuthenticatedData(page);
    await mockReports(page);

    await page.goto("/");
    await page.getByRole("button", { name: "Reports", exact: true }).click();
    await expect(page.getByRole("heading", { name: "Reports & analytics" })).toBeVisible();
    await expect(page.getByText("Active products")).toBeVisible();
    await expect(page.getByText("Current on-hand value")).toBeVisible();
    await expect(page.getByText("₹1,500.00")).toBeVisible();
    await expect(page.getByText("Synthetic Supplier")).toBeVisible();
    await expect(page.getByText("0-30")).toBeVisible();
    await expect(page.getByText("Recent movement activity")).toBeVisible();
    await expect(page.getByText("Review stock risk before ordering")).toBeVisible();
    await expect(page.getByText("Review history")).toBeVisible();
    await expect(page.getByRole("status")).toContainText("0 reviews shown.");
  });
});
