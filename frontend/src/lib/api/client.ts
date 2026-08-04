import axios from "axios";

const rawBaseURL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";
const normalizedBaseURL = rawBaseURL.endsWith("/api")
  ? rawBaseURL
  : `${rawBaseURL.replace(/\/$/, "")}/api`;

export const api = axios.create({
  baseURL: normalizedBaseURL,
  withCredentials: true,
});

api.interceptors.request.use((config) => {
  if (typeof window !== "undefined") {
    const token = window.localStorage.getItem("gomess_access_token");
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
  }
  return config;
});

api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      window.localStorage.removeItem("gomess_access_token");
      window.localStorage.removeItem("gomess_user");
    }
    return Promise.reject(error);
  },
);

export function unwrapResponse<T>(response: any): T {
  return response?.data?.data ?? response?.data?.result ?? response?.data ?? response;
}
