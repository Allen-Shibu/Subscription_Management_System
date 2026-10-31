import assert from "node:assert/strict";
import { test } from "node:test";
import { summarizeDashboard } from "./dashboardSummary.ts";

const subscriptions = [
  { customer: "A", plan: "Growth", status: "Active", started: "2026-01-01" },
  { customer: "B", plan: "Starter", status: "Paused", started: "2026-01-01" },
];
const invoices = [
  { id: "1", customer: "A", plan: "Growth", amount: 100, issued: "2026-10-01", due: "2026-10-10", paid: "" },
  { id: "2", customer: "A", plan: "Growth", amount: 200, issued: "2026-10-02", due: "2026-10-12", paid: "2026-10-15" },
  { id: "3", customer: "B", plan: "Starter", amount: 50, issued: "2026-10-03", due: "2026-11-01", paid: "" },
  { id: "4", customer: "A", plan: "Growth", amount: 75, issued: "2026-09-29", due: "2026-10-10", paid: "2026-10-03" },
];
const filters = { from: "2026-10-01", to: "2026-10-31", plan: "All plans", customer: "All customers", status: "All statuses" };

test("dashboard totals and filters use the same invoices", () => {
  const all = summarizeDashboard(subscriptions, invoices, "2026-10-31", filters);
  assert.equal(all.active, 1);
  assert.equal(all.revenue, 275);
  assert.equal(all.outstanding, 150);
  assert.deepEqual(all.overdue.map(item => item.id), ["1"]);
  assert.deepEqual(all.chart, [{ month: "2026-10", revenue: 275 }]);

  const growth = summarizeDashboard(subscriptions, invoices, "2026-10-31", { ...filters, plan: "Growth", status: "Overdue" });
  assert.equal(growth.outstanding, 100);
  assert.equal(growth.overdue.length, 1);
  assert.equal(growth.revenue, 0);
});
