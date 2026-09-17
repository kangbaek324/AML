import { apiClient } from "./client";
import type { Alert, AlertDetail } from "../../types/alert";

export async function fetchAlerts(): Promise<Alert[]> {
  const { data } = await apiClient.get<Alert[]>("/api/v1/alerts");
  return data;
}

export async function fetchAlertDetail(alertId: number): Promise<AlertDetail> {
  const { data } = await apiClient.get<AlertDetail>(`/api/v1/alerts/${alertId}`);
  return data;
}

export async function updateAlertStatus(
  alertId: number,
  status: "NORMAL" | "ABNORMAL",
): Promise<void> {
  await apiClient.patch(`/api/v1/alerts/${alertId}/status`, { status });
}
