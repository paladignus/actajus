import { expect } from "@playwright/test";

import { test } from "./helpers/env";

test.describe("public smoke", () => {
  test("renders login page", async ({ page }) => {
    await page.goto("/login");

    await expect(page).toHaveURL(/\/login$/);
    await expect(page.getByRole("heading", { name: "Entrar" })).toBeVisible();
    await expect(page.getByLabel("Email")).toBeVisible();
    await expect(page.getByLabel("Senha")).toBeVisible();
  });

  test("renders 404 page for invalid route", async ({ page }) => {
    await page.goto("/rota-invalida-playwright");

    await expect(page).toHaveURL(/\/rota-invalida-playwright$/);
    await expect(page.getByRole("heading", { name: "404" })).toBeVisible();
  });

  test("loads favicon and bootstrap endpoint", async ({ request }) => {
    const favicon = await request.get("/favicon.svg");
    expect(favicon.ok()).toBeTruthy();
    expect(favicon.headers()["content-type"] || "").toContain("image/svg+xml");

    const bootstrap = await request.get("/bootstrap");
    expect(bootstrap.ok()).toBeTruthy();
    const payload = await bootstrap.json();
    expect(payload).toHaveProperty("authenticated");
  });
});
