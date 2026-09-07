import { test, expect } from "@playwright/test";

test.describe("signed-out workspace", () => {
  test.beforeEach(async ({ page }) => {
    await page.route("**/api/v1/auth/me", async (route) => {
      await route.fulfill({
        status: 401,
        contentType: "application/json",
        body: JSON.stringify({ error: "unauthorized" }),
      });
    });
  });

  test("renders a labeled sign-in form", async ({ page }) => {
    await page.goto("/");

    await expect(page.getByRole("heading", { name: "Sign in to StockPilot" })).toBeVisible();
    await expect(page.getByLabel("Email")).toHaveAttribute("autocomplete", "username");
    await expect(page.getByLabel("Password")).toHaveAttribute("autocomplete", "current-password");
    await expect(page.getByRole("button", { name: "Sign in" })).toBeDisabled();
  });

  test("supports keyboard traversal through the sign-in controls", async ({ page }) => {
    await page.goto("/");

    const homeLink = page.getByRole("link", { name: "StockPilot home" });
    await homeLink.focus();
    await expect(homeLink).toBeFocused();

    await page.keyboard.press("Tab");
    await expect(page.getByLabel("Email")).toBeFocused();
    await page.keyboard.press("Tab");
    await expect(page.getByLabel("Password")).toBeFocused();
    await page.keyboard.press("Tab");
    await expect(page.getByRole("button", { name: "Sign in" })).toBeFocused();
  });
});
