"use client";

import { useMemo, useState } from "react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { useAuth } from "@/features/auth/context/auth-context";
import { acceptFriendRequest, createConversation, listConversations, listReceivedRequests, rejectFriendRequest, searchUsers, type UserSearchResult } from "@/lib/api";
import { ConversationList } from "./conversation-list";
import { MessageComposer } from "./message-composer";
import { MessageList } from "./message-list";

export function ChatShell() {
  const queryClient = useQueryClient();
  const { user } = useAuth();
  const [conversationId, setConversationId] = useState<string | undefined>();
  const [keyword, setKeyword] = useState("");
  const [results, setResults] = useState<UserSearchResult[]>([]);
  const [statusMessage, setStatusMessage] = useState("");

  const { data } = useQuery({
    queryKey: ["conversations"],
    queryFn: listConversations,
  });

  const { data: requestsData } = useQuery({
    queryKey: ["friend-requests"],
    queryFn: listReceivedRequests,
  });

  const conversations = useMemo(() => {
    return Array.isArray(data) ? data : (data as any)?.conversations || (data as any)?.items || [];
  }, [data]);

  const pendingRequests = useMemo(() => {
    return Array.isArray(requestsData) ? requestsData : [];
  }, [requestsData]);

  const selectedConversation = conversations.find((item: any) => item.id === conversationId || item._id === conversationId || item.conversationId === conversationId);

  const createMutation = useMutation({
    mutationFn: (value: string) => createConversation(value),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["conversations"] });
      setStatusMessage("Friend request sent.");
    },
    onError: (error: any) => {
      setStatusMessage(error?.response?.data?.error || error?.message || "Unable to send request.");
    },
  });

  const acceptMutation = useMutation({
    mutationFn: (requestId: string | number) => acceptFriendRequest(requestId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["friend-requests"] });
      queryClient.invalidateQueries({ queryKey: ["conversations"] });
      setStatusMessage("Friend request accepted.");
    },
  });

  const rejectMutation = useMutation({
    mutationFn: (requestId: string | number) => rejectFriendRequest(requestId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["friend-requests"] });
      setStatusMessage("Friend request rejected.");
    },
  });

  const handleSearch = async (value: string) => {
    const trimmed = value.trim();
    if (!trimmed) {
      setResults([]);
      setStatusMessage("");
      return;
    }

    try {
      const users = await searchUsers(trimmed);
      setResults(users);
      setStatusMessage(users.length ? "" : "No users found.");
    } catch (error: any) {
      setResults([]);
      setStatusMessage(error?.response?.data?.error || error?.message || "Search failed.");
    }
  };

  return (
    <div className="grid gap-4 overflow-hidden rounded-[28px] border border-sky-100 bg-white/90 shadow-[0_20px_60px_-20px_rgba(14,116,144,0.25)] lg:grid-cols-[360px_1fr]">
      <aside className="flex h-[760px] flex-col border-b border-sky-100 bg-slate-50/80 p-4 lg:border-b-0 lg:border-r">
        <div className="mb-4 flex items-center gap-3 rounded-2xl border border-sky-100 bg-white p-3 shadow-sm">
          {user?.avatar ? (
            <img src={user.avatar} alt={user?.name || "Me"} className="h-11 w-11 rounded-full object-cover" />
          ) : (
            <div className="flex h-11 w-11 items-center justify-center rounded-full bg-sky-600 text-sm font-semibold text-white">
              {(user?.name || user?.email || "M").charAt(0).toUpperCase()}
            </div>
          )}
          <div>
            <div className="text-sm font-semibold text-slate-800">{user?.name || user?.email || "Me"}</div>
            <div className="text-xs text-slate-500">Online</div>
          </div>
        </div>

        <div className="mb-4 rounded-2xl border border-sky-100 bg-white p-3 shadow-sm">
          <div className="mb-2 text-sm font-semibold text-slate-700">Find friends</div>
          <div className="flex gap-2">
            <Input
              value={keyword}
              onChange={(e) => setKeyword(e.target.value)}
              onKeyDown={(event) => {
                if (event.key === "Enter") {
                  event.preventDefault();
                  handleSearch(keyword);
                }
              }}
              placeholder="Search by name or email"
              className="h-9"
            />
            <Button onClick={() => handleSearch(keyword)} className="h-9">
              Search
            </Button>
          </div>
          {statusMessage ? <div className="mt-2 text-sm text-slate-600">{statusMessage}</div> : null}
          <div className="mt-3 space-y-2">
            {results.map((person) => (
              <div key={person.id} className="flex items-center justify-between rounded-xl border border-sky-100 bg-sky-50/70 p-2">
                <div className="flex items-center gap-2">
                  {person.avatar ? (
                    <img src={person.avatar} alt={person.name} className="h-9 w-9 rounded-full object-cover" />
                  ) : (
                    <div className="flex h-9 w-9 items-center justify-center rounded-full bg-sky-100 text-sm font-semibold text-sky-700">
                      {person.name?.charAt(0)?.toUpperCase() || "U"}
                    </div>
                  )}
                  <div>
                    <div className="text-sm font-medium text-slate-700">{person.name || person.account}</div>
                    <div className="text-xs text-slate-500">{person.account}</div>
                  </div>
                </div>
                <Button size="sm" onClick={() => createMutation.mutate(String(person.id))} disabled={createMutation.isPending}>
                  {createMutation.isPending ? "..." : "Add"}
                </Button>
              </div>
            ))}
          </div>
        </div>

        {pendingRequests.length > 0 ? (
          <div className="mb-4 rounded-2xl border border-amber-200 bg-amber-50/80 p-3 shadow-sm">
            <div className="mb-2 text-sm font-semibold text-amber-700">Friend requests</div>
            <div className="space-y-2">
              {pendingRequests.map((request) => (
                <div key={request.id} className="flex items-center justify-between rounded-xl bg-white px-2 py-2">
                  <div className="text-sm text-slate-700">Request #{request.id}</div>
                  <div className="flex gap-2">
                    <Button size="sm" variant="default" onClick={() => acceptMutation.mutate(request.id)} disabled={acceptMutation.isPending}>
                      Accept
                    </Button>
                    <Button size="sm" variant="outline" onClick={() => rejectMutation.mutate(request.id)} disabled={rejectMutation.isPending}>
                      Reject
                    </Button>
                  </div>
                </div>
              ))}
            </div>
          </div>
        ) : null}

        <div className="flex-1 overflow-auto">
          <ConversationList value={conversationId} onSelect={setConversationId} />
        </div>
      </aside>

      <section className="flex h-[760px] flex-col bg-white">
        <div className="flex items-center justify-between border-b border-sky-100 px-5 py-4">
          <div className="flex items-center gap-3">
            {selectedConversation?.avatar ? (
              <img src={selectedConversation.avatar} alt={selectedConversation.name || "Friend"} className="h-10 w-10 rounded-full object-cover" />
            ) : (
              <div className="flex h-10 w-10 items-center justify-center rounded-full bg-sky-100 text-sm font-semibold text-sky-700">
                {(selectedConversation?.name || "F").charAt(0).toUpperCase()}
              </div>
            )}
            <div>
              <div className="font-semibold text-slate-800">{selectedConversation?.name || "Select a conversation"}</div>
              <div className="text-sm text-slate-500">{selectedConversation ? "Active now" : "Choose a chat"}</div>
            </div>
          </div>
        </div>
        <div className="flex-1 overflow-hidden p-4">
          <MessageList conversationId={conversationId} />
        </div>
        <div className="border-t border-sky-100 p-4">
          <MessageComposer conversationId={conversationId} />
        </div>
      </section>
    </div>
  );
}
