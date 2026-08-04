"use client";

import { useQuery } from "@tanstack/react-query";
import { listMessages } from "@/lib/api";

export function MessageList({ conversationId }: { conversationId?: string }) {
  const { data } = useQuery({
    queryKey: ["messages", conversationId],
    queryFn: () => listMessages(conversationId!),
    enabled: Boolean(conversationId),
  });

  const messages = Array.isArray(data) ? data : (data as any)?.messages || (data as any)?.items || [];

  if (!conversationId) {
    return (
      <div className="flex h-full items-center justify-center rounded-2xl border border-dashed border-sky-200 bg-sky-50/70 text-sm text-slate-500">
        Select a conversation to start chatting.
      </div>
    );
  }

  return (
    <div className="flex h-full flex-col gap-3 overflow-auto rounded-2xl border border-sky-100 bg-sky-50/70 p-4">
      {messages.map((message: any, index: number) => (
        <div key={message.id || message._id || `${message.text || "msg"}-${index}`} className="max-w-[80%] rounded-2xl bg-white p-3 shadow-sm">
          <div className="text-sm font-medium text-sky-700">{message.sender?.name || message.user?.name || "You"}</div>
          <div className="mt-1 text-sm text-slate-600">{message.text || message.content || message.body || "..."}</div>
        </div>
      ))}
    </div>
  );
}
