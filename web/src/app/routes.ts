import type { RouteEntry } from "./types";

const RECENT_ROUTES_KEY = "actajus_recent_routes";

export function buildRouteCatalog(): RouteEntry[] {
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

export function trackRecentRoute() {
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

export function loadRecentRoutes(): RouteEntry[] {
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
