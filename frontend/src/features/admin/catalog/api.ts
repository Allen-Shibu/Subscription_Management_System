import { mockCatalog } from "./mock";
import { CatalogError, type GetToken, type ListParams, type PlanInput, type ProductInput } from "./types";

export const usingMock = import.meta.env?.VITE_USE_MOCK !== "false";
const base = (import.meta.env?.VITE_API_URL ?? "http://localhost:3000").replace(/\/$/, "");

export async function request<T>(path: string, token: GetToken, init: RequestInit = {}): Promise<T> {
  const bearer = await token();
  if (!bearer) throw new CatalogError("Sign in again to continue", 401);
  const response = await fetch(`${base}/api${path}`, {
    ...init,
    headers: { Authorization: `Bearer ${bearer}`, ...(init.body ? { "Content-Type": "application/json" } : {}), ...init.headers },
  });
  if (!response.ok) {
    const body = await response.json().catch(() => ({}));
    throw new CatalogError(body.error ?? `Request failed (${response.status})`, response.status, body.fields ?? body.errors ?? {});
  }
  return response.status === 204 ? undefined as T : response.json() as Promise<T>;
}
export const catalog = {
  listProducts: (params: ListParams, _token: GetToken, _signal?: AbortSignal) => mockCatalog.listProducts(params),
  listPlans: (params: ListParams, _token: GetToken, _signal?: AbortSignal) => mockCatalog.listPlans(params),
  getProduct: (id: string, _token: GetToken) => mockCatalog.getProduct(id),
  getPlan: (id: string, _token: GetToken) => mockCatalog.getPlan(id),
  saveProduct: (id: string | undefined, input: ProductInput, _token: GetToken) => mockCatalog.saveProduct(id, input),
  savePlan: (id: string | undefined, input: PlanInput, _token: GetToken) => mockCatalog.savePlan(id, input),
  deleteProduct: (id: string, _token: GetToken) => mockCatalog.deleteProduct(id),
  deletePlan: (id: string, _token: GetToken) => mockCatalog.deletePlan(id),
};
