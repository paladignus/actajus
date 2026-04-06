export async function loadBootstrap() {
  const response = await fetch("/bootstrap", {
    headers: { Accept: "application/json" }
  });
  if (!response.ok) {
    return null;
  }
  return response.json() as Promise<{ authenticated?: boolean }>;
}
