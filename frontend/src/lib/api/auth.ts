import { api, unwrapResponse } from "./client";

export interface LoginPayload {
  provider?: string;
  token?: string;
  email?: string;
  password?: string;
}

export interface RegisterPayload {
  provider?: string;
  token?: string;
  name?: string;
  email?: string;
  password?: string;
}

function buildAuthPayload(payload: LoginPayload | RegisterPayload) {
  const provider = payload.provider || "google";
  const token = payload.token || payload.password || payload.email || "";

  return {
    token,
  };
}

export async function login(payload: LoginPayload) {
  const response = await api.post(`/auth/${payload.provider || "google"}`, buildAuthPayload(payload));
  return unwrapResponse(response);
}

export async function register(payload: RegisterPayload) {
  const response = await api.post(`/auth/${payload.provider || "google"}`, buildAuthPayload(payload));
  return unwrapResponse(response);
}

export async function getMe() {
  const response = await api.get("/user/me");
  return unwrapResponse(response);
}
