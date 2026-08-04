"use client";

import Link from "next/link";
import { Button } from "@/components/ui/button";
import { useAuth } from "@/features/auth/context/auth-context";

export function AppShell({ children }: { children: React.ReactNode }) {
  const { user, signOut } = useAuth();

  const avatarUrl = user?.avatar;
  const initial = (user?.name || user?.email || "U").charAt(0).toUpperCase();

  return (
    <div className="min-h-screen bg-[radial-gradient(circle_at_top,_#f8fbff,_#eef5ff_40%,_#e6efff)]">
      <header className="border-b border-sky-100 bg-white/80 backdrop-blur">
        <div className="mx-auto flex max-w-7xl items-center justify-between px-6 py-4">
          <Link href="/dashboard" className="text-lg font-semibold text-slate-800">
            GoMess
          </Link>
          <div className="flex items-center gap-3">
            <div className="flex items-center gap-2 rounded-full border border-sky-100 bg-sky-50 px-3 py-2">
              {avatarUrl ? (
                <img src={avatarUrl} alt={user?.name || "User"} className="h-8 w-8 rounded-full object-cover" />
              ) : (
                <div className="flex h-8 w-8 items-center justify-center rounded-full bg-sky-600 text-sm font-semibold text-white">
                  {initial}
                </div>
              )}
              <span className="text-sm font-medium text-slate-700">{user?.name || user?.email || "User"}</span>
            </div>
            <Button variant="outline" size="sm" onClick={() => signOut()}>
              Logout
            </Button>
          </div>
        </div>
      </header>
      <main className="mx-auto max-w-7xl p-4 sm:p-6">{children}</main>
    </div>
  );
}
