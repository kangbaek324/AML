import { useEffect, useState } from "react";
import type { Alert } from "../../types/alert";
import { fetchAlerts } from "../../services/api/alerts";
import { AlertRow } from "./AlertRow";

export function AlertList() {
  const [alerts, setAlerts] = useState<Alert[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    fetchAlerts()
      .then(setAlerts)
      .catch(() => setError("Alert 목록을 불러오지 못했습니다."))
      .finally(() => setLoading(false));
  }, []);

  if (loading) {
    return <p className="text-sm text-gray-500">로딩중...</p>;
  }

  if (error) {
    return <p className="text-sm text-gray-500">{error}</p>;
  }

  return (
    <div className="rounded border border-gray-200">
      <div className="grid grid-cols-[3rem_4rem_9rem_1fr_10rem_5rem] gap-4 bg-gray-50 px-4 py-2 text-xs font-medium text-gray-500">
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
