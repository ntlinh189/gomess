"use client";

import { useMemo } from "react";
import { useQuery } from "@tanstack/react-query";
import { listConversations } from "@/lib/api";
import { cn } from "@/lib/utils";

export function ConversationList({ value, onSelect }: { value?: string; onSelect: (id: string) => void }) {
  const { data } = useQuery({
    queryKey: ["conversations"],
    queryFn: listConversations,
  });

  const conversations = useMemo(() => {
    return Array.isArray(data) ? data : (data as any)?.conversations || (data as any)?.items || [];
  }, [data]);

  return (
    <div className="space-y-2">
      {conversations.map((conversation: any, index: number) => {
        const id = conversation.id || conversation._id || conversation.conversationId || `${index}`;
        const name = conversation.name || `Conversation ${index + 1}`;

        return (
          <button
            key={id}
            className={cn(
              "flex w-full items-center gap-3 rounded-2xl border border-sky-100 bg-white/90 p-3 text-left shadow-sm transition",
              value === id && "border-sky-400 bg-sky-50 shadow-md",
            )}
            onClick={() => onSelect(id)}
          >
            <div className="flex h-10 w-10 items-center justify-center rounded-full bg-sky-100 text-sm font-semibold text-sky-700">
              {(name || "F").charAt(0).toUpperCase()}
            </div>
            <div className="min-w-0 flex-1">
              <div className="truncate font-medium text-slate-800">{name}</div>
              <div className="truncate text-sm text-slate-500">{conversation.lastMessage || "Start a conversation"}</div>
            </div>
          </button>
        );
      })}
    </div>
  );
}
