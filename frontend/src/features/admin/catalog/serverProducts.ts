import { request } from "./api";
import type { GetToken } from "./types";

export type ServerProduct = { id: number; name: string; description: string; active: boolean };
export type ServerVariant = { id: number; productId: number; name: string; description: string; priceCents: number; currency: string; billingInterval: "month" | "year"; active: boolean };
export type ProductFields = Pick<ServerProduct, "name" | "description" | "active">;
export type VariantFields = Pick<ServerVariant, "name" | "description" | "priceCents" | "currency" | "billingInterval" | "active">;

export const serverProducts = {
  list: (token: GetToken, signal?: AbortSignal) => request<ServerProduct[]>("/products", token, { signal }),
  get: (id: string, token: GetToken) => request<ServerProduct>(`/products/${encodeURIComponent(id)}`, token),
  variants: (id: string, token: GetToken) => request<ServerVariant[]>(`/products/${encodeURIComponent(id)}/variants`, token),
  save: (id: string | undefined, input: ProductFields, token: GetToken) => request<ServerProduct>(`/admin/products${id ? `/${encodeURIComponent(id)}` : ""}`, token, { method: id ? "PUT" : "POST", body: JSON.stringify(input) }),
  delete: (id: number, token: GetToken) => request<void>(`/admin/products/${id}`, token, { method: "DELETE" }),
  saveVariant: (productId: string, id: number | undefined, input: VariantFields, token: GetToken) => request<ServerVariant>(id ? `/admin/variants/${id}` : `/admin/products/${encodeURIComponent(productId)}/variants`, token, { method: id ? "PUT" : "POST", body: JSON.stringify(input) }),
  deleteVariant: (id: number, token: GetToken) => request<void>(`/admin/variants/${id}`, token, { method: "DELETE" }),
};
