import { test as base } from "@playwright/test";

export type E2ECredentials = {
  email: string;
  password: string;
};

export function e2eCredentials(): E2ECredentials | null {
  const email = process.env.ACTAJUS_E2E_EMAIL;
  const password = process.env.ACTAJUS_E2E_PASSWORD;
  if (!email || !password) {
    return null;
  }
  return { email, password };
}

export const test = base;
