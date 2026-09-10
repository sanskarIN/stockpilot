// Replace the existing mockReports function with this:
async function mockReports(page: Page) {
  // Generic handler for all /api/v1/reports/* endpoints so query-string variants are matched
  await page.route("**/api/v1/reports/**", async (route) => {
    const url = route.request().url();

    if (url.includes("/inventory-valuation")) {
      await route.fulfill({
        status: 200,
        contentType: "application/json",
        body: JSON.stringify({
          items: [
            {
              productId: product.id,
              sku: product.sku,
              name: product.name,
              unit: product.unit,
              onHand: 12,
              unitCostMinor: 12500,
              currency: "INR",
              valueMinor: 150000, // 1500.00 INR
            },
          ],
          totals: [{ currency: "INR", valueMinor: 150000 }],
        }),
      });
      return;
    }

    if (url.includes("/inventory-aging")) {
      await route.fulfill({
        status: 200,
        contentType: "application/json",
        body: JSON.stringify({
          items: [
            {
              productId: product.id,
              sku: product.sku,
              name: product.name,
              locationId: location.id,
              quantity: 12,
              ageDays: 12,
              bucket: "0-30",
              asOf: "2026-09-01T00:00:00Z",
              lastMovementAt: "2026-08-20T00:00:00Z",
            },
          ],
        }),
      });
      return;
    }

    if (url.includes("/stock-movement-history")) {
      await route.fulfill({
        status: 200,
        contentType: "application/json",
        body: JSON.stringify({
          asOf: "2026-09-01T00:00:00Z",
          windowDays: 30,
          items: [
            {
              productId: product.id,
              sku: product.sku,
              name: product.name,
              locationId: location.id,
              movementCount: 4,
              inboundUnits: 12,
              outboundUnits: 5,
              netUnits: 7,
              averageDailyOutbound: 0.17,
              lastMovementAt: "2026-08-20T00:00:00Z",
            },
          ],
        }),
      });
      return;
    }

    if (url.includes("/supplier-performance")) {
      await route.fulfill({
        status: 200,
        contentType: "application/json",
        body: JSON.stringify({
          asOf: "2026-09-01T00:00:00Z",
          windowDays: 30,
          items: [
            {
              supplierId: supplier.id,
              supplierCode: supplier.code,
              supplierName: supplier.name,
              orderCount: 2,
              orderedUnits: 6,
              receivedUnits: 3,
              openUnits: 3,
              orderedValueMinor: 75000,
              receivedValueMinor: 37500,
            },
          ],
        }),
      });
      return;
    }

    if (url.includes("/replenishment-readiness")) {
      await route.fulfill({
        status: 200,
        contentType: "application/json",
        body: JSON.stringify({
          asOf: "2026-09-01T00:00:00Z",
          windowDays: 30,
          items: [
            {
              productId: product.id,
              sku: product.sku,
              name: product.name,
              supplierId: supplier.id,
              unit: product.unit,
              onHand: 12,
              reorderPoint: 5,
              reorderQuantity: 10,
              targetStock: 15,
              suggestedQuantity: 3,
            },
          ],
        }),
      });
      return;
    }

    // Fallback empty response so UI can still render
    await route.fulfill({ status: 200, contentType: "application/json", body: JSON.stringify({ items: [] }) });
  });

  // Keep replenishment reviews mock
  await page.route("**/api/v1/replenishment/reviews**", async (route) => {
    await route.fulfill({ status: 200, contentType: "application/json", body: JSON.stringify({ items: [] }) });
  });

  // Keep reports/overview explicit (used by the UI)
  await page.route("**/api/v1/reports/overview", async (route) => {
    await route.fulfill({
      status: 200,
      contentType: "application/json",
      body: JSON.stringify({
        inventory: {
          productCount: 1,
          activeProductCount: 1,
          activeWarehouseCount: 1,
          activeLocationCount: 1,
          totalUnits: 12,
          lowStockBalanceCount: 0,
          outOfStockCount: 0,
        },
        purchasing: {
          totalOrders: 2,
          draftOrders: 1,
          orderedOrders: 1,
          partiallyReceivedOrders: 0,
          receivedOrders: 0,
          cancelledOrders: 0,
          outstandingUnits: 3,
        },
      }),
    });
  });
}

// Then, update the failing assertion near the end of the test to:
await expect(page.getByText("₹1,500.00", { exact: true }).first()).toBeVisible({ timeout: 10_000 });
