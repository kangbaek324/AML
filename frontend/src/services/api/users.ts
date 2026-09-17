import { apiClient } from "./client";
import type { AccountDetail, User, UserDetail } from "../../types/user";

export async function fetchUsers(): Promise<User[]> {
  const { data } = await apiClient.get<User[]>("/api/v1/users");
  return data;
}

export async function fetchUserDetail(userId: number): Promise<UserDetail> {
  const { data } = await apiClient.get<UserDetail>(`/api/v1/users/${userId}`);
  return data;
}

export async function fetchAccountDetail(accountId: number): Promise<AccountDetail> {
  const { data } = await apiClient.get<AccountDetail>(`/api/v1/accounts/${accountId}`);
  return data;
}
