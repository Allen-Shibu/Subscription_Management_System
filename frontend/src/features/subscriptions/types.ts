export type SubscriptionStatus = "Draft" | "Quotation" | "Confirmed" | "Active" | "Closed";

export type SubscriptionLine = {
  product: string;
  quantity: number;
  unitPrice: number;
  taxRate: number;
};

export type Subscription = {
  id: string;
  number: string;
  customer: string;
  initials: string;
  plan: string;
  startDate: string;
  expiryDate: string;
  paymentTerms: string;
  status: SubscriptionStatus;
  sourceTemplate?: string;
  lines: SubscriptionLine[];
};
