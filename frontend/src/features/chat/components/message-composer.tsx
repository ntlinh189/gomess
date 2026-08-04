"use client";

import { useState } from "react";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { sendMessage } from "@/lib/api";

export function MessageComposer({ conversationId }: { conversationId?: string }) {
  const [text, setText] = useState("");
  const queryClient = useQueryClient();

  const mutation = useMutation({
    mutationFn: ({ conversationId, text }: { conversationId: string; text: string }) => sendMessage(conversationId, text),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["messages", conversationId] });
      setText("");
    },
  });

  const handleSubmit = (event: React.FormEvent) => {
    event.preventDefault();
    if (!conversationId || !text.trim()) return;
    mutation.mutate({ conversationId, text: text.trim() });
  };

  return (
    <form onSubmit={handleSubmit} className="flex gap-2">
      <Input value={text} onChange={(e) => setText(e.target.value)} placeholder="Type a message..." className="flex-1 rounded-full border-sky-200" />
      <Button type="submit" variant="default" className="rounded-full px-4">
        Send
      </Button>
    </form>
  );
}
