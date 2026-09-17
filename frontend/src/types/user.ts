export type RiskLevel = "LOW" | "MEDIUM" | "HIGH";
export type AssetTier = "LOW" | "MEDIUM" | "HIGH";
export type AccountStatus = "PENDING" | "ACTIVE";

export const LEVEL_RANK: Record<RiskLevel | AssetTier, number> = {
  LOW: 0,
  MEDIUM: 1,
  HIGH: 2,
};

export interface User {
  id: number;
  average_asset: string;
  risk_level: RiskLevel;
  asset_tier: AssetTier;
  updated_at: string;
}

export interface Account {
  id: number;
  account_number: number;
  balance: number;
  available_balance: number;
  status: AccountStatus;
}

export interface UserDetail {
  id: number;
  average_asset: string;
  risk_level: RiskLevel;
  asset_tier: AssetTier;
  accounts: Account[];
}

export interface Stock {
  stock_id: number;
  name: string;
  quantity: number;
  average: number;
  price: number;
  valuation: number;
}

export interface AccountDetail {
  id: number;
  user_id: number;
  account_number: number;
  balance: number;
  available_balance: number;
  status: AccountStatus;
  stocks: Stock[];
}
