import type { ReactNode } from "react";
import { useState } from "react";
import type { Alert, AlertDetail } from "../../types/alert";
import { fetchAlertDetail } from "../../services/api/alerts";
import { formatDateTime, formatWon } from "../../utils/format";

interface AlertRowProps {
  alert: Alert;
}

function DetailCard({ title, children }: { title: string; children: ReactNode }) {
  return (
    <div className="max-w-xs rounded border border-gray-200 bg-white p-3">
      <p className="mb-2 text-xs font-medium text-gray-500">{title}</p>
      <dl className="space-y-1">{children}</dl>
    </div>
  );
}

function DetailRow({ label, children }: { label: string; children: ReactNode }) {
  return (
    <div className="flex justify-between gap-4">
      <dt className="text-gray-500">{label}</dt>
      <dd className="text-gray-700">{children}</dd>
    </div>
  );
}

export function AlertRow({ alert }: AlertRowProps) {
  const [expanded, setExpanded] = useState(false);
  const [detail, setDetail] = useState<AlertDetail | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  function handleToggle() {
    setExpanded((prev) => !prev);

    if (!detail && !loading) {
      setLoading(true);
      setError(null);
      fetchAlertDetail(alert.id)
        .then(setDetail)
        .catch(() => setError("Alert 상세를 불러오지 못했습니다."))
        .finally(() => setLoading(false));
    }
  }

  return (
    <div className="border-b border-gray-200">
      <button
        type="button"
        onClick={handleToggle}
        className="grid w-full grid-cols-[3rem_4rem_9rem_1fr_10rem_5rem] gap-4 px-4 py-3 text-left text-sm text-gray-900 hover:bg-gray-100"
      >
        <span>
          <span className="mr-1 inline-block w-3 text-gray-400">{expanded ? "▾" : "▸"}</span>
          {alert.id}
        </span>
        <span>{alert.user_id}</span>
        <span>
          <span className="rounded border border-gray-300 px-2 py-0.5 text-xs text-gray-600">
            {alert.type}
          </span>
        </span>
        <span className="truncate text-gray-700" title={alert.reason}>
          {alert.reason}
        </span>
        <span className="text-gray-500">{formatDateTime(alert.alerted_at)}</span>
        <span className="text-right">
          <span className="rounded border border-gray-300 px-2 py-0.5 text-xs text-gray-600">
            {alert.status}
          </span>
        </span>
      </button>

      {expanded && (
        <div className="space-y-2 bg-gray-50 py-3 pl-10 pr-4 text-sm">
          {loading && <p className="text-gray-500">로딩중...</p>}
          {error && <p className="text-gray-500">{error}</p>}
          {detail && detail.trades.length === 0 && detail.transfers.length === 0 && (
            <p className="text-gray-500">연결된 거래/송금 정보가 없습니다.</p>
          )}
          {detail?.trades.map((t) => (
            <DetailCard key={t.trade_id} title={`Trade #${t.trade_id}`}>
              <DetailRow label="종목">{t.stock_id}</DetailRow>
              <DetailRow label="수량 × 가격">
                {t.quantity.toLocaleString()}주 × {formatWon(t.price)}
              </DetailRow>
              <DetailRow label="Maker">User #{t.maker_user_id}</DetailRow>
              <DetailRow label="Taker">User #{t.taker_user_id}</DetailRow>
              <DetailRow label="체결일시">{formatDateTime(t.matched_at)}</DetailRow>
            </DetailCard>
          ))}
          {detail?.transfers.map((tr) => (
            <DetailCard key={tr.transfer_id} title={`Transfer #${tr.transfer_id}`}>
              <DetailRow label="From">User #{tr.sender_user_id}</DetailRow>
              <DetailRow label="To">User #{tr.recipient_user_id}</DetailRow>
              <DetailRow label="금액">{formatWon(tr.amount)}</DetailRow>
              <DetailRow label="상태">{tr.status}</DetailRow>
              <DetailRow label="일시">{formatDateTime(tr.created_at)}</DetailRow>
            </DetailCard>
          ))}
        </div>
      )}
    </div>
  );
}
