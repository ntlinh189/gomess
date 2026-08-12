"use client";
/* eslint-disable @next/next/no-img-element */

import type { Friend } from "@/types/user";

export function ConversationList({ friends, selectedId, onSelect }: { friends: Friend[]; selectedId: number | null; onSelect: (friend: Friend) => void }) {
  if (!friends.length) return <p className="px-3 py-8 text-center text-sm text-slate-500">No friends yet. Search above to connect.</p>;
  return <div className="space-y-1">{friends.map((friend) => <button key={friend.id} onClick={() => onSelect(friend)} className={`flex w-full items-center gap-3 rounded-xl px-3 py-3 text-left transition ${selectedId === friend.id ? "bg-sky-100 text-sky-950" : "hover:bg-slate-100"}`}><Avatar name={friend.name} src={friend.avatar} /><span className="min-w-0"><span className="block truncate font-medium">{friend.name}</span><span className="block truncate text-xs text-slate-500">{friend.account}</span></span></button>)}</div>;
}

export function Avatar({ name, src, className = "" }: { name: string; src?: string; className?: string }) {
  if (src) return <img src={src} alt="" className={`h-10 w-10 rounded-full object-cover ${className}`} />;
  return <span aria-hidden className={`grid h-10 w-10 shrink-0 place-items-center rounded-full bg-sky-200 font-semibold text-sky-800 ${className}`}>{name.slice(0, 1).toUpperCase()}</span>;
}
