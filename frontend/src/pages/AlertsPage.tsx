import { AlertList } from "../components/alert-list";

export function AlertsPage() {
  return (
    <div className="p-6">
      <h2 className="text-xl font-semibold">Alerts</h2>
      <section className="mt-6">
        <AlertList />
      </section>
    </div>
  );
}
