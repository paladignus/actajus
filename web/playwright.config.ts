import { defineConfig, devices } from "@playwright/test";

const baseURL = process.env.ACTAJUS_E2E_BASE_URL || "http://127.0.0.1:8080";
const reuseExistingServer = process.env.ACTAJUS_E2E_REUSE_SERVER === "true";
const executablePath = process.env.ACTAJUS_E2E_CHROMIUM_PATH || "/usr/bin/chromium";

export default defineConfig({
  testDir: "./tests/e2e",
  fullyParallel: true,
  forbidOnly: !!process.env.CI,
  retries: process.env.CI ? 2 : 0,
  workers: process.env.CI ? 1 : undefined,
  reporter: [["list"], ["html", { open: "never" }]],
  use: {
    baseURL,
    trace: "retain-on-failure",
    screenshot: "only-on-failure",
    video: "retain-on-failure"
  },
  webServer: reuseExistingServer
    ? undefined
    : {
        command: "go run ../cmd/http/server.go",
        url: `${baseURL}/health`,
        reuseExistingServer: !process.env.CI,
        stdout: "pipe",
        stderr: "pipe",
        timeout: 120_000
      },
  projects: [
    {
      name: "chromium",
      use: {
        ...devices["Desktop Chrome"],
        channel: undefined,
        executablePath,
        launchOptions: {
          args: ["--no-sandbox", "--disable-setuid-sandbox", "--disable-dev-shm-usage", "--disable-gpu"]
        }
      }
    }
  ]
});
