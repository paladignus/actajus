export type ToastKind = "success" | "error" | "info" | "alert";

export type ToastInput = {
  kind?: ToastKind;
  title?: string;
  message: string;
  source?: string;
  code?: string;
  timeoutMs?: number;
};

export type RouteEntry = {
  href: string;
  label: string;
  keywords: string[];
};

declare global {
  interface Window {
    __ACTAJUS_BOOTSTRAP__?: {
      appName?: string;
      apiBaseURL?: string;
      authenticated?: boolean;
      canManageSessions?: boolean;
      canManageRBAC?: boolean;
    };
    actajusNotify?: (toast: ToastInput) => void;
  }
}

export {};
