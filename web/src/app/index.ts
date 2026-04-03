import { loadBootstrap } from "./bootstrap";
import { setupCommandPalette } from "./command-palette";
import { setupFormConfirmations } from "./confirm";
import { setupSocialMediaForm } from "./company-form";
import { setupErrorRouteTools } from "./error-tools";
import { enhanceForms } from "./forms";
import { trackRecentRoute } from "./routes";
import { setupSidebarToggle } from "./sidebar";
import { setupToasts } from "./toast";

export async function initApp() {
  trackRecentRoute();
  enhanceForms();
  setupFormConfirmations();
  setupToasts();
  setupSidebarToggle();
  setupCommandPalette();
  setupErrorRouteTools();
  setupSocialMediaForm();

  const payload = await loadBootstrap();
  if (!payload) {
    return;
  }
  document.body.dataset.authenticated = String(Boolean(payload.authenticated));
}
