const currencyFormatter = new Intl.NumberFormat("ko-KR");

export function formatWon(value: number | string): string {
  const numeric = typeof value === "string" ? Number(value) : value;
  return `${currencyFormatter.format(numeric)}원`;
}

export function formatNumber(value: number): string {
  return currencyFormatter.format(value);
}

export function formatDateTime(value: string): string {
  return new Date(value).toLocaleString("ko-KR");
}
