import type { ToastInput, ToastKind } from "./types";

export function setupToasts() {
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
