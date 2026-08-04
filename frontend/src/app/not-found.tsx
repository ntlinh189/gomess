import Link from "next/link";

export default function NotFound() {
  return (
    <main className="flex min-h-screen items-center justify-center bg-slate-950 px-6">
      <div className="text-center">
        <h2 className="text-2xl font-semibold">Page not found</h2>
        <Link href="/" className="mt-4 inline-block text-slate-300 underline">
          Go home
        </Link>
      </div>
    </main>
  );
}
