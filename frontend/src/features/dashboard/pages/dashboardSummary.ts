export type Subscription = { customer: string; plan: string; status: string; started: string };
export type Invoice = { id: string; customer: string; plan: string; amount: number; issued: string; due: string; paid: string };
export type Filters = { from: string; to: string; plan: string; customer: string; status: string };

export function summarizeDashboard(subscriptions: Subscription[], invoices: Invoice[], today: string, filters: Filters) {
  const matches = (item: { customer: string; plan: string }) =>
    (filters.plan === "All plans" || item.plan === filters.plan) &&
    (filters.customer === "All customers" || item.customer === filters.customer);
  const active = subscriptions.filter(item => matches(item) && item.started <= filters.to &&
    item.status === "Active" && ["All statuses", "Active"].includes(filters.status)).length;
  const selected = invoices.filter(item => {
    if (!matches(item)) return false;
    switch (filters.status) {
      case "Paid": return Boolean(item.paid);
      case "Unpaid": return !item.paid;
      case "Overdue": return !item.paid && item.due < today;
      case "All statuses": return true;
      default: return subscriptions.some(sub => sub.customer === item.customer && sub.status === filters.status);
    }
  });
  const unpaid = selected.filter(item => !item.paid && item.issued >= filters.from && item.issued <= filters.to);
  const paid = selected.filter(item => item.paid >= filters.from && item.paid <= filters.to);
  const overdue = unpaid.filter(item => item.due < today).sort((a, b) => a.due.localeCompare(b.due));
  const months = new Map<string, number>();
  for (const item of paid) {
    const month = item.paid.slice(0, 7);
    months.set(month, (months.get(month) ?? 0) + item.amount);
  }
  return {
    active,
    revenue: paid.filter(item => item.paid.slice(0, 7) === today.slice(0, 7)).reduce((sum, item) => sum + item.amount, 0),
    outstanding: unpaid.reduce((sum, item) => sum + item.amount, 0),
    overdue,
    chart: [...months].sort(([a], [b]) => a.localeCompare(b)).map(([month, revenue]) => ({ month, revenue })),
  };
}
