import { mockCatalog } from "./mock";
import { CatalogError, type GetToken, type ListParams, type Page, type Plan, type PlanInput, type Product, type ProductInput } from "./types";

const mock = import.meta.env.VITE_USE_MOCK !== "false";
const base = (import.meta.env.VITE_API_URL ?? "http://localhost:3000").replace(/\/$/, "");

async function request<T>(path: string, token: GetToken, init: RequestInit = {}): Promise<T> {
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
const query = ({ q, page, pageSize }: ListParams) => `?q=${encodeURIComponent(q)}&page=${page}&pageSize=${pageSize}`;
const paged = <T,>(value: unknown): Page<T> => {
  if (!value || typeof value !== "object" || !("items" in value) || !Array.isArray(value.items) || !("total" in value)) {
    throw new CatalogError("The Go API does not yet return a paginated catalog with the required price and variant fields.", 501);
  }
  return value as Page<T>;
};
const product = (value: unknown): Product => {
  if (!value || typeof value !== "object" || !("salesPrice" in value) || !("costPrice" in value) || !("variants" in value)) {
    throw new CatalogError("The Go product API does not yet expose the requested catalog fields.", 501);
  }
  return value as Product;
};
const ensureProducts = async (token: GetToken) => { paged<Product>(await request<unknown>(`/products${query({ q: "", page: 1, pageSize: 1 })}`, token)); };
const ensurePlans = async (token: GetToken) => { paged<Plan>(await request<unknown>(`/plans${query({ q: "", page: 1, pageSize: 1 })}`, token)); };

export const catalog = {
  listProducts: (params: ListParams, token: GetToken, signal?: AbortSignal): Promise<Page<Product>> =>
    mock ? mockCatalog.listProducts(params) : request<unknown>(`/products${query(params)}`, token, { signal }).then(value => {
      const result = paged<Product>(value);
      result.items.forEach(product);
      return result;
    }),
  listPlans: (params: ListParams, token: GetToken, signal?: AbortSignal): Promise<Page<Plan>> =>
    mock ? mockCatalog.listPlans(params) : request<unknown>(`/plans${query(params)}`, token, { signal }).then(paged<Plan>),
  getProduct: (id: string, token: GetToken): Promise<Product> =>
    mock ? mockCatalog.getProduct(id) : request<unknown>(`/products/${encodeURIComponent(id)}`, token).then(product),
  getPlan: (id: string, token: GetToken): Promise<Plan> =>
    mock ? mockCatalog.getPlan(id) : request(`/plans/${encodeURIComponent(id)}`, token),
  saveProduct: async (id: string | undefined, input: ProductInput, token: GetToken): Promise<Product> => {
    if (mock) return mockCatalog.saveProduct(id, input);
    await ensureProducts(token);
    return request(`/admin/products${id ? `/${encodeURIComponent(id)}` : ""}`, token, { method: id ? "PUT" : "POST", body: JSON.stringify(input) });
  },
  savePlan: async (id: string | undefined, input: PlanInput, token: GetToken): Promise<Plan> => {
    if (mock) return mockCatalog.savePlan(id, input);
    await ensurePlans(token);
    return request(`/admin/plans${id ? `/${encodeURIComponent(id)}` : ""}`, token, { method: id ? "PUT" : "POST", body: JSON.stringify(input) });
  },
  deleteProduct: async (id: string, token: GetToken): Promise<void> => {
    if (mock) return mockCatalog.deleteProduct(id);
    await ensureProducts(token);
    return request(`/admin/products/${encodeURIComponent(id)}`, token, { method: "DELETE" });
  },
  deletePlan: async (id: string, token: GetToken): Promise<void> => {
    if (mock) return mockCatalog.deletePlan(id);
    await ensurePlans(token);
    return request(`/admin/plans/${encodeURIComponent(id)}`, token, { method: "DELETE" });
  },
};
