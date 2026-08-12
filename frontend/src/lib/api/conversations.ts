import { api } from "./client";
import type { Friend, FriendRequest, OAuthProvider, UserSearchResult } from "@/types/user";

export async function listFriends() {
  const { data } = await api.get<Friend[]>("/friends");
  return data;
}

export async function searchUsers(provider: OAuthProvider, keyword: string, skip = 0, limit = 20) {
  const { data } = await api.get<UserSearchResult[]>("/user/search", { params: { provider, keyword, skip, limit } });
  return data;
}

export async function sendFriendRequest(receiverId: number) {
  await api.post("/friends/requests", { receiver_id: receiverId });
}

export async function listReceivedRequests() {
  const { data } = await api.get<FriendRequest[]>("/friends/requests/received");
  return data;
}

export async function acceptFriendRequest(requestId: number) { await api.post(`/friends/requests/${requestId}/accept`); }
export async function rejectFriendRequest(requestId: number) { await api.post(`/friends/requests/${requestId}/reject`); }
export async function removeFriend(friendId: number) { await api.delete(`/friends/${friendId}`); }
