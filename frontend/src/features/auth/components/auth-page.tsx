"use client";

import { useEffect, useRef, useState } from "react";
import { useRouter } from "next/navigation";
import { MessageCircle } from "lucide-react";
import { useAuth } from "@/features/auth/context/auth-context";
import { errorMessage } from "@/lib/api/client";

declare global {
  interface Window {
    google?: {
      accounts: {
        id: {
          initialize: (options: { client_id: string; callback: (result: { credential?: string }) => void }) => void;
          renderButton: (element: HTMLElement, options: Record<string, string>) => void;
        };
      };
    };
  }
}

export function AuthPage() {
  const router = useRouter();
  const googleButton = useRef<HTMLDivElement>(null);
  const { signIn } = useAuth();
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);
  const googleClientId = process.env.NEXT_PUBLIC_GOOGLE_CLIENT_ID;

  useEffect(() => {
    if (!googleClientId) return;
    const setupGoogle = () => {
      if (!window.google || !googleButton.current) return;
      window.google.accounts.id.initialize({
        client_id: googleClientId,
        callback: ({ credential }) => {
          if (!credential) return;
          setError("");
          setLoading(true);
          void signIn("google", credential)
            .then(() => router.replace("/chat"))
            .catch((reason: unknown) => setError(errorMessage(reason, "We could not sign you in. Please try again.")))
            .finally(() => setLoading(false));
        },
      });
      window.google.accounts.id.renderButton(googleButton.current, {
        theme: "outline",
        size: "large",
        width: "360",
        text: "continue_with",
        shape: "pill",
      });
    };
    const script = document.createElement("script");
    script.src = "https://accounts.google.com/gsi/client";
    script.async = true;
    script.onload = setupGoogle;
    document.body.appendChild(script);
    return () => script.remove();
  }, [googleClientId, router, signIn]);

  return <main className="grid min-h-screen place-items-center bg-[linear-gradient(135deg,#eff9ff,#fff_48%,#e8f5ff)] px-5 py-10 text-slate-900"><section className="w-full max-w-md rounded-3xl border border-sky-100 bg-white p-8 shadow-xl shadow-sky-100/70"><div className="mb-8 text-center"><div className="mx-auto mb-4 grid h-14 w-14 place-items-center rounded-2xl bg-sky-500 text-white"><MessageCircle size={28} /></div><h1 className="text-2xl font-bold">Welcome to GoMess</h1><p className="mt-2 text-sm text-slate-500">A private space for your closest conversations.</p></div>{googleClientId ? <div ref={googleButton} className="flex min-h-11 justify-center" /> : <p className="rounded-xl bg-sky-50 p-3 text-sm text-sky-800">Set NEXT_PUBLIC_GOOGLE_CLIENT_ID to enable Google Sign-In.</p>}{loading ? <p className="mt-4 text-center text-sm text-slate-500">Signing you in…</p> : null}{error ? <p role="alert" className="mt-4 rounded-xl bg-red-50 p-3 text-sm text-red-700">{error}</p> : null}<p className="mt-8 text-center text-xs text-slate-400">GoMess uses your Google credential only to authenticate with the GoMess API.</p></section></main>;
}
