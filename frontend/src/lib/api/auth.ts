import { api } from "./client";
import type { OAuthProvider, CurrentUser } from "@/types/user";

export async function login(provider: OAuthProvider, token: string) {
  const { data } = await api.post<{ access_token: string }>(`/auth/${provider}`, { token });
  return data;
}

export async function getMe() {
  const { data } = await api.get<CurrentUser>("/user/me");
  return data;
}

export async function logout() {
  await api.post("/auth/logout");
}

export async function deleteMe() {
  await api.delete("/user/me", { data: { confirm: true } });
}
