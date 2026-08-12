"use client";

import { useEffect, useMemo, useState } from "react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { LogOut, Search, Settings, UserPlus, Users, X } from "lucide-react";
import { useAuth } from "@/features/auth/context/auth-context";
import { acceptFriendRequest, deleteMessage, listFriends, listMessages, listReceivedRequests, presignUpload, rejectFriendRequest, removeFriend, revokeMessage, searchUsers, sendFriendRequest, sendMessage, uploadFile } from "@/lib/api";
import { deleteMe } from "@/lib/api/auth";
import { errorMessage, getApiBaseUrl } from "@/lib/api/client";
import type { Message, WebSocketEvent } from "@/types/chat";
import type { CurrentUser, Friend, FriendRequest, UserSearchResult } from "@/types/user";
import { Avatar, ConversationList } from "./conversation-list";
import { MessageComposer } from "./message-composer";
import { MessageList } from "./message-list";

type Panel = "empty" | "search" | "requests" | "chat";

export function ChatShell() {
  const queryClient = useQueryClient();
  const { user, accessToken, signOut } = useAuth();
  const [panel, setPanel] = useState<Panel>("empty");
  const [selected, setSelected] = useState<Friend | null>(null);
  const [keyword, setKeyword] = useState("");
  const [results, setResults] = useState<UserSearchResult[]>([]);
  const [notice, setNotice] = useState("");
  const [showSettings, setShowSettings] = useState(false);
  const friendsQuery = useQuery({ queryKey: ["friends", user?.id], queryFn: listFriends, enabled: Boolean(user?.id && accessToken), refetchOnMount: "always", staleTime: 0 });
  const requestsQuery = useQuery({ queryKey: ["received-requests"], queryFn: listReceivedRequests, enabled: Boolean(accessToken) });
  const historyQuery = useQuery({ queryKey: ["messages", selected?.id], queryFn: () => listMessages(selected!.id), enabled: panel === "chat" && Boolean(selected) });
  const friends = friendsQuery.data ?? [];
  const requests = requestsQuery.data ?? [];
  const messages = useMemo(() => [...(historyQuery.data ?? [])].reverse(), [historyQuery.data]);

  useEffect(() => {
    if (!accessToken) return;
    const url = new URL(getApiBaseUrl());
    url.protocol = url.protocol === "https:" ? "wss:" : "ws:";
    url.pathname = `${url.pathname.replace(/\/api$/, "")}/api/ws`;
    url.searchParams.set("token", accessToken);
    const socket = new WebSocket(url);
    socket.onmessage = ({ data }) => {
      try {
        const event = JSON.parse(data) as WebSocketEvent;
        if (event.event !== "message.new") return;
        const message = event.data;
        const friendID = message.sender_id === user?.id ? message.receiver_id : message.sender_id;
        queryClient.setQueryData<Message[]>(["messages", friendID], (old = []) => old.some((item) => item.id === message.id) ? old : [message, ...old]);
      } catch { /* Ignore invalid socket messages. */ }
    };
    return () => socket.close();
  }, [accessToken, queryClient, user?.id]);

  const addRequest = useMutation({ mutationFn: sendFriendRequest, onSuccess: () => setNotice("Friend request sent."), onError: (error) => setNotice(errorMessage(error, "Could not send friend request.")) });
  async function refreshSocial() { await Promise.all([queryClient.invalidateQueries({ queryKey: ["friends"] }), queryClient.invalidateQueries({ queryKey: ["received-requests"] })]); }
  async function search() { if (!keyword.trim()) return; try { setResults(await searchUsers("google", keyword.trim(), 0, 20)); setPanel("search"); } catch (error) { setNotice(errorMessage(error, "Search failed.")); } }
  async function requestAction(task: () => Promise<void>, success: string) { try { await task(); await refreshSocial(); setNotice(success); } catch (error) { setNotice(errorMessage(error, "That action could not be completed.")); } }
  async function send(content: string, files: File[]) { if (!selected) return; try { const attachments = await Promise.all(files.map(async (file) => { const upload = await presignUpload(file.name); await uploadFile(upload.upload_url, file); return { object_key: upload.object_key, file_name: file.name }; })); const message = await sendMessage(selected.id, content, attachments); queryClient.setQueryData<Message[]>(["messages", selected.id], (old = []) => old.some((item) => item.id === message.id) ? old : [message, ...old]); } catch (error) { setNotice(errorMessage(error, "Your message could not be sent.")); throw error; } }
  function chooseFriend(friend: Friend) { setSelected(friend); setPanel("chat"); }

  return <div className="grid h-full w-full overflow-hidden rounded-3xl border border-sky-100 bg-white shadow-xl shadow-sky-100/60 lg:grid-cols-[340px_minmax(0,1fr)]"><aside className="flex min-h-0 flex-col border-r border-sky-100 bg-white"><header className="flex items-center justify-between p-4"><div className="flex min-w-0 items-center gap-3"><Avatar name={user?.name ?? ""} src={user?.avatar} className="h-11 w-11" /><div className="min-w-0"><p className="truncate font-semibold text-slate-900">{user?.name}</p><p className="truncate text-sm text-slate-500">{user?.account}</p></div></div><button onClick={() => setShowSettings(true)} aria-label="Settings" className="rounded-full p-2 text-sky-800 hover:bg-sky-50"><Settings size={19} /></button></header><div className="px-3"><div className="flex items-center gap-2 rounded-xl bg-sky-50 p-1"><input value={keyword} onChange={(event) => setKeyword(event.target.value)} onKeyDown={(event) => event.key === "Enter" && void search()} placeholder="Search users" className="min-w-0 flex-1 bg-transparent px-2 py-1.5 text-sm outline-none placeholder:text-slate-400" /><button onClick={() => void search()} aria-label="Search users" className="rounded-lg bg-white p-2 text-sky-800 shadow-sm"><Search size={16} /></button></div></div><button onClick={() => setPanel("requests")} className={`mx-3 mt-3 flex items-center justify-between rounded-xl px-3 py-2.5 text-left text-sm font-semibold ${panel === "requests" ? "bg-sky-100 text-sky-950" : "bg-sky-50 text-sky-800 hover:bg-sky-100"}`}><span className="flex items-center gap-2"><Users size={17} /> Requests</span><span className="grid h-5 min-w-5 place-items-center rounded-full bg-sky-600 px-1 text-xs text-white">{requests.length}</span></button><div className="mt-4 min-h-0 flex-1 overflow-y-auto px-3 pb-4"><p className="mb-2 px-2 text-xs font-semibold uppercase tracking-wide text-sky-700">Friends</p>{friendsQuery.isLoading ? <p className="p-3 text-sm text-slate-500">Loading friends…</p> : <ConversationList friends={friends} selectedId={selected?.id ?? null} onSelect={chooseFriend} />}</div></aside><main className="flex min-h-0 flex-col bg-sky-50/40">{panel === "chat" && selected ? <ChatPanel friend={selected} messages={messages} loading={historyQuery.isLoading} currentUserID={user!.id} onDelete={(id) => void requestAction(() => deleteMessage(id), "Message deleted.")} onRevoke={(id) => void requestAction(() => revokeMessage(id), "Message revoked.")} onRemove={() => { if (confirm(`Remove ${selected.name}? This also removes your conversation.`)) void requestAction(() => removeFriend(selected.id), "Friend removed."); }} onSend={send} /> : panel === "search" ? <SearchPanel results={results} keyword={keyword} pending={addRequest.isPending} onAdd={(id) => addRequest.mutate(id)} /> : panel === "requests" ? <RequestsPanel requests={requests} onAccept={(id) => void requestAction(() => acceptFriendRequest(id), "Friend request accepted.")} onReject={(id) => void requestAction(() => rejectFriendRequest(id), "Friend request rejected.")} /> : <EmptyPanel />}</main>{notice ? <div role="status" className="fixed bottom-6 left-1/2 z-30 -translate-x-1/2 rounded-full bg-sky-900 px-4 py-2 text-sm text-white shadow-lg">{notice}<button onClick={() => setNotice("")} aria-label="Dismiss" className="ml-3"><X size={15} /></button></div> : null}{showSettings ? <SettingsDialog user={user} onClose={() => setShowSettings(false)} onSignOut={() => void signOut()} onDelete={() => { if (confirm("Permanently delete your GoMess account?")) void deleteMe().then(signOut); }} /> : null}</div>;
}

function EmptyPanel() { return <div className="grid h-full place-items-center p-6 text-center"><div><div className="mx-auto mb-3 grid h-16 w-16 place-items-center rounded-2xl bg-sky-100 text-sky-700"><Users /></div><h2 className="font-semibold">Your messages</h2><p className="mt-1 text-sm text-slate-500">Choose a friend, search users, or review requests.</p></div></div>; }
function SearchPanel({ results, keyword, pending, onAdd }: { results: UserSearchResult[]; keyword: string; pending: boolean; onAdd: (id: number) => void }) { return <section className="h-full overflow-y-auto p-6"><h2 className="text-xl font-semibold">Search results</h2><p className="mt-1 text-sm text-slate-500">{results.length ? `People matching “${keyword}”` : "No users found."}</p><div className="mt-5 grid gap-3 sm:grid-cols-2 xl:grid-cols-3">{results.map((person) => <article key={person.id} className="flex items-center gap-3 rounded-2xl border border-sky-100 bg-white p-4"><Avatar name={person.name} src={person.avatar} /><div className="min-w-0 flex-1"><p className="truncate font-medium">{person.name}</p><p className="truncate text-sm text-slate-500">{person.account}</p></div>{person.request_status === "pending" ? <button type="button" disabled title="Already sent request" aria-label="Already sent request" className="cursor-not-allowed rounded-lg bg-slate-300 p-2 text-slate-500"><UserPlus size={17} /></button> : <button disabled={pending} onClick={() => onAdd(person.id)} className="rounded-lg bg-sky-600 p-2 text-white disabled:opacity-50"><UserPlus size={17} /></button>}</article>)}</div></section>; }
function RequestsPanel({ requests, onAccept, onReject }: { requests: FriendRequest[]; onAccept: (id: number) => void; onReject: (id: number) => void }) { return <section className="h-full overflow-y-auto p-6"><h2 className="text-xl font-semibold">Friend requests</h2><p className="mt-1 text-sm text-slate-500">{requests.length ? `${requests.length} pending request${requests.length === 1 ? "" : "s"}` : "You have no pending requests."}</p><div className="mt-5 space-y-3">{requests.map((request) => { const sender = request.sender; return <article key={request.id} className="flex items-center gap-3 rounded-2xl border border-sky-100 bg-white p-4"><Avatar name={sender?.name ?? "Unknown sender"} src={sender?.avatar} className="h-12 w-12" /><div className="min-w-0 flex-1"><p className="truncate font-medium">{sender?.name ?? "Sender details unavailable"}</p><p className="truncate text-sm text-slate-500">{sender ? `${sender.provider} · ${sender.account}` : "Restart the API to load sender details"}</p></div><button onClick={() => onAccept(request.id)} className="rounded-lg bg-sky-600 px-3 py-2 text-sm text-white">Accept</button><button onClick={() => onReject(request.id)} className="rounded-lg bg-sky-50 px-3 py-2 text-sm text-sky-900">Decline</button></article>; })}</div></section>; }
function ChatPanel({ friend, messages, loading, currentUserID, onDelete, onRevoke, onRemove, onSend }: { friend: Friend; messages: Message[]; loading: boolean; currentUserID: number; onDelete: (id: number) => void; onRevoke: (id: number) => void; onRemove: () => void; onSend: (content: string, files: File[]) => Promise<void> }) { return <><header className="flex items-center justify-between border-b border-sky-100 bg-white px-5 py-4"><div className="flex items-center gap-3"><Avatar name={friend.name} src={friend.avatar} /><div><h2 className="font-semibold">{friend.name}</h2><p className="text-xs text-slate-500">{friend.account}</p></div></div><button onClick={onRemove} className="text-sm text-red-600">Remove friend</button></header><div className="min-h-0 flex-1 overflow-y-auto p-5">{loading ? <p className="text-center text-sm text-slate-500">Loading messages…</p> : <MessageList messages={messages} currentUserId={currentUserID} onDelete={onDelete} onRevoke={onRevoke} />}</div><MessageComposer disabled={false} sending={false} onSend={onSend} /></>; }
function SettingsDialog({ user, onClose, onSignOut, onDelete }: { user: CurrentUser | null; onClose: () => void; onSignOut: () => void; onDelete: () => void }) { return <div className="fixed inset-0 z-40 grid place-items-center bg-sky-950/20 p-4"><section role="dialog" aria-modal="true" className="w-full max-w-sm rounded-2xl bg-white p-6 shadow-xl"><div className="flex items-center justify-between"><h2 className="text-lg font-bold">Account</h2><button onClick={onClose}><X /></button></div><div className="mt-5 flex items-center gap-3"><Avatar name={user?.name ?? ""} src={user?.avatar} className="h-12 w-12" /><div><p className="font-medium">{user?.name}</p><p className="text-sm text-slate-500">{user?.account}</p></div></div><button onClick={onSignOut} className="mt-6 flex w-full items-center justify-center gap-2 rounded-xl bg-sky-600 py-2.5 text-sm font-semibold text-white"><LogOut size={16} /> Sign out</button><button onClick={onDelete} className="mt-3 w-full text-sm text-red-600">Delete account</button></section></div>; }
