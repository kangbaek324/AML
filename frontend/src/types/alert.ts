export type AlertType = "CROSS_TRADING" | "LARGE_TRANSACTION" | "ABNORMAL_TRANSFER";
export type AlertStatus = "PENDING" | "NORMAL" | "ABNORMAL";

export interface Alert {
  id: number;
  user_id: number;
  type: AlertType;
  reason: string;
  status: AlertStatus;
  alerted_at: string;
  processed_at: string | null;
}

export interface Trade {
  trade_id: number;
  stock_id: number;
  quantity: number;
  price: number;
  matched_at: string;
  maker_user_id: number;
  taker_user_id: number;
}

export interface Transfer {
  transfer_id: number;
  sender_user_id: number;
  recipient_user_id: number;
  amount: number;
  status: "RECEIVED" | "REJECTED" | "COMPLETED";
  created_at: string;
  completed_at: string | null;
}

export interface AlertDetail extends Alert {
  trades: Trade[];
  transfers: Transfer[];
}
