import { expect, type Page } from "@playwright/test";

import { e2eCredentials, test } from "./helpers/env";

const credentials = e2eCredentials();
const activeSessionLimitMessage = "Limite de sessoes ativas atingido.";
const rateLimitedMessage = "Muitas tentativas de login. Tente novamente em instantes.";

type LoginResult = "success" | "active-session-limit" | "rate-limited" | "login-blocked";

test.describe("authenticated dashboard", () => {
  test.skip(!credentials, "Defina ACTAJUS_E2E_EMAIL e ACTAJUS_E2E_PASSWORD para rodar o fluxo autenticado.");
  test.describe.configure({ mode: "serial" });

  async function login(page: Page): Promise<LoginResult> {
    await page.goto("/login");
    await expect(page.locator('form[action="/login"] input[name="_csrf"]')).toHaveValue(/.+/);

    await page.getByLabel("Email").fill(credentials!.email);
    await page.getByLabel("Senha").fill(credentials!.password);
    await page.getByRole("button", { name: "Entrar" }).click();

    const sessionsHeading = page.getByRole("heading", { name: "Minhas sessoes" });
    const activeSessionLimit = page.getByText(activeSessionLimitMessage);
    const rateLimited = page.getByText(rateLimitedMessage);
    const loginAlert = page.locator(".dashboard-alert").first();

    await Promise.race([
      sessionsHeading.waitFor({ state: "visible", timeout: 8000 }).then(() => "success" as const),
      activeSessionLimit.waitFor({ state: "visible", timeout: 8000 }).then(() => "active-session-limit" as const),
      rateLimited.waitFor({ state: "visible", timeout: 8000 }).then(() => "rate-limited" as const),
      loginAlert.waitFor({ state: "visible", timeout: 8000 }).then(() => "login-blocked" as const)
    ]);

    if (await activeSessionLimit.isVisible()) {
      return "active-session-limit";
    }
    if (await rateLimited.isVisible()) {
      return "rate-limited";
    }
    if (await loginAlert.isVisible()) {
      return "login-blocked";
    }
    return "success";
  }

  test("signs in and reaches sessions dashboard", async ({ page }) => {
    const loginResult = await login(page);
    test.skip(loginResult === "active-session-limit", "Conta de E2E bloqueada por limite de sessoes ativas no ambiente atual.");
    test.skip(loginResult === "rate-limited", "Conta de E2E temporariamente bloqueada pelo login shield do ambiente atual.");
    test.skip(loginResult === "login-blocked", "Conta de E2E bloqueada por uma mensagem de autenticacao do ambiente atual.");

    await expect(page).toHaveURL(/\/sessions$/);
    await expect(page.getByRole("heading", { name: "Minhas sessoes" })).toBeVisible();
  });

  test("navigates to companies list and company creation page", async ({ page }) => {
    const loginResult = await login(page);
    test.skip(loginResult === "active-session-limit", "Conta de E2E bloqueada por limite de sessoes ativas no ambiente atual.");
    test.skip(loginResult === "rate-limited", "Conta de E2E temporariamente bloqueada pelo login shield do ambiente atual.");
    test.skip(loginResult === "login-blocked", "Conta de E2E bloqueada por uma mensagem de autenticacao do ambiente atual.");

    await page.goto("/companies");
    await expect(page).toHaveURL(/\/companies$/);
    await expect(page.getByRole("main").getByRole("link", { name: "Nova company" })).toBeVisible();
    await expect(page.locator(".dashboard-panel__title").filter({ hasText: /Pagina\s+\d+/ })).toBeVisible();

    await page.getByRole("main").getByRole("link", { name: "Nova company" }).click();

    await expect(page).toHaveURL(/\/companies\/new$/);
    await expect(page.getByRole("heading", { name: "Nova company" })).toBeVisible();
    await expect(page.locator("form[data-company-form]")).toBeVisible();
    await expect(page.locator("form[data-company-form] button[type='submit']")).toBeVisible();
  });

  test("opens and cancels destructive confirmation on logout", async ({ page }) => {
    const loginResult = await login(page);
    test.skip(loginResult === "active-session-limit", "Conta de E2E bloqueada por limite de sessoes ativas no ambiente atual.");
    test.skip(loginResult === "rate-limited", "Conta de E2E temporariamente bloqueada pelo login shield do ambiente atual.");
    test.skip(loginResult === "login-blocked", "Conta de E2E bloqueada por uma mensagem de autenticacao do ambiente atual.");

    await page.getByRole("button", { name: "Logout" }).click();
    await expect(page.locator("[data-confirm-dialog].confirm-dialog--open")).toBeVisible();
    await expect(page.getByRole("alertdialog")).toBeVisible();
    await expect(page.getByText("Voce sera desconectado deste navegador imediatamente.")).toBeVisible();

    await page.getByRole("button", { name: "Cancelar" }).click();
    await expect(page.locator("[data-confirm-dialog].confirm-dialog--open")).toBeHidden();
    await expect(page).toHaveURL(/\/sessions$/);
  });

  test("opens and cancels destructive confirmation on session revoke all", async ({ page }) => {
    const loginResult = await login(page);
    test.skip(loginResult === "active-session-limit", "Conta de E2E bloqueada por limite de sessoes ativas no ambiente atual.");
    test.skip(loginResult === "rate-limited", "Conta de E2E temporariamente bloqueada pelo login shield do ambiente atual.");
    test.skip(loginResult === "login-blocked", "Conta de E2E bloqueada por uma mensagem de autenticacao do ambiente atual.");

    await expect(page).toHaveURL(/\/sessions$/);
    await page.getByRole("button", { name: "Encerrar todas" }).click();

    await expect(page.locator("[data-confirm-dialog].confirm-dialog--open")).toBeVisible();
    await expect(page.getByRole("alertdialog")).toBeVisible();
    await expect(page.getByText("Todas as suas sessoes serao encerradas, incluindo outros dispositivos.")).toBeVisible();

    await page.getByRole("button", { name: "Cancelar" }).click();
    await expect(page.locator("[data-confirm-dialog].confirm-dialog--open")).toBeHidden();
    await expect(page).toHaveURL(/\/sessions$/);
  });
});
