import { api } from "./client";
import type { Message, PresignedUpload } from "@/types/chat";

export async function listMessages(friendId: number, beforeId?: number) {
  const { data } = await api.get<Message[]>(`/messages/${friendId}`, { params: { before_id: beforeId, limit: 100 } });
  return data;
}

export async function sendMessage(receiverId: number, content: string, attachments: Array<{ object_key: string; file_name: string }> = []) {
  const { data } = await api.post<Message>("/messages", { receiver_id: receiverId, content, attachments });
  return data;
}

export async function deleteMessage(messageId: number) { await api.delete(`/messages/${messageId}`); }
export async function revokeMessage(messageId: number) { await api.post(`/messages/${messageId}/revoke`); }

export async function presignUpload(fileName: string) {
  const { data } = await api.post<PresignedUpload>("/uploads/presign", { file_name: fileName });
  return data;
}

export async function uploadFile(uploadUrl: string, file: File) {
  const response = await fetch(uploadUrl, { method: "PUT", body: file, headers: { "Content-Type": file.type || "application/octet-stream" } });
  if (!response.ok) throw new Error("File upload failed");
}
