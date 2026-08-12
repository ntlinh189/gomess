import axios, { AxiosError, type InternalAxiosRequestConfig } from "axios";

const rawApiUrl = process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080";
const apiBaseUrl = rawApiUrl.replace(/\/$/, "").endsWith("/api")
  ? rawApiUrl.replace(/\/$/, "")
  : `${rawApiUrl.replace(/\/$/, "")}/api`;

export const api = axios.create({ baseURL: apiBaseUrl, withCredentials: true });

let accessToken: string | null = null;
let refreshPromise: Promise<string | null> | null = null;

export function setAccessToken(token: string | null) {
  accessToken = token;
}

export function getApiBaseUrl() {
  return apiBaseUrl;
}

api.interceptors.request.use((config: InternalAxiosRequestConfig) => {
  if (accessToken) config.headers.Authorization = `Bearer ${accessToken}`;
  return config;
});

async function refreshAccessToken() {
  if (!refreshPromise) {
    refreshPromise = api
      .post<{ access_token: string }>("/auth/refresh", undefined, { headers: { Authorization: undefined } })
      .then(({ data }) => {
        setAccessToken(data.access_token);
        return data.access_token;
      })
      .catch(() => {
        setAccessToken(null);
        return null;
      })
      .finally(() => {
        refreshPromise = null;
      });
  }
  return refreshPromise;
}

api.interceptors.response.use(
  (response) => response,
  async (error: AxiosError) => {
    const request = error.config as (InternalAxiosRequestConfig & { _retried?: boolean }) | undefined;
    if (!request || error.response?.status !== 401 || request._retried || request.url?.includes("/auth/refresh")) {
      return Promise.reject(error);
    }
    request._retried = true;
    const token = await refreshAccessToken();
    if (!token) return Promise.reject(error);
    request.headers.Authorization = `Bearer ${token}`;
    return api(request);
  },
);

export function errorMessage(error: unknown, fallback: string) {
  if (axios.isAxiosError<{ error?: string }>(error)) return error.response?.data?.error ?? fallback;
  return fallback;
}
