"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { useAuth } from "@/features/auth/context/auth-context";

declare global {
  interface Window {
    google?: {
      accounts: {
        id: {
          initialize: (config: { client_id: string; callback: (response: { credential?: string }) => void }) => void;
          renderButton: (element: HTMLElement | null, options: { theme?: string; size?: string; text?: string; shape?: string }) => void;
        };
      };
    };
  }
}

export function AuthPage({ mode }: { mode: "login" | "register" }) {
  const router = useRouter();
  const { signIn, signUp } = useAuth();
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    const clientId = process.env.NEXT_PUBLIC_GOOGLE_CLIENT_ID;
    if (!clientId || typeof window === "undefined") {
      setError("Google client ID is not configured.");
      return;
    }

    const loadGoogleScript = () => {
      const existingScript = document.getElementById("google-gsi-script");
      if (existingScript) {
        if (window.google?.accounts?.id) {
          initializeGoogle();
        } else {
          existingScript.addEventListener("load", initializeGoogle, { once: true });
        }
        return;
      }

      const script = document.createElement("script");
      script.id = "google-gsi-script";
      script.src = "https://accounts.google.com/gsi/client";
      script.async = true;
      script.defer = true;
      script.onload = initializeGoogle;
      document.body.appendChild(script);
    };

    const initializeGoogle = () => {
      if (!window.google?.accounts?.id) return;

      window.google.accounts.id.initialize({
        client_id: clientId,
        callback: async (response: { credential?: string }) => {
          if (!response.credential) {
            setError("Google did not return a credential.");
            return;
          }

          setError("");
          setLoading(true);
          try {
            if (mode === "login") {
              await signIn({ provider: "google", token: response.credential });
            } else {
              await signUp({ provider: "google", token: response.credential });
            }
            router.replace("/chat");
          } catch (err: any) {
            setError(err?.response?.data?.error || err?.message || "Google authentication failed");
          } finally {
            setLoading(false);
          }
        },
      });

      const container = document.getElementById("google-signin-button");
      if (container) {
        window.google.accounts.id.renderButton(container, {
          theme: "filled_blue",
          size: "large",
          text: mode === "login" ? "signin_with" : "signup_with",
          shape: "pill",
        });
      }
    };

    loadGoogleScript();
  }, [mode, router, signIn, signUp]);

  return (
    <main className="flex min-h-screen items-center justify-center px-6 py-10">
      <Card className="w-full max-w-md">
        <CardHeader>
          <CardTitle>{mode === "login" ? "Sign in with Google" : "Sign up with Google"}</CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          <p className="text-sm text-slate-500">
            {mode === "login"
              ? "Continue securely with your Google account."
              : "Create your account quickly with Google."}
          </p>
          <div id="google-signin-button" className="flex justify-center" />
          {loading ? <p className="text-center text-sm text-slate-600">Signing you in...</p> : null}
          {error ? <p className="text-sm text-red-500">{error}</p> : null}
        </CardContent>
      </Card>
    </main>
  );
}
