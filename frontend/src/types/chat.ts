export interface Conversation {
  id?: string;
  _id?: string;
  conversationId?: string;
  name?: string;
  lastMessage?: string;
}

export interface Message {
  id?: string;
  _id?: string;
  text?: string;
  content?: string;
  body?: string;
  sender?: {
    name?: string;
    email?: string;
  };
  user?: {
    name?: string;
    email?: string;
  };
}
