type PendingSubmit = {
  form: HTMLFormElement;
  submitter: HTMLElement | null;
};

export function setupFormConfirmations() {
  let pending: PendingSubmit | null = null;
  let confirmedForm: HTMLFormElement | null = null;

  const dialog = ensureDialog();
  const title = dialog.querySelector<HTMLElement>("[data-confirm-title]");
  const message = dialog.querySelector<HTMLElement>("[data-confirm-message]");
  const confirmButton = dialog.querySelector<HTMLButtonElement>("[data-confirm-accept]");
  const cancelButton = dialog.querySelector<HTMLButtonElement>("[data-confirm-cancel]");

  if (!title || !message || !confirmButton || !cancelButton) {
    return;
  }

  const close = () => {
    pending = null;
    dialog.classList.remove("confirm-dialog--open");
    dialog.setAttribute("aria-hidden", "true");
  };

  const submitPending = () => {
    if (!pending) {
      return;
    }
    confirmedForm = pending.form;
    const { form, submitter } = pending;
    close();
    if (submitter instanceof HTMLButtonElement || submitter instanceof HTMLInputElement) {
      form.requestSubmit(submitter);
      return;
    }
    form.requestSubmit();
  };

  document.addEventListener(
    "submit",
    (event) => {
      const form = event.target;
      if (!(form instanceof HTMLFormElement)) {
        return;
      }
      if (confirmedForm === form) {
        confirmedForm = null;
        return;
      }
      const confirmationMessage = form.dataset.confirmMessage;
      if (!confirmationMessage) {
        return;
      }

      event.preventDefault();
      pending = {
        form,
        submitter: event instanceof SubmitEvent ? (event.submitter as HTMLElement | null) : null
      };
      title.textContent = form.dataset.confirmTitle || "Confirmar acao";
      message.textContent = confirmationMessage;
      confirmButton.textContent = form.dataset.confirmAccept || "Confirmar";
      cancelButton.textContent = form.dataset.confirmCancel || "Cancelar";
      dialog.classList.add("confirm-dialog--open");
      dialog.setAttribute("aria-hidden", "false");
      window.requestAnimationFrame(() => confirmButton.focus());
    },
    true
  );

  confirmButton.addEventListener("click", submitPending);
  cancelButton.addEventListener("click", close);
  dialog.addEventListener("click", (event) => {
    if (event.target === dialog || event.target === dialog.querySelector(".confirm-dialog__backdrop")) {
      close();
    }
  });
  document.addEventListener("keydown", (event) => {
    if (event.key === "Escape" && dialog.classList.contains("confirm-dialog--open")) {
      close();
    }
  });
}

function ensureDialog() {
  const existing = document.querySelector<HTMLElement>("[data-confirm-dialog]");
  if (existing) {
    return existing;
  }

  const dialog = document.createElement("div");
  dialog.className = "confirm-dialog";
  dialog.setAttribute("data-confirm-dialog", "true");
  dialog.setAttribute("aria-hidden", "true");
  dialog.innerHTML = `
    <div class="confirm-dialog__backdrop"></div>
    <section class="confirm-dialog__panel" role="alertdialog" aria-modal="true" aria-labelledby="actajus-confirm-title">
      <div class="confirm-dialog__content">
        <p class="confirm-dialog__eyebrow">Confirmacao</p>
        <h2 class="confirm-dialog__title" id="actajus-confirm-title" data-confirm-title>Confirmar acao</h2>
        <p class="confirm-dialog__message" data-confirm-message></p>
      </div>
      <div class="confirm-dialog__actions">
        <button class="button button--ghost" type="button" data-confirm-cancel>Cancelar</button>
        <button class="button button--danger" type="button" data-confirm-accept>Confirmar</button>
      </div>
    </section>
  `;
  document.body.appendChild(dialog);
  return dialog;
}
