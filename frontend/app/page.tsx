"use client";

import { useEffect } from "react";
import { useRouter } from "next/navigation";
import { useAuth } from "@/features/auth/context/auth-context";
import { AuthPage } from "@/features/auth/components/auth-page";

export default function HomePage() {
  const router = useRouter();
  const { isAuthenticated, loading } = useAuth();

  useEffect(() => {
    if (!loading && isAuthenticated) {
      router.replace("/chat");
    }
  }, [isAuthenticated, loading, router]);

  if (loading) {
    return (
      <main className="flex min-h-screen items-center justify-center px-6 py-10 text-slate-600">
        Checking your session...
      </main>
    );
  }

  if (isAuthenticated) {
    return null;
  }

  return <AuthPage />;
}
