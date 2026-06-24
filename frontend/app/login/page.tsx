"use client"

import { useEffect, useState } from "react";

type CredentialResponse = {
  credential: string;
};

const GOOGLE_SCRIPT_ID = "google-client-script";

type Provider = {
  id: string;
  label: string;
  // type controls how the provider is handled: 'gsi' uses Google Identity Services,
  // 'oauth' would be a redirect/popup flow handled by backend, 'placeholder' is not implemented yet.
  type: "gsi" | "oauth" | "placeholder";
  disabled?: boolean;
};

const PROVIDERS: Provider[] = [
  { id: "google", label: "Google", type: "gsi" },
  { id: "facebook", label: "Facebook", type: "placeholder", disabled: true },
  { id: "tiktok", label: "TikTok", type: "placeholder", disabled: true },
];

export default function Login() {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState(false);
  const googleClientId = process.env.NEXT_PUBLIC_GOOGLE_CLIENT_ID;
  const apiHost = process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080";

  const signInEndpoint = (providerId: string) => `${apiHost.replace(/\/$/, "")}/api/auth/${providerId}`;

  const handleCredentialResponse = async (providerId: string, response: CredentialResponse) => {
    if (!response?.credential) {
      setError("Không nhận được token từ provider. Vui lòng thử lại.");
      return;
    }

    setLoading(true);
    setError(null);

    try {
      const res = await fetch(signInEndpoint(providerId), {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        credentials: "include",
        body: JSON.stringify({ token: response.credential }),
      });

      if (!res.ok) {
        const data = await res.json().catch(() => null);
        throw new Error(data?.error || "Đăng nhập thất bại. Vui lòng thử lại.");
      }

      const data = await res.json();
      window.localStorage.setItem("gomess_access_token", data.accessToken);
      setSuccess(true);
    } catch (err) {
      setError(err instanceof Error ? err.message : String(err));
    } finally {
      setLoading(false);
    }
  };

  // Google-specific initialization
  useEffect(() => {
    if (!PROVIDERS.some((p) => p.id === "google")) return;
    if (!googleClientId) return;

    const existingScript = document.getElementById(GOOGLE_SCRIPT_ID);
    const initializeGoogle = () => {
      const google = (window as any).google;
      if (!google?.accounts?.id) {
        setError("Không thể khởi tạo Google Sign-In.");
        return;
      }

      google.accounts.id.initialize({
        client_id: googleClientId,
        callback: (res: CredentialResponse) => handleCredentialResponse("google", res),
      });

      google.accounts.id.renderButton(document.getElementById("google-signin-button"), {
        theme: "outline",
        size: "large",
        type: "standard",
        text: "signin_with",
        shape: "rectangular",
      });
    };

    if (existingScript) {
      initializeGoogle();
      return;
    }

    const script = document.createElement("script");
    script.id = GOOGLE_SCRIPT_ID;
    script.src = "https://accounts.google.com/gsi/client";
    script.async = true;
    script.defer = true;
    script.onload = initializeGoogle;
    script.onerror = () => setError("Không thể tải Google Sign-In. Vui lòng kiểm tra kết nối mạng.");
    document.body.appendChild(script);

    return () => {
      const scriptNode = document.getElementById(GOOGLE_SCRIPT_ID);
      if (scriptNode) scriptNode.remove();
    };
  }, [googleClientId]);

  const handleProviderClick = (provider: Provider) => {
    setError(null);
    setSuccess(false);
    if (provider.disabled) {
      setError(`${provider.label} chưa được hỗ trợ. Đang chờ cài đặt.`);
      return;
    }

    if (provider.type === "gsi") {
      const google = (window as any).google;
      if (!google?.accounts?.id) {
        setError("Google Sign-In chưa sẵn sàng. Vui lòng tải lại trang.");
        return;
      }
      google.accounts.id.prompt();
      return;
    }

    if (provider.type === "oauth") {
      // Placeholder for future redirect/popup OAuth flows handled by backend.
      // Example: window.location.href = `${apiHost}/api/auth/${provider.id}/start`;
      setError("OAuth flow chưa được cấu hình cho provider này.");
      return;
    }
  };

  return (
    <main className="min-h-screen flex items-center justify-center bg-slate-50 px-4 py-10 text-slate-900">
      <div className="w-full max-w-md rounded-3xl border border-slate-200 bg-white p-8 shadow-xl shadow-slate-200/80">
        <div className="space-y-6 text-center">
          <div>
            <h1 className="text-3xl font-semibold">Đăng nhập</h1>
            <p className="mt-2 text-sm text-slate-600">Chọn phương thức đăng nhập. Hiện tại chỉ hỗ trợ Google.</p>
          </div>

          <div className="flex flex-col gap-3">
            {PROVIDERS.map((p) => (
              <div key={p.id} className="flex items-center justify-center">
                {p.id === "google" ? (
                  <div id="google-signin-button" />
                ) : null}

                <button
                  type="button"
                  onClick={() => handleProviderClick(p)}
                  disabled={p.disabled || loading}
                  className={`w-full max-w-sm rounded-full border px-4 py-3 text-sm font-medium transition ${p.disabled ? "opacity-60 cursor-not-allowed" : "hover:bg-slate-50"}`}
                >
                  {loading && p.id === "google" ? "Đang đăng nhập..." : `Đăng nhập bằng ${p.label}`}
                </button>
              </div>
            ))}
          </div>

          {error ? <p className="rounded-2xl bg-rose-50 px-4 py-3 text-sm text-rose-700">{error}</p> : null}
          {success ? <p className="rounded-2xl bg-emerald-50 px-4 py-3 text-sm text-emerald-700">Đăng nhập thành công! Bạn đã được chuyển hướng.</p> : null}

          <div className="rounded-3xl border border-slate-100 bg-slate-50 px-4 py-4 text-left text-sm text-slate-600">
            <p>Đường dẫn backend (ví dụ):</p>
            <p className="break-all text-slate-500">{signInEndpoint("google")}</p>
          </div>
        </div>
      </div>
    </main>
  );
}
