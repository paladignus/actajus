export function setupSidebarToggle() {
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
