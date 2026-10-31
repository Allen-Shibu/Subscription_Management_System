import assert from "node:assert/strict";
import { test } from "node:test";
import { productSchema, planSchema, planAvailability } from "./validation.ts";
import { mockCatalog } from "./mock.ts";
import { pageProducts } from "./pageProducts.ts";

const storage = new Map();
globalThis.localStorage = {
  getItem: key => storage.get(key) ?? null,
  setItem: (key, value) => storage.set(key, value),
  clear: () => storage.clear(),
};

test("product variants and plan dates validate", () => {
  const product = { name: "Cloud Suite", type: "Service", salesPrice: 1000, costPrice: 1100, notes: "", variants: [
    { attribute: "Brand", value: "Odoo", extraPrice: 560 },
  ] };
  assert.equal(productSchema.safeParse(product).success, true); // negative margin is a warning
  assert.equal(productSchema.safeParse({ ...product, salesPrice: 100.01 }).success, true);
  assert.equal(product.salesPrice + product.variants[0].extraPrice, 1560);
  assert.equal(productSchema.safeParse({ ...product, variants: [...product.variants, { ...product.variants[0], value: "odoo" }] }).success, false);
  assert.equal(productSchema.safeParse({ ...product, variants: [{ ...product.variants[0], extraPrice: -1 }] }).success, false);

  const plan = { name: "Growth", price: 50, billingPeriod: "Monthly", minQuantity: 1,
    availableFrom: "2026-10-01", availableTo: "2026-10-31", autoClose: false, earlyClose: false, allowPause: true, allowRenew: true };
  assert.equal(planSchema.safeParse(plan).success, true);
  assert.equal(planSchema.safeParse({ ...plan, minQuantity: 0 }).success, false);
  assert.equal(planSchema.safeParse({ ...plan, availableTo: "2026-09-30" }).success, false);
  assert.equal(planAvailability(plan, "2026-09-30"), "Scheduled");
  assert.equal(planAvailability(plan, "2026-10-15"), "Available");
  assert.equal(planAvailability(plan, "2026-11-01"), "Ended");
});

test("mock catalog creates, edits, searches, paginates, and deletes", async () => {
  storage.clear();
  const product = await mockCatalog.saveProduct(undefined, {
    name: "Test product", type: "Digital", salesPrice: 100, costPrice: 30, notes: "",
    variants: [{ attribute: "Brand", value: "Odoo", extraPrice: 20 }],
  });
  assert.equal((await mockCatalog.listProducts({ q: "test", page: 1, pageSize: 1 })).items[0].id, product.id);
  assert.equal((await mockCatalog.listProducts({ q: "", page: 2, pageSize: 1 })).items.length, 1);
  await mockCatalog.saveProduct(product.id, { ...product, name: "Edited product" });
  assert.equal((await mockCatalog.getProduct(product.id)).name, "Edited product");
  await mockCatalog.deleteProduct(product.id);
  assert.equal((await mockCatalog.listProducts({ q: "Edited", page: 1, pageSize: 5 })).total, 0);

  const plan = await mockCatalog.savePlan(undefined, {
    name: "Test plan", price: 45, billingPeriod: "Weekly", minQuantity: 2,
    availableFrom: "", availableTo: "", autoClose: false, earlyClose: false, allowPause: true, allowRenew: true,
  });
  assert.equal((await mockCatalog.listPlans({ q: "test", page: 1, pageSize: 5 })).total, 1);
  await mockCatalog.savePlan(plan.id, { ...plan, price: 60 });
  assert.equal((await mockCatalog.getPlan(plan.id)).price, 60);
  await mockCatalog.deletePlan(plan.id);
  assert.equal((await mockCatalog.listPlans({ q: "test", page: 1, pageSize: 5 })).total, 0);
});

test("Go product list is searched and paged in the browser", () => {
  const items = [{ id: 1, name: "Alpha", description: "", active: true }, { id: 2, name: "Beta", description: "", active: false }, { id: 3, name: "Alpine", description: "", active: true }];
  assert.deepEqual(pageProducts(items, { q: "al", page: 2, pageSize: 1 }), { items: [items[2]], total: 2, page: 2, pageSize: 1 });
});
