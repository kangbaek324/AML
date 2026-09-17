import type { Stock } from "../../types/user";
import { formatNumber, formatWon } from "../../utils/format";

interface StockHoldingsTableProps {
  stocks: Stock[];
}

export function StockHoldingsTable({ stocks }: StockHoldingsTableProps) {
  if (stocks.length === 0) {
    return <p className="py-3 pl-20 text-sm text-gray-500">보유 종목이 없습니다.</p>;
  }

  return (
    <div className="pl-20 pr-4 pb-3">
      <div className="grid grid-cols-5 gap-4 border-b border-gray-200 px-3 py-1.5 text-xs font-medium text-gray-500">
        <span>종목명</span>
        <span className="text-right">수량</span>
        <span className="text-right">평균단가</span>
        <span className="text-right">현재가</span>
        <span className="text-right">평가금액</span>
      </div>
      {stocks.map((stock) => (
        <div
          key={stock.stock_id}
          className="grid grid-cols-5 gap-4 px-3 py-2 text-sm text-gray-700"
        >
          <span>{stock.name}</span>
          <span className="text-right">{formatNumber(stock.quantity)}</span>
          <span className="text-right">{formatWon(stock.average)}</span>
          <span className="text-right">{formatWon(stock.price)}</span>
          <span className="text-right">{formatWon(stock.valuation)}</span>
        </div>
      ))}
    </div>
  );
}
