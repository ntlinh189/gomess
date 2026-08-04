import { api, unwrapResponse } from "./client";

export async function listMessages(friendId: string) {
  const response = await api.get(`/messages/${friendId}`);
  return unwrapResponse(response);
}

export async function sendMessage(receiverId: string, content: string) {
  const response = await api.post("/messages", { receiver_id: Number(receiverId), content });
  return unwrapResponse(response);
}
