import { useCallback, useEffect, useState } from "react";
import type { Alert } from "../../types/alert";
import { fetchAlerts } from "../../services/api/alerts";
import { AlertRow } from "./AlertRow";

export function AlertList() {
  const [alerts, setAlerts] = useState<Alert[]>([]);
  const [loading, setLoading] = useState(true);
  const [refreshing, setRefreshing] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const loadAlerts = useCallback(() => {
    return fetchAlerts()
      .then(setAlerts)
      .catch(() => setError("Alert 목록을 불러오지 못했습니다."));
  }, []);

  useEffect(() => {
    loadAlerts().finally(() => setLoading(false));
  }, [loadAlerts]);

  function handleRefresh() {
    setRefreshing(true);
    setError(null);
    // 로컬 환경에서는 요청이 너무 빨리 끝나 회전 애니메이션이 안 보이므로 최소 노출 시간을 둔다.
    const minSpin = new Promise((resolve) => setTimeout(resolve, 500));
    Promise.all([loadAlerts(), minSpin]).finally(() => setRefreshing(false));
  }

  if (loading) {
    return <p className="text-sm text-gray-500">로딩중...</p>;
  }

  return (
    <div className="rounded border border-gray-200">
      {error && <p className="border-b border-gray-200 px-4 py-1.5 text-xs text-red-600">{error}</p>}
      <div className="flex items-center justify-end border-b border-gray-200 px-4 py-1.5">
        <button
          type="button"
          onClick={handleRefresh}
          disabled={refreshing}
          title="새로고침"
          className="text-gray-500 hover:text-gray-900 disabled:opacity-50"
        >
          <svg
            xmlns="http://www.w3.org/2000/svg"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            strokeWidth={2}
            strokeLinecap="round"
            strokeLinejoin="round"
            className={`h-3.5 w-3.5 ${refreshing ? "animate-spin" : ""}`}
          >
            <path d="M3 12a9 9 0 0 1 9-9 9.75 9.75 0 0 1 6.74 2.74L21 8" />
            <path d="M21 3v5h-5" />
            <path d="M21 12a9 9 0 0 1-9 9 9.75 9.75 0 0 1-6.74-2.74L3 16" />
            <path d="M8 16H3v5" />
          </svg>
        </button>
      </div>
      <div className="grid grid-cols-[3rem_4rem_9rem_1fr_12rem_11rem] gap-4 bg-gray-50 px-4 py-2 text-xs font-medium text-gray-500">
        <span>ID</span>
        <span>유저 ID</span>
        <span>유형</span>
        <span>사유</span>
        <span>발생일시</span>
        <span className="text-right">상태</span>
      </div>
      {alerts.length === 0 ? (
        <p className="px-4 py-3 text-sm text-gray-500">Alert가 없습니다.</p>
      ) : (
        alerts.map((alert) => <AlertRow key={alert.id} alert={alert} />)
      )}
    </div>
  );
}
