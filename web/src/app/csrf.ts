const CSRF_COOKIE = "actajus_csrf";
const CSRF_FIELD = "_csrf";

function readCookie(name: string): string {
  const prefix = `${name}=`;
  const match = document.cookie
    .split(";")
    .map((item) => item.trim())
    .find((item) => item.startsWith(prefix));

  return match ? decodeURIComponent(match.slice(prefix.length)) : "";
}

function ensureHiddenField(form: HTMLFormElement, token: string) {
  if (!token) {
    return;
  }
  let input = form.querySelector<HTMLInputElement>(`input[name="${CSRF_FIELD}"]`);
  if (!input) {
    input = document.createElement("input");
    input.type = "hidden";
    input.name = CSRF_FIELD;
    form.append(input);
  }
  input.value = token;
}

export function ensurePostFormsCSRF() {
  const token = readCookie(CSRF_COOKIE);
  if (!token) {
    return;
  }

  const forms = document.querySelectorAll<HTMLFormElement>('form[method="post"]');
  forms.forEach((form) => ensureHiddenField(form, token));
}
