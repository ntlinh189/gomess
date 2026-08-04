"use client";

import { createContext, useContext, useEffect, useMemo, useState } from "react";
import { getMe, login, register, type LoginPayload, type RegisterPayload } from "@/lib/api";
import { api } from "@/lib/api/client";

interface AuthContextValue {
  user: any | null;
  token: string | null;
  loading: boolean;
  isAuthenticated: boolean;
  signIn: (payload: LoginPayload) => Promise<void>;
  signUp: (payload: RegisterPayload) => Promise<void>;
  signOut: () => void;
}

const AuthContext = createContext<AuthContextValue | undefined>(undefined);

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<any | null>(null);
  const [token, setToken] = useState<string | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (typeof window === "undefined") return;

    const restoreSession = async () => {
      const storedToken = window.localStorage.getItem("gomess_access_token");
      const storedUser = window.localStorage.getItem("gomess_user");

      if (storedToken) {
        setToken(storedToken);
      }

      if (storedUser) {
        try {
          setUser(JSON.parse(storedUser));
        } catch {
          setUser(null);
        }
      }

      try {
        const response = await api.post("/auth/refresh");
        const nextToken = (response.data as any)?.access_token ?? (response.data as any)?.token ?? null;
        if (nextToken) {
          setToken(nextToken);
          window.localStorage.setItem("gomess_access_token", nextToken);
        }

        const me = await getMe();
        const meUser = (me as any)?.user ?? (me as any);
        if (meUser) {
          setUser(meUser);
          window.localStorage.setItem("gomess_user", JSON.stringify(meUser));
        }
      } catch {
        setToken(null);
        setUser(null);
        window.localStorage.removeItem("gomess_access_token");
        window.localStorage.removeItem("gomess_user");
      } finally {
        setLoading(false);
      }
    };

    restoreSession();
  }, []);

  const signIn = async (payload: LoginPayload) => {
    setLoading(true);
    try {
      const result = await login(payload);
      const nextUser = (result as any)?.user ?? null;
      const nextToken = (result as any)?.access_token ?? (result as any)?.accessToken ?? (result as any)?.token ?? null;

      if (nextToken) {
        setToken(nextToken);
        window.localStorage.setItem("gomess_access_token", nextToken);
      }

      setUser(nextUser);
      window.localStorage.setItem("gomess_user", JSON.stringify(nextUser));

      if (!nextToken) {
        const me = await getMe();
        const meUser = (me as any)?.user ?? (me as any);
        setUser(meUser);
        window.localStorage.setItem("gomess_user", JSON.stringify(meUser));
      }
    } finally {
      setLoading(false);
    }
  };

  const signUp = async (payload: RegisterPayload) => {
    setLoading(true);
    try {
      const result = await register(payload);
      const nextUser = (result as any)?.user ?? null;
      const nextToken = (result as any)?.access_token ?? (result as any)?.accessToken ?? (result as any)?.token ?? null;

      if (nextToken) {
        setToken(nextToken);
        window.localStorage.setItem("gomess_access_token", nextToken);
      }

      setUser(nextUser);
      window.localStorage.setItem("gomess_user", JSON.stringify(nextUser));
    } finally {
      setLoading(false);
    }
  };

  const signOut = () => {
    void (async () => {
      try {
        await api.post("/auth/logout");
      } catch {
        // Ignore logout errors and clear client state.
      } finally {
        setToken(null);
        setUser(null);
        window.localStorage.removeItem("gomess_access_token");
        window.localStorage.removeItem("gomess_user");
      }
    })();
  };

  const value = useMemo<AuthContextValue>(
    () => ({
      user,
      token,
      loading,
      isAuthenticated: Boolean(token),
      signIn,
      signUp,
      signOut,
    }),
    [loading, token, user],
  );

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function useAuth() {
  const ctx = useContext(AuthContext);
  if (!ctx) throw new Error("useAuth must be used within AuthProvider");
  return ctx;
}
