import { useState } from "react";
import type { Account, AccountDetail } from "../../types/user";
import { fetchAccountDetail } from "../../services/api/users";
import { formatWon } from "../../utils/format";
import { StockHoldingsTable } from "./StockHoldingsTable";

interface AccountRowProps {
  account: Account;
}

export function AccountRow({ account }: AccountRowProps) {
  const [expanded, setExpanded] = useState(false);
  const [detail, setDetail] = useState<AccountDetail | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  function handleToggle() {
    setExpanded((prev) => !prev);

    if (!detail && !loading) {
      setLoading(true);
      setError(null);
      fetchAccountDetail(account.id)
        .then(setDetail)
        .catch(() => setError("계좌 정보를 불러오지 못했습니다."))
        .finally(() => setLoading(false));
    }
  }

  return (
    <div className="border-t border-gray-100">
      <button
        type="button"
        onClick={handleToggle}
        className="grid w-full grid-cols-4 gap-4 py-2 pl-10 pr-4 text-left text-sm text-gray-700 hover:bg-gray-100"
      >
        <span>
          <span className="mr-2 inline-block w-3 text-gray-400">{expanded ? "▾" : "▸"}</span>
          {account.account_number}
        </span>
        <span className="text-right">{formatWon(account.balance)}</span>
        <span className="text-right">{formatWon(account.available_balance)}</span>
        <span className="text-right">
          <span className="rounded border border-gray-300 px-2 py-0.5 text-xs text-gray-600">
            {account.status}
          </span>
        </span>
      </button>

      {expanded && (
        <div className="bg-gray-50">
          {loading && <p className="py-3 pl-20 text-sm text-gray-500">로딩중...</p>}
          {error && <p className="py-3 pl-20 text-sm text-gray-500">{error}</p>}
          {detail && <StockHoldingsTable stocks={detail.stocks} />}
        </div>
      )}
    </div>
  );
}
