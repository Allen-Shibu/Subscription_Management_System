import type { ListParams, Page, Plan, PlanInput, Product, ProductInput } from "./types.ts";
import { CatalogError } from "./types.ts";

const key = "recur-admin-catalog-v1";
type Data = { products: Product[]; plans: Plan[] };
const seed: Data = {
  products: [
    { id: "product-1", name: "Cloud Suite", type: "Service", salesPrice: 1000, costPrice: 400, notes: "Sample product", variants: [{ id: "variant-1", attribute: "Brand", value: "Odoo", extraPrice: 560 }] },
    { id: "product-2", name: "Setup Kit", type: "Physical", salesPrice: 250, costPrice: 90, notes: "", variants: [] },
  ],
  plans: [
    { id: "plan-1", name: "Growth", price: 2400, billingPeriod: "Monthly", minQuantity: 1, availableFrom: "", availableTo: "", autoClose: false, earlyClose: true, allowPause: true, allowRenew: true },
    { id: "plan-2", name: "Starter", price: 890, billingPeriod: "Yearly", minQuantity: 1, availableFrom: "", availableTo: "", autoClose: false, earlyClose: false, allowPause: false, allowRenew: true },
  ],
};

function read(): Data {
  try {
    const stored = localStorage.getItem(key);
    return stored ? JSON.parse(stored) as Data : structuredClone(seed);
  } catch {
    return structuredClone(seed);
  }
}
function write(data: Data) { localStorage.setItem(key, JSON.stringify(data)); }
function page<T extends { name: string }>(items: T[], params: ListParams): Page<T> {
  const filtered = items.filter(item => item.name.toLocaleLowerCase().includes(params.q.toLocaleLowerCase()));
  return { items: filtered.slice((params.page - 1) * params.pageSize, params.page * params.pageSize), total: filtered.length, page: params.page, pageSize: params.pageSize };
}
const pause = () => new Promise(resolve => setTimeout(resolve, 180));

export const mockCatalog = {
  async listProducts(params: ListParams) { await pause(); return page(read().products, params); },
  async listPlans(params: ListParams) { await pause(); return page(read().plans, params); },
  async getProduct(id: string) { await pause(); const item = read().products.find(item => item.id === id); if (!item) throw new CatalogError("Product not found", 404); return item; },
  async getPlan(id: string) { await pause(); const item = read().plans.find(item => item.id === id); if (!item) throw new CatalogError("Plan not found", 404); return item; },
  async saveProduct(id: string | undefined, input: ProductInput) {
    const data = read();
    if (id && !data.products.some(item => item.id === id)) throw new CatalogError("Product not found", 404);
    const item: Product = { ...input, id: id ?? crypto.randomUUID(), variants: input.variants.map(variant => ({ ...variant, id: variant.id ?? crypto.randomUUID() })) };
    data.products = id ? data.products.map(current => current.id === id ? item : current) : [item, ...data.products];
    write(data);
    return item;
  },
  async savePlan(id: string | undefined, input: PlanInput) {
    const data = read();
    if (id && !data.plans.some(item => item.id === id)) throw new CatalogError("Plan not found", 404);
    const item: Plan = { ...input, id: id ?? crypto.randomUUID() };
    data.plans = id ? data.plans.map(current => current.id === id ? item : current) : [item, ...data.plans];
    write(data);
    return item;
  },
  async deleteProduct(id: string) { const data = read(); data.products = data.products.filter(item => item.id !== id); write(data); },
  async deletePlan(id: string) { const data = read(); data.plans = data.plans.filter(item => item.id !== id); write(data); },
};
