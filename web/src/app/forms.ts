export function enhanceForms() {
  const loginEmail = document.querySelector<HTMLInputElement>('input[name="email"]');
  if (loginEmail && !loginEmail.value) {
    loginEmail.focus();
  }

  const search = document.querySelector<HTMLInputElement>('input[type="search"][name="email"]');
  if (search && window.__ACTAJUS_BOOTSTRAP__?.canManageSessions) {
    search.setAttribute("aria-label", "Filtrar sessoes por email");
  }
}
