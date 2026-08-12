export type OAuthProvider = "google";

export interface CurrentUser {
  id: number;
  account: string;
  name: string;
  avatar: string;
}

export interface UserSearchResult extends CurrentUser {
  provider: OAuthProvider;
  request_status?: "pending";
}

export type Friend = UserSearchResult;

export interface FriendRequest {
  id: number;
  receiver_id: number;
  status: "pending" | "accepted" | "rejected";
  created_at: string;
  // Optional only while older API processes are restarted during rollout.
  sender?: Friend;
}
