export type ApiSubscription = {
  id: number; customerId: number; variantId: number; status: string; quantity: number;
  unitPriceCents: number; currency: string; billingInterval: string;
  currentPeriodStart: string; currentPeriodEnd: string; cancelAtPeriodEnd: boolean;
};
export type ApiInvoice = {
  id: number; invoiceNumber: string; subscriptionId: number; status: string;
  periodStart: string; periodEnd: string; subtotalCents: number; discountCents: number;
  totalCents: number; issuedAt?: string; dueAt?: string; paidAt?: string;
};
export type ApiVariant = { id: number; productId: number; name: string; priceCents: number; currency: string; billingInterval: string };
export type ApiProduct = { id: number; name: string };
export type GetToken = () => Promise<string | null>;
const base = (import.meta.env.VITE_API_URL ?? "http://localhost:3000").replace(/\/$/, "");

export async function billingRequest<T>(path: string, token: GetToken, init: RequestInit = {}): Promise<T> {
  const bearer = await token();
  if (!bearer) throw new Error("Your session expired. Sign in again to continue.");
  const response = await fetch(`${base}/api/staff${path}`, {
    ...init,
    headers: { Authorization: `Bearer ${bearer}`, ...(init.body ? { "Content-Type": "application/json" } : {}), ...init.headers },
  });
  if (!response.ok) {
    const body = await response.json().catch(() => ({}));
    throw new Error(body.error ?? `Request failed (${response.status})`);
  }
  return response.status === 204 ? undefined as T : response.json() as Promise<T>;
}

export async function loadBilling(token: GetToken) {
  const subscriptions = await billingRequest<ApiSubscription[]>("/subscriptions", token);
  const enriched = await Promise.all(subscriptions.map(async subscription => {
    const [invoices, variant] = await Promise.all([
      billingRequest<ApiInvoice[]>(`/subscriptions/${subscription.id}/invoices`, token),
      fetch(`${base}/api/variants/${subscription.variantId}`, { headers: { Authorization: `Bearer ${await token()}` } })
        .then(async response => response.ok ? response.json() as Promise<ApiVariant> : null).catch(() => null),
    ]);
    const product = variant ? await fetch(`${base}/api/products/${variant.productId}`, { headers: { Authorization: `Bearer ${await token()}` } })
      .then(async response => response.ok ? response.json() as Promise<ApiProduct> : null).catch(() => null) : null;
    return { subscription, invoices, variant, product };
  }));
  return enriched;
}
