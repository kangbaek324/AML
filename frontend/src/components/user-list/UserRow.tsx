import { useState } from "react";
import type { User, UserDetail } from "../../types/user";
import { fetchUserDetail } from "../../services/api/users";
import { formatDateTime, formatWon } from "../../utils/format";
import { AccountRow } from "./AccountRow";

interface UserRowProps {
  user: User;
}

export function UserRow({ user }: UserRowProps) {
  const [expanded, setExpanded] = useState(false);
  const [detail, setDetail] = useState<UserDetail | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  function handleToggle() {
    setExpanded((prev) => !prev);

    if (!detail && !loading) {
      setLoading(true);
      setError(null);
      fetchUserDetail(user.id)
        .then(setDetail)
        .catch(() => setError("유저 정보를 불러오지 못했습니다."))
        .finally(() => setLoading(false));
    }
  }

  return (
    <div className="border-b border-gray-200">
      <button
        type="button"
        onClick={handleToggle}
        className="grid w-full grid-cols-6 gap-4 px-4 py-3 text-left text-sm text-gray-900 hover:bg-gray-100"
      >
        <span>
          <span className="mr-2 inline-block w-3 text-gray-400">{expanded ? "▾" : "▸"}</span>
          {user.id}
        </span>
        <span className="text-right">{formatWon(user.average_asset)}</span>
        <span className="text-right">
          <span className="rounded border border-gray-300 px-2 py-0.5 text-xs text-gray-600">
            {user.risk_level}
          </span>
        </span>
        <span className="text-right">
          <span className="rounded border border-gray-300 px-2 py-0.5 text-xs text-gray-600">
            {user.asset_tier}
          </span>
        </span>
        <span className="col-span-2 text-right text-gray-500">
          {formatDateTime(user.updated_at)}
        </span>
      </button>

      {expanded && (
        <div className="bg-gray-50">
          {loading && <p className="py-3 pl-10 text-sm text-gray-500">로딩중...</p>}
          {error && <p className="py-3 pl-10 text-sm text-gray-500">{error}</p>}
          {detail && (
            <div>
              <div className="grid grid-cols-4 gap-4 border-y border-gray-200 pl-10 pr-4 py-1.5 text-xs font-medium text-gray-500">
                <span>계좌번호</span>
                <span className="text-right">잔액</span>
                <span className="text-right">가용잔액</span>
                <span className="text-right">상태</span>
              </div>
              {detail.accounts.length === 0 ? (
                <p className="py-3 pl-10 text-sm text-gray-500">보유 계좌가 없습니다.</p>
              ) : (
                detail.accounts.map((account) => <AccountRow key={account.id} account={account} />)
              )}
            </div>
          )}
        </div>
      )}
    </div>
  );
}
