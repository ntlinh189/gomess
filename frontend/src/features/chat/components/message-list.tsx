"use client";
/* eslint-disable @next/next/no-img-element */

import { FileIcon, MoreHorizontal } from "lucide-react";
import type { Message } from "@/types/chat";

export function MessageList({ messages, currentUserId, onDelete, onRevoke }: { messages: Message[]; currentUserId: number; onDelete: (id: number) => void; onRevoke: (id: number) => void }) {
  if (!messages.length) return <div className="grid h-full place-items-center text-center text-sm text-slate-500">No messages yet.<br />Say hello to get the conversation started.</div>;
  return <div className="space-y-3">{messages.map((message) => <MessageBubble key={message.id} message={message} mine={message.sender_id === currentUserId} onDelete={onDelete} onRevoke={onRevoke} />)}</div>;
}

function MessageBubble({ message, mine, onDelete, onRevoke }: { message: Message; mine: boolean; onDelete: (id: number) => void; onRevoke: (id: number) => void }) {
  return <article className={`group flex ${mine ? "justify-end" : "justify-start"}`}><div className={`max-w-[85%] rounded-2xl px-4 py-2.5 text-sm ${mine ? "bg-sky-600 text-white" : "bg-slate-100 text-slate-800"}`}><div className="flex items-start gap-2"><div className="min-w-0 flex-1">{message.revoked ? <em>This message was revoked.</em> : <><p className="whitespace-pre-wrap break-words">{message.content}</p>{message.attachments.map((attachment) => attachment.type === "image" ? <a key={attachment.id} href={attachment.url} target="_blank" rel="noreferrer"><img src={attachment.url} alt={attachment.file_name} className="mt-2 max-h-64 rounded-lg object-cover" /></a> : <a key={attachment.id} href={attachment.url} target="_blank" rel="noreferrer" className={`mt-2 flex items-center gap-2 rounded-lg p-2 ${mine ? "bg-sky-700" : "bg-white"}`}><FileIcon size={16} />{attachment.file_name}</a>)}</>}<time className={`mt-1 block text-[11px] ${mine ? "text-sky-100" : "text-slate-400"}`}>{new Date(message.created_at).toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" })}</time></div>{!message.revoked ? <details className="relative"><summary aria-label="Message actions" className="cursor-pointer list-none opacity-0 transition group-hover:opacity-100"><MoreHorizontal size={16} /></summary><div className="absolute right-0 z-10 mt-1 w-32 rounded-lg bg-white p-1 text-slate-800 shadow-lg"><button onClick={() => onDelete(message.id)} className="w-full rounded px-2 py-1 text-left hover:bg-slate-100">Delete for me</button>{mine ? <button onClick={() => onRevoke(message.id)} className="w-full rounded px-2 py-1 text-left hover:bg-slate-100">Revoke</button> : null}</div></details> : null}</div></div></article>;
}
