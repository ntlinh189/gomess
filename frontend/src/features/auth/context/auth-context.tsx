"use client";

import { createContext, useCallback, useContext, useEffect, useMemo, useState } from "react";
import { getMe, login, logout } from "@/lib/api/auth";
import { setAccessToken } from "@/lib/api/client";
import type { CurrentUser, OAuthProvider } from "@/types/user";

interface AuthContextValue {
  user: CurrentUser | null;
  accessToken: string | null;
  loading: boolean;
  isAuthenticated: boolean;
  signIn: (provider: OAuthProvider, providerToken: string) => Promise<void>;
  signOut: () => Promise<void>;
  deleteLocalSession: () => void;
}

const AuthContext = createContext<AuthContextValue | null>(null);

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<CurrentUser | null>(null);
  const [accessToken, setToken] = useState<string | null>(null);
  const [loading, setLoading] = useState(true);

  const applyToken = useCallback((token: string | null) => {
    setAccessToken(token);
    setToken(token);
  }, []);

  const deleteLocalSession = useCallback(() => {
    applyToken(null);
    setUser(null);
  }, [applyToken]);

  useEffect(() => {
    let active = true;
    async function restoreSession() {
      try {
        // The refresh token is HTTP-only; only the browser sends it with this request.
        const { data } = await import("@/lib/api/client").then(({ api }) => api.post<{ access_token: string }>("/auth/refresh"));
        if (!active) return;
        applyToken(data.access_token);
        const currentUser = await getMe();
        if (active) setUser(currentUser);
      } catch {
        if (active) deleteLocalSession();
      } finally {
        if (active) setLoading(false);
      }
    }
    void restoreSession();
    return () => { active = false; };
  }, [applyToken, deleteLocalSession]);

  const signIn = useCallback(async (provider: OAuthProvider, providerToken: string) => {
    const { access_token } = await login(provider, providerToken);
    applyToken(access_token);
    try {
      setUser(await getMe());
    } catch (error) {
      deleteLocalSession();
      throw error;
    }
  }, [applyToken, deleteLocalSession]);

  const signOut = useCallback(async () => {
    try { await logout(); } finally { deleteLocalSession(); }
  }, [deleteLocalSession]);

  const value = useMemo(() => ({ user, accessToken, loading, isAuthenticated: Boolean(user && accessToken), signIn, signOut, deleteLocalSession }), [accessToken, deleteLocalSession, loading, signIn, signOut, user]);
  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function useAuth() {
  const value = useContext(AuthContext);
  if (!value) throw new Error("useAuth must be used within AuthProvider");
  return value;
}
