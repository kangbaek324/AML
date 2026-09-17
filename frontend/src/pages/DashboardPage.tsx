import { UserList } from "../components/user-list";

export function DashboardPage() {
  return (
    <div className="p-6">
      <h2 className="text-xl font-semibold">Dashboard</h2>
      <section className="mt-6">
        <h3 className="mb-3 text-base font-semibold text-gray-900">전체 유저 리스트</h3>
        <UserList />
      </section>
    </div>
  );
}
