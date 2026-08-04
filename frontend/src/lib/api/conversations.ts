import { api, unwrapResponse } from "./client";

export interface UserSearchResult {
  id: number;
  provider: string;
  account: string;
  name: string;
  avatar?: string;
}

export interface FriendRequestItem {
  id: number;
  senderID: number;
  receiverID: number;
  status: string;
  createdAt?: string;
}

export async function listConversations() {
  const response = await api.get("/friends");
  return unwrapResponse(response);
}

export async function listReceivedRequests() {
  const response = await api.get("/friends/requests/received");
  return unwrapResponse(response) as FriendRequestItem[];
}

export async function searchUsers(keyword: string) {
  const response = await api.get("/user/search", {
    params: {
      provider: "google",
      keyword,
      limit: 10,
    },
  });
  return unwrapResponse(response) as UserSearchResult[];
}

export async function createConversation(userId: string) {
  const response = await api.post("/friends/requests", { receiver_id: Number(userId) });
  return unwrapResponse(response);
}

export async function acceptFriendRequest(requestId: string | number) {
  const response = await api.post(`/friends/requests/${requestId}/accept`);
  return unwrapResponse(response);
}

export async function rejectFriendRequest(requestId: string | number) {
  const response = await api.post(`/friends/requests/${requestId}/reject`);
  return unwrapResponse(response);
}
