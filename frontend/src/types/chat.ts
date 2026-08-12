export interface Attachment {
  id: number;
  type: "image" | "video" | "audio" | "file";
  url: string;
  file_name: string;
  mime_type: string;
  size_bytes: number;
}

export interface Message {
  id: number;
  sender_id: number;
  receiver_id: number;
  content: string;
  attachments: Attachment[];
  created_at: string;
  revoked: boolean;
}

export interface PresignedUpload {
  object_key: string;
  upload_url: string;
}

export type WebSocketEvent =
  | { event: "message.new"; data: Message }
  | { event: "message.deleted"; data: { id: number } }
  | { event: "message.revoked"; data: { id: number; sender_id: number; receiver_id: number } };
