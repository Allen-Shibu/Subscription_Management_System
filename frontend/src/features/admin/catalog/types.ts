export type Variant = { id?: string; attribute: string; value: string; extraPrice: number };
export type ProductInput = {
  name: string;
  type: "Service" | "Physical" | "Digital";
  salesPrice: number;
  costPrice: number;
  notes: string;
  variants: Variant[];
};
export type Product = ProductInput & { id: string };

export type PlanInput = {
  name: string;
  price: number;
  billingPeriod: "Daily" | "Weekly" | "Monthly" | "Yearly";
  minQuantity: number;
  availableFrom: string;
  availableTo: string;
  autoClose: boolean;
  earlyClose: boolean;
  allowPause: boolean;
  allowRenew: boolean;
};
export type Plan = PlanInput & { id: string };
export type ListParams = { q: string; page: number; pageSize: number };
export type Page<T> = { items: T[]; total: number; page: number; pageSize: number };
export type GetToken = () => Promise<string | null>;

export class CatalogError extends Error {
  status: number;
  fields: Record<string, string>;
  constructor(message: string, status: number, fields: Record<string, string> = {}) {
    super(message);
    this.status = status;
    this.fields = fields;
  }
}
