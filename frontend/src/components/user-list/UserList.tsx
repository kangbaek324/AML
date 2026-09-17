import { useEffect, useState } from "react";
import type { User } from "../../types/user";
import { fetchUsers } from "../../services/api/users";
import { UserRow } from "./UserRow";

export function UserList() {
  const [users, setUsers] = useState<User[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    fetchUsers()
      .then(setUsers)
      .catch(() => setError("유저 목록을 불러오지 못했습니다."))
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
      <div className="grid grid-cols-6 gap-4 bg-gray-50 px-4 py-2 text-xs font-medium text-gray-500">
        <span>ID</span>
        <span className="text-right">평균 자산</span>
        <span className="text-right">위험도</span>
        <span className="text-right">자산 등급</span>
        <span className="col-span-2 text-right">갱신일시</span>
      </div>
      {users.map((user) => (
        <UserRow key={user.id} user={user} />
      ))}
    </div>
  );
}
