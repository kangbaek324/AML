import { useEffect, useMemo, useState } from "react";
import { LEVEL_RANK, type User } from "../../types/user";
import { fetchUsers } from "../../services/api/users";
import { UserRow } from "./UserRow";

type SortKey = "id" | "risk_level" | "asset_tier";
type SortDir = "asc" | "desc";

export function UserList() {
  const [users, setUsers] = useState<User[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [sortKey, setSortKey] = useState<SortKey>("id");
  const [sortDir, setSortDir] = useState<SortDir>("asc");

  useEffect(() => {
    fetchUsers()
      .then(setUsers)
      .catch(() => setError("유저 목록을 불러오지 못했습니다."))
      .finally(() => setLoading(false));
  }, []);

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

  if (error) {
    return <p className="text-sm text-gray-500">{error}</p>;
  }

  return (
    <div className="rounded border border-gray-200">
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
        <span className="col-span-2 text-right">갱신일시</span>
      </div>
      {sortedUsers.map((user) => (
        <UserRow key={user.id} user={user} />
      ))}
    </div>
  );
}
