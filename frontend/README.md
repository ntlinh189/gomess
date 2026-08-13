# GoMess frontend

Next.js client for the GoMess API. Copy `.env.example` to `.env.local`, configure your Google client ID, then run:

```bash
npm install
npm run dev
```

Browser API requests use the same-origin `/api` path. Next.js proxies them to the local API at `http://localhost:8080/api`.
