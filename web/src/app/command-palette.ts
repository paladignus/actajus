import { buildRouteCatalog, loadRecentRoutes } from "./routes";
import type { RouteEntry } from "./types";

export function setupCommandPalette() {
  if (!document.body.classList.contains("dashboard-shell")) {
    return;
  }

  const catalog = buildRouteCatalog();
  if (catalog.length === 0) {
    return;
  }

  const isEditableTarget = (target: EventTarget | null) => {
    if (!(target instanceof HTMLElement)) {
      return false;
    }
    if (target.isContentEditable) {
      return true;
    }
    return Boolean(target.closest("input, textarea, select, [contenteditable='true']"));
  };

  const overlay = document.createElement("div");
  overlay.className = "command-palette";
  overlay.hidden = true;
  overlay.innerHTML = `
    <div class="command-palette__backdrop" data-command-palette-close></div>
    <section class="command-palette__dialog" role="dialog" aria-modal="true" aria-label="Busca global">
      <header class="command-palette__header">
        <div class="command-palette__search">
          <span class="command-palette__icon" aria-hidden="true">/</span>
          <input class="command-palette__input" type="search" placeholder="Buscar rota, area ou modulo">
        </div>
        <button class="command-palette__close" type="button" data-command-palette-close aria-label="Fechar busca global">Fechar</button>
      </header>
      <div class="command-palette__body">
        <div class="command-palette__section">
          <div class="command-palette__section-head">
            <strong class="command-palette__title">Resultados</strong>
            <span class="command-palette__hint">Use setas e Enter</span>
          </div>
          <div class="command-palette__results" data-command-palette-results></div>
        </div>
        <div class="command-palette__section">
          <div class="command-palette__section-head">
            <strong class="command-palette__title">Ultimas areas</strong>
          </div>
          <div class="command-palette__results" data-command-palette-recent></div>
        </div>
      </div>
    </section>
  `;
  document.body.appendChild(overlay);

  const input = overlay.querySelector<HTMLInputElement>(".command-palette__input");
  const results = overlay.querySelector<HTMLElement>("[data-command-palette-results]");
  const recent = overlay.querySelector<HTMLElement>("[data-command-palette-recent]");
  const openButtons = document.querySelectorAll<HTMLElement>("[data-command-palette-open]");
  const closeButtons = overlay.querySelectorAll<HTMLElement>("[data-command-palette-close]");
  if (!input || !results || !recent) {
    overlay.remove();
    return;
  }

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
      link.className = "command-palette__item";
      link.href = item.href;
      link.dataset.routeIndex = String(index);

      const label = document.createElement("span");
      label.className = "command-palette__item-label";
      label.textContent = item.label;

      const href = document.createElement("span");
      href.className = "command-palette__item-href";
      href.textContent = item.href;

      link.append(label, href);
      target.appendChild(link);
    });
  };

  const syncActive = () => {
    const links = Array.from(results.querySelectorAll<HTMLAnchorElement>(".command-palette__item"));
    links.forEach((link, index) => {
      const isActive = index === activeIndex;
      link.classList.toggle("command-palette__item--active", isActive);
      link.setAttribute("aria-selected", String(isActive));
    });
    if (activeIndex < 0 || activeIndex >= links.length) {
      return;
    }
    links[activeIndex].scrollIntoView({ block: "nearest" });
  };

  const updateResults = (items: RouteEntry[], emptyMessage: string) => {
    currentItems = items;
    activeIndex = items.length > 0 ? 0 : -1;
    renderLinks(items, results, emptyMessage);
    syncActive();
  };

  const runSearch = () => {
    const query = input.value.trim().toLowerCase();
    if (!query) {
      updateResults(catalog, "Nenhuma rota disponivel.");
      return;
    }
    const filtered = catalog.filter((item) => {
      const haystack = [item.label, item.href, ...item.keywords].join(" ").toLowerCase();
      return haystack.includes(query);
    });
    updateResults(filtered, "Nenhuma rota encontrada.");
  };

  const open = () => {
    overlay.hidden = false;
    document.body.classList.add("command-palette-open");
    input.focus();
    input.select();
    renderLinks(loadRecentRoutes(), recent, "Nenhuma area visitada recentemente.");
    runSearch();
  };

  const close = () => {
    overlay.hidden = true;
    document.body.classList.remove("command-palette-open");
    input.value = "";
  };

  openButtons.forEach((button) => button.addEventListener("click", open));
  closeButtons.forEach((button) => button.addEventListener("click", close));
  input.addEventListener("input", runSearch);
  input.addEventListener("keydown", (event) => {
    if (currentItems.length === 0) {
      if (event.key === "Escape") {
        event.preventDefault();
        close();
      }
      return;
    }
    if (event.key === "ArrowDown") {
      event.preventDefault();
      activeIndex = (activeIndex + 1 + currentItems.length) % currentItems.length;
      syncActive();
      return;
    }
    if (event.key === "ArrowUp") {
      event.preventDefault();
      activeIndex = (activeIndex - 1 + currentItems.length) % currentItems.length;
      syncActive();
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
      event.preventDefault();
      close();
    }
  });
  results.addEventListener("mousemove", (event) => {
    const target = event.target;
    if (!(target instanceof Element)) {
      return;
    }
    const link = target.closest<HTMLAnchorElement>(".command-palette__item");
    if (!link) {
      return;
    }
    const nextIndex = Number(link.dataset.routeIndex);
    if (Number.isNaN(nextIndex)) {
      return;
    }
    activeIndex = nextIndex;
    syncActive();
  });
  document.addEventListener("keydown", (event) => {
    if (event.key === "/" && !isEditableTarget(event.target)) {
      event.preventDefault();
      open();
      return;
    }
    if ((event.ctrlKey || event.metaKey) && event.key.toLowerCase() === "k") {
      event.preventDefault();
      open();
      return;
    }
    if (event.key === "Escape" && !overlay.hidden) {
      event.preventDefault();
      close();
    }
  });
}
