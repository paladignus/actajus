import { buildRouteCatalog, loadRecentRoutes } from "./routes";
import type { RouteEntry } from "./types";

export function setupErrorRouteTools() {
  const searchRoot = document.querySelector<HTMLElement>("[data-error-route-search]");
  const results = document.querySelector<HTMLElement>("[data-error-route-results]");
  const input = document.querySelector<HTMLInputElement>("[data-error-route-input]");
  const recentRoot = document.querySelector<HTMLElement>("[data-error-recent-routes]");
  if (!searchRoot || !results || !input || !recentRoot) {
    return;
  }

  const catalog = buildRouteCatalog();
  let activeIndex = -1;
  let currentItems: RouteEntry[] = [];

  const renderLinks = (items: RouteEntry[], target: HTMLElement, emptyMessage: string) => {
    target.innerHTML = "";
    if (items.length === 0) {
      const empty = document.createElement("p");
      empty.className = "dashboard-empty";
      empty.textContent = emptyMessage;
      target.appendChild(empty);
      return;
    }
    items.forEach((item, index) => {
      const link = document.createElement("a");
      link.className = "error-card__shortcut";
      link.href = item.href;
      link.dataset.routeIndex = String(index);
      link.setAttribute("tabindex", "-1");

      const label = document.createElement("span");
      label.className = "error-card__shortcut-label";
      label.textContent = item.label;

      const href = document.createElement("span");
      href.className = "error-card__shortcut-href";
      href.textContent = item.href;

      link.append(label, href);
      target.appendChild(link);
    });
  };

  const syncActiveRoute = () => {
    const links = Array.from(results.querySelectorAll<HTMLAnchorElement>(".error-card__shortcut"));
    links.forEach((link, index) => {
      const isActive = index === activeIndex;
      link.classList.toggle("error-card__shortcut--active", isActive);
      link.setAttribute("aria-selected", String(isActive));
    });
    if (activeIndex < 0 || activeIndex >= links.length) {
      return;
    }
    links[activeIndex].scrollIntoView({ block: "nearest", inline: "nearest" });
  };

  const updateResults = (items: RouteEntry[], emptyMessage: string) => {
    currentItems = items;
    activeIndex = items.length > 0 ? 0 : -1;
    renderLinks(items, results, emptyMessage);
    syncActiveRoute();
  };

  const focusSearch = () => {
    searchRoot.classList.add("error-card__search--focused");
    input.focus();
    input.select();
  };

  const clearSearchFocus = () => {
    if (document.activeElement === input) {
      input.blur();
    }
    searchRoot.classList.remove("error-card__search--focused");
  };

  const isEditableTarget = (target: EventTarget | null) => {
    if (!(target instanceof HTMLElement)) {
      return false;
    }
    if (target.isContentEditable) {
      return true;
    }
    return Boolean(target.closest("input, textarea, select, [contenteditable='true']"));
  };

  const search = () => {
    const query = input.value.trim().toLowerCase();
    if (!query) {
      updateResults(catalog.slice(0, 6), "Nenhuma rota disponivel.");
      return;
    }
    const filtered = catalog.filter((item) => {
      const haystack = [item.label, item.href, ...item.keywords].join(" ").toLowerCase();
      return haystack.includes(query);
    });
    updateResults(filtered, "Nenhuma rota encontrada para essa busca.");
  };

  updateResults(catalog.slice(0, 6), "Nenhuma rota disponivel.");
  renderLinks(loadRecentRoutes(), recentRoot, "Nenhuma area visitada recentemente.");
  input.addEventListener("focus", () => searchRoot.classList.add("error-card__search--focused"));
  input.addEventListener("blur", () => searchRoot.classList.remove("error-card__search--focused"));
  input.addEventListener("input", search);
  input.addEventListener("keydown", (event) => {
    if (currentItems.length === 0) {
      if (event.key === "Escape" && input.value) {
        event.preventDefault();
        input.value = "";
        search();
      }
      return;
    }
    if (event.key === "ArrowDown") {
      event.preventDefault();
      activeIndex = (activeIndex + 1 + currentItems.length) % currentItems.length;
      syncActiveRoute();
      return;
    }
    if (event.key === "ArrowUp") {
      event.preventDefault();
      activeIndex = (activeIndex - 1 + currentItems.length) % currentItems.length;
      syncActiveRoute();
      return;
    }
    if (event.key === "Enter") {
      event.preventDefault();
      const nextRoute = currentItems[Math.max(activeIndex, 0)];
      if (nextRoute) {
        window.location.assign(nextRoute.href);
      }
      return;
    }
    if (event.key === "Escape") {
      if (input.value) {
        event.preventDefault();
        input.value = "";
        search();
        return;
      }
      clearSearchFocus();
      activeIndex = -1;
      syncActiveRoute();
    }
  });
  document.addEventListener("keydown", (event) => {
    if (event.key === "/" && !isEditableTarget(event.target)) {
      event.preventDefault();
      focusSearch();
      return;
    }
    if ((event.ctrlKey || event.metaKey) && event.key.toLowerCase() === "k") {
      event.preventDefault();
      focusSearch();
    }
  });
  results.addEventListener("mousemove", (event) => {
    const target = event.target;
    if (!(target instanceof Element)) {
      return;
    }
    const link = target.closest<HTMLAnchorElement>(".error-card__shortcut");
    if (!link) {
      return;
    }
    const nextIndex = Number(link.dataset.routeIndex);
    if (Number.isNaN(nextIndex) || nextIndex === activeIndex) {
      return;
    }
    activeIndex = nextIndex;
    syncActiveRoute();
  });
  searchRoot.addEventListener("click", (event) => {
    const target = event.target;
    if (target instanceof HTMLAnchorElement || target instanceof HTMLButtonElement) {
      return;
    }
    if (target instanceof HTMLElement && target.closest(".error-card__route-results")) {
      return;
    }
    focusSearch();
  });
}
