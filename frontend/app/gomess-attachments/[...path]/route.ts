import type { NextRequest } from "next/server";
import http from "node:http";

const minioOrigin = "http://localhost:9000";

async function proxyToMinio(request: NextRequest, path: string[]) {
  const target = new URL(`/gomess-attachments/${path.map(encodeURIComponent).join("/")}`, minioOrigin);
  target.search = request.nextUrl.search;

  const headers = new Headers(request.headers);
  const signedHost = request.headers.get("x-forwarded-host") ?? request.headers.get("host");
  if (signedHost) headers.set("host", signedHost);
  headers.delete("connection");
  headers.delete("content-length");

  const body = request.method === "GET" || request.method === "HEAD"
    ? undefined
    : Buffer.from(await request.arrayBuffer());
  if (body) headers.set("content-length", String(body.byteLength));

  return new Promise<Response>((resolve, reject) => {
    const proxyRequest = http.request(target, {
      method: request.method,
      headers: Object.fromEntries(headers.entries()),
    }, (proxyResponse) => {
      const chunks: Buffer[] = [];
      proxyResponse.on("data", (chunk: Buffer) => chunks.push(chunk));
      proxyResponse.on("end", () => {
        const responseHeaders = new Headers();
        for (const [name, value] of Object.entries(proxyResponse.headers)) {
          if (value !== undefined) responseHeaders.set(name, Array.isArray(value) ? value.join(", ") : value);
        }
        responseHeaders.delete("transfer-encoding");
        resolve(new Response(Buffer.concat(chunks), {
          status: proxyResponse.statusCode ?? 502,
          headers: responseHeaders,
        }));
      });
    });
    proxyRequest.on("error", reject);
    if (body) proxyRequest.write(body);
    proxyRequest.end();
  });
}

type RouteContext = { params: Promise<{ path: string[] }> };

export async function GET(request: NextRequest, context: RouteContext) {
  return proxyToMinio(request, (await context.params).path);
}

export async function HEAD(request: NextRequest, context: RouteContext) {
  return proxyToMinio(request, (await context.params).path);
}

export async function PUT(request: NextRequest, context: RouteContext) {
  return proxyToMinio(request, (await context.params).path);
}
