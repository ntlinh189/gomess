"use client";

import { useEffect } from "react";
import { useRouter } from "next/navigation";
import { useAuth } from "@/features/auth/context/auth-context";
import { ChatShell } from "@/features/chat/components/chat-shell";

export default function DashboardPage() {
  const router = useRouter();
  const { isAuthenticated, loading } = useAuth();

  useEffect(() => {
    if (!loading && !isAuthenticated) {
      router.replace("/login");
    }
  }, [isAuthenticated, loading, router]);

  if (loading) {
    return (
      <main className="flex min-h-screen items-center justify-center px-6 py-10 text-slate-600">
        Checking your session...
      </main>
    );
  }

  if (!isAuthenticated) {
    return null;
  }

  return (
    <main className="h-dvh p-3 md:p-4">
      <ChatShell />
    </main>
  );
}
