const passwordRules = {
  length: (value: string) => value.length >= 8,
  lower: (value: string) => /[a-z]/.test(value),
  upper: (value: string) => /[A-Z]/.test(value),
  number: (value: string) => /[0-9]/.test(value),
  special: (value: string) => /[^A-Za-z0-9]/.test(value)
} satisfies Record<string, (value: string) => boolean>;

export function setupRegisterForm() {
  const form = document.querySelector<HTMLFormElement>("[data-register-form]");
  const input = document.querySelector<HTMLInputElement>("[data-password-input]");
  const checklist = document.querySelector<HTMLElement>("[data-password-checklist]");
  if (!form || !input || !checklist) {
    return;
  }

  const updateRules = () => {
    const value = input.value;
    Object.entries(passwordRules).forEach(([rule, matches]) => {
      const item = checklist.querySelector<HTMLElement>(`[data-password-rule="${rule}"]`);
      if (!item) {
        return;
      }
      const valid = matches(value);
      item.dataset.valid = String(valid);
    });
  };

  input.addEventListener("input", updateRules);
  input.addEventListener("focus", updateRules);
  form.addEventListener("submit", updateRules);
  updateRules();
}
