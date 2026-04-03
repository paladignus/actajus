export function setupSocialMediaForm() {
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
