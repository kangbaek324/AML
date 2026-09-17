import { useCallback, useEffect, useMemo, useState } from "react";
import { LEVEL_RANK, type User } from "../../types/user";
import { fetchUsers, refreshUsers } from "../../services/api/users";
import { UserRow } from "./UserRow";

type SortKey = "id" | "risk_level" | "asset_tier";
type SortDir = "asc" | "desc";

export function UserList() {
  const [users, setUsers] = useState<User[]>([]);
  const [loading, setLoading] = useState(true);
  const [refreshing, setRefreshing] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [sortKey, setSortKey] = useState<SortKey>("id");
  const [sortDir, setSortDir] = useState<SortDir>("asc");

  const loadUsers = useCallback(() => {
    return fetchUsers()
      .then(setUsers)
      .catch(() => setError("유저 목록을 불러오지 못했습니다."));
  }, []);

  useEffect(() => {
    loadUsers().finally(() => setLoading(false));
  }, [loadUsers]);

  function handleRefresh() {
    setRefreshing(true);
    setError(null);
    // 로컬 환경에서는 요청이 너무 빨리 끝나 회전 애니메이션이 안 보이므로 최소 노출 시간을 둔다.
    const minSpin = new Promise((resolve) => setTimeout(resolve, 500));
    Promise.all([refreshUsers().then(loadUsers), minSpin])
      .catch(() => setError("갱신에 실패했습니다."))
      .finally(() => setRefreshing(false));
  }

  const sortedUsers = useMemo(() => {
    const value = (user: User) =>
      sortKey === "id" ? user.id : LEVEL_RANK[user[sortKey]];

    return [...users].sort((a, b) => {
      const diff = value(a) - value(b);
      return sortDir === "asc" ? diff : -diff;
    });
  }, [users, sortKey, sortDir]);

  function handleSort(key: SortKey) {
    if (key === sortKey) {
      setSortDir((prev) => (prev === "asc" ? "desc" : "asc"));
    } else {
      setSortKey(key);
      setSortDir("asc");
    }
  }

  function sortIndicator(key: SortKey) {
    const active = sortKey === key;
    const arrow = active ? (sortDir === "asc" ? "▲" : "▼") : "▲";
    return (
      <span className={`ml-1 text-[8px] ${active ? "text-gray-700" : "text-gray-300"}`}>
        {arrow}
      </span>
    );
  }

  if (loading) {
    return <p className="text-sm text-gray-500">로딩중...</p>;
  }

  return (
    <div className="rounded border border-gray-200">
      {error && <p className="border-b border-gray-200 px-4 py-1.5 text-xs text-red-600">{error}</p>}
      <div className="grid grid-cols-6 gap-4 bg-gray-50 px-4 py-2 text-xs font-medium text-gray-500">
        <button type="button" onClick={() => handleSort("id")} className="text-left hover:text-gray-900">
          ID{sortIndicator("id")}
        </button>
        <span className="text-right">평균 자산</span>
        <button
          type="button"
          onClick={() => handleSort("risk_level")}
          className="text-right hover:text-gray-900"
        >
          위험도{sortIndicator("risk_level")}
        </button>
        <button
          type="button"
          onClick={() => handleSort("asset_tier")}
          className="text-right hover:text-gray-900"
        >
          자산 등급{sortIndicator("asset_tier")}
        </button>
        <span className="col-span-2 flex items-center justify-end gap-1.5">
          갱신일시
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
        </span>
      </div>
      {sortedUsers.map((user) => (
        <UserRow key={user.id} user={user} />
      ))}
    </div>
  );
}
