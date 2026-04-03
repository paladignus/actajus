import "./styles/app.css";

declare global {
  interface Window {
    __ACTAJUS_BOOTSTRAP__?: {
      appName?: string;
      apiBaseURL?: string;
      authenticated?: boolean;
      canManageSessions?: boolean;
      canManageRBAC?: boolean;
    };
    actajusNotify?: (toast: ToastInput) => void;
  }
}

type ToastKind = "success" | "error" | "info" | "alert";

type ToastInput = {
  kind?: ToastKind;
  title?: string;
  message: string;
  source?: string;
  code?: string;
  timeoutMs?: number;
};

type RouteEntry = {
  href: string;
  label: string;
  keywords: string[];
};

const RECENT_ROUTES_KEY = "actajus_recent_routes";

function enhanceForms() {
  const loginEmail = document.querySelector<HTMLInputElement>('input[name="email"]');
  if (loginEmail && !loginEmail.value) {
    loginEmail.focus();
  }

  const search = document.querySelector<HTMLInputElement>('input[type="search"][name="email"]');
  if (search && window.__ACTAJUS_BOOTSTRAP__?.canManageSessions) {
    search.setAttribute("aria-label", "Filtrar sessoes por email");
  }
}

function buildRouteCatalog() {
  const bootstrap = window.__ACTAJUS_BOOTSTRAP__;
  const authenticated = Boolean(bootstrap?.authenticated);
  const routes: RouteEntry[] = [
    { href: "/", label: "Inicio", keywords: ["home", "inicio", "painel"] },
    { href: "/login", label: "Login", keywords: ["login", "entrar", "acesso"] }
  ];

  if (!authenticated) {
    return routes;
  }

  routes.push(
    { href: "/sessions", label: "Sessoes", keywords: ["sessoes", "sessions", "dispositivos"] },
    { href: "/companies", label: "Companies", keywords: ["companies", "empresa", "empresas"] },
    { href: "/companies/new", label: "Nova company", keywords: ["nova company", "criar company", "empresa"] }
  );

  if (bootstrap?.canManageRBAC) {
    routes.push(
      { href: "/users", label: "Usuarios", keywords: ["usuarios", "users"] },
      { href: "/roles", label: "Roles", keywords: ["roles", "rbac"] },
      { href: "/permissions", label: "Permissoes", keywords: ["permissoes", "permissions", "rbac"] }
    );
  }

  return routes;
}

function trackRecentRoute() {
  const path = window.location.pathname;
  if (!path || path.startsWith("/assets") || path === "/bootstrap") {
    return;
  }

  const catalog = buildRouteCatalog();
  const match = catalog.find((item) => item.href === path);
  const route: RouteEntry = match || {
    href: path,
    label: path === "/" ? "Inicio" : path,
    keywords: [path]
  };

  const current = loadRecentRoutes().filter((item) => item.href !== route.href);
  current.unshift(route);
  localStorage.setItem(RECENT_ROUTES_KEY, JSON.stringify(current.slice(0, 6)));
}

function loadRecentRoutes(): RouteEntry[] {
  try {
    const raw = localStorage.getItem(RECENT_ROUTES_KEY);
    if (!raw) {
      return [];
    }
    const parsed = JSON.parse(raw) as RouteEntry[];
    return Array.isArray(parsed) ? parsed.filter((item) => item?.href && item?.label) : [];
  } catch {
    return [];
  }
}

function setupToasts() {
  const notifier = createToastNotifier();
  window.actajusNotify = notifier;

  const flashes = document.querySelectorAll<HTMLElement>(".dashboard-alert[data-flash-kind]");
  if (flashes.length === 0) {
    return;
  }

  flashes.forEach((flash, index) => {
    const kind = normalizeToastKind(flash.dataset.flashKind);
    notifier({
      kind,
      source: flash.dataset.flashSource || undefined,
      code: flash.dataset.flashCode || undefined,
      message: flash.textContent?.trim() || "",
      timeoutMs: toastTimeout(kind) + index * 400
    });
  });
}

function createToastNotifier() {
  let viewport = document.querySelector<HTMLElement>("[data-toast-viewport]");

  const ensureViewport = () => {
    if (viewport) {
      return viewport;
    }
    viewport = document.createElement("div");
    viewport.className = "toast-viewport";
    viewport.setAttribute("data-toast-viewport", "true");
    viewport.setAttribute("aria-live", "polite");
    viewport.setAttribute("aria-atomic", "true");
    document.body.appendChild(viewport);
    return viewport;
  };

  return (input: ToastInput) => {
    const kind = normalizeToastKind(input.kind);
    const holder = ensureViewport();
    const toast = document.createElement("section");
    toast.className = `toast toast--${kind}`;
    toast.setAttribute("role", kind === "error" ? "alert" : "status");
    toast.dataset.toastKind = kind;
    toast.dataset.toastCode = input.code || "";
    toast.dataset.toastSource = input.source || "";

    const progress = document.createElement("div");
    progress.className = "toast__progress";
    progress.style.setProperty("--toast-duration", `${input.timeoutMs ?? toastTimeout(kind)}ms`);

    const header = document.createElement("div");
    header.className = "toast__header";

    const leading = document.createElement("div");
    leading.className = "toast__leading";

    const icon = document.createElement("span");
    icon.className = "toast__icon";
    icon.textContent = toastIcon(kind);
    icon.setAttribute("aria-hidden", "true");

    const text = document.createElement("div");
    text.className = "toast__text";

    const title = document.createElement("strong");
    title.className = "toast__title";
    title.textContent = input.title || toastTitle(kind, input.source || input.code || "Actajus");

    const meta = document.createElement("span");
    meta.className = "toast__meta";
    meta.textContent = input.code || input.source || "";

    const closeButton = document.createElement("button");
    closeButton.className = "toast__close";
    closeButton.type = "button";
    closeButton.setAttribute("aria-label", "Fechar notificacao");
    closeButton.textContent = "Fechar";

    const body = document.createElement("p");
    body.className = "toast__message";
    body.textContent = input.message;

    text.append(title);
    if (meta.textContent) {
      text.append(meta);
    }
    leading.append(icon, text);
    header.append(leading, closeButton);
    toast.append(progress, header, body);
    holder.prepend(toast);

    let dismissed = false;
    let startX = 0;
    let currentX = 0;
    let dragging = false;
    let activePointerId: number | null = null;

    const dismiss = () => {
      if (dismissed) {
        return;
      }
      dismissed = true;
      toast.classList.add("toast--closing");
      window.setTimeout(() => {
        toast.remove();
        if (viewport && viewport.childElementCount === 0) {
          viewport.remove();
          viewport = null;
        }
      }, 180);
    };

    const onPointerMove = (event: PointerEvent) => {
      if (!dragging) {
        return;
      }
      currentX = event.clientX - startX;
      toast.style.transform = `translateX(${currentX}px)`;
      toast.style.opacity = String(Math.max(0.35, 1 - Math.abs(currentX) / 180));
    };

    const onPointerUp = () => {
      if (!dragging) {
        return;
      }
      dragging = false;
      if (activePointerId !== null) {
        toast.releasePointerCapture?.(activePointerId);
      }
      activePointerId = null;
      if (Math.abs(currentX) > 110) {
        dismiss();
      } else {
        toast.style.transform = "";
        toast.style.opacity = "";
      }
      currentX = 0;
    };

    toast.addEventListener("pointerdown", (event) => {
      dragging = true;
      startX = event.clientX;
      activePointerId = event.pointerId;
      toast.setPointerCapture?.(event.pointerId);
    });
    toast.addEventListener("pointermove", onPointerMove);
    toast.addEventListener("pointerup", onPointerUp);
    toast.addEventListener("pointercancel", onPointerUp);
    closeButton.addEventListener("click", dismiss);
    window.setTimeout(dismiss, input.timeoutMs ?? toastTimeout(kind));
  };
}

function normalizeToastKind(kind?: string): ToastKind {
  if (kind === "success" || kind === "error" || kind === "alert" || kind === "info") {
    return kind;
  }
  return "info";
}

function toastTitle(kind: ToastKind, source: string) {
  switch (kind) {
    case "success":
      return `Sucesso: ${source}`;
    case "error":
      return `Erro: ${source}`;
    case "alert":
      return `Alerta: ${source}`;
    case "info":
      return `Info: ${source}`;
    default:
      return source;
  }
}

function toastIcon(kind: ToastKind) {
  switch (kind) {
    case "success":
      return "OK";
    case "error":
      return "ER";
    case "alert":
      return "AL";
    case "info":
      return "IN";
  }
}

function toastTimeout(kind: ToastKind) {
  switch (kind) {
    case "error":
      return 6400;
    case "alert":
      return 5600;
    case "success":
      return 3800;
    case "info":
      return 4200;
  }
}

function setupSidebarToggle() {
  const shell = document.querySelector<HTMLElement>(".dashboard-shell");
  const toggle = document.querySelector<HTMLButtonElement>("[data-sidebar-toggle]");
  if (shell && toggle) {
    toggle.addEventListener("click", () => {
      shell.classList.toggle("dashboard-shell--sidebar-open");
    });
  }

  const groups = document.querySelectorAll<HTMLElement>("[data-sidebar-group]");
  groups.forEach((group) => {
    const submenuToggle = group.querySelector<HTMLButtonElement>("[data-sidebar-submenu]");
    if (!submenuToggle) {
      return;
    }
    submenuToggle.addEventListener("click", () => {
      const isOpen = group.classList.toggle("dashboard-sidebar__group--open");
      submenuToggle.setAttribute("aria-expanded", String(isOpen));
    });
  });
}

function setupErrorRouteTools() {
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
  input.addEventListener("focus", () => {
    searchRoot.classList.add("error-card__search--focused");
  });
  input.addEventListener("blur", () => {
    searchRoot.classList.remove("error-card__search--focused");
  });
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

function setupCommandPalette() {
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

  openButtons.forEach((button) => {
    button.addEventListener("click", open);
  });
  closeButtons.forEach((button) => {
    button.addEventListener("click", close);
  });
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

function setupSocialMediaForm() {
  const list = document.querySelector<HTMLElement>("[data-social-media-list]");
  const addButton = document.querySelector<HTMLButtonElement>("[data-social-media-add]");
  const template = document.querySelector<HTMLTemplateElement>("#social-media-item-template");
  const form = document.querySelector<HTMLFormElement>("[data-company-form]");
  if (!list || !addButton || !template) {
    return;
  }

  let isDirty = false;
  let warnedAboutDirty = false;

  const markDirty = () => {
    isDirty = true;
  };

  if (form) {
    form.addEventListener("input", markDirty);
    form.addEventListener("change", markDirty);
    form.addEventListener("submit", () => {
      isDirty = false;
    });
    window.addEventListener("beforeunload", (event) => {
      if (!isDirty) {
        return;
      }
      event.preventDefault();
      event.returnValue = "";
    });
    document.addEventListener("click", (event) => {
      if (!isDirty || warnedAboutDirty) {
        return;
      }
      const target = event.target;
      if (!(target instanceof Element)) {
        return;
      }
      const anchor = target.closest<HTMLAnchorElement>("a[href]");
      if (!anchor || anchor.target === "_blank") {
        return;
      }
      const href = anchor.getAttribute("href");
      if (!href || href.startsWith("#")) {
        return;
      }
      warnedAboutDirty = true;
      window.actajusNotify?.({
        kind: "alert",
        source: "companies.form.leave",
        code: "companies_form_unsaved_changes",
        message: "Voce tem alteracoes nao salvas. Se sair agora, elas serao perdidas.",
        timeoutMs: 5200
      });
      window.setTimeout(() => {
        warnedAboutDirty = false;
      }, 2500);
    });
  }

  const bindRemoveHandlers = () => {
    const items = list.querySelectorAll<HTMLElement>("[data-social-media-item]");
    items.forEach((item) => {
      const removeButton = item.querySelector<HTMLButtonElement>("[data-social-media-remove]");
      if (!removeButton || removeButton.dataset.bound === "true") {
        return;
      }
      removeButton.dataset.bound = "true";
      removeButton.addEventListener("click", () => {
        if (list.querySelectorAll("[data-social-media-item]").length <= 1) {
          const inputs = item.querySelectorAll<HTMLInputElement>('input[name="social_platform"], input[name="social_url"], input[name="social_media_id"]');
          inputs.forEach((input) => {
            input.value = "";
          });
          markDirty();
          window.actajusNotify?.({
            kind: "info",
            source: "companies.social_media.reset",
            code: "companies_social_media_reset",
            message: "A social media restante foi limpa para edicao."
          });
          return;
        }
        item.remove();
        markDirty();
        window.actajusNotify?.({
          kind: "alert",
          source: "companies.social_media.remove",
          code: "companies_social_media_removed",
          message: "Social media removida do formulario."
        });
      });
    });
  };

  addButton.addEventListener("click", () => {
    const fragment = template.content.cloneNode(true);
    list.appendChild(fragment);
    bindRemoveHandlers();
    markDirty();
    window.actajusNotify?.({
      kind: "success",
      source: "companies.social_media.add",
      code: "companies_social_media_added",
      message: "Nova social media adicionada ao formulario."
    });
  });

  bindRemoveHandlers();
}

async function loadBootstrap() {
  const response = await fetch("/bootstrap", {
    headers: { Accept: "application/json" }
  });
  if (!response.ok) {
    return null;
  }
  return response.json() as Promise<{ authenticated?: boolean }>;
}

async function main() {
  trackRecentRoute();
  enhanceForms();
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

void main();
