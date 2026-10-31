import { useAuth } from "@clerk/react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { billingRequest, loadBilling, type ApiInvoice } from "../../billing/api";
import type { Subscription } from "../types";
import { SubscriptionSummary } from "../components/SubscriptionSummary";
import { SubscriptionsHeader } from "../components/SubscriptionsHeader";
import { SubscriptionTable } from "../components/SubscriptionTable";
import styles from "./SubscriptionsPage.module.css";

export function SubscriptionsPage() {
  const { getToken } = useAuth();
  const client = useQueryClient();
  const query = useQuery({ queryKey: ["billing-data"], queryFn: () => loadBilling(getToken) });
  const issue = useMutation({ mutationFn: (id: string) => billingRequest<ApiInvoice>(`/subscriptions/${id}/invoices`, getToken, { method: "POST" }), onSuccess: () => client.invalidateQueries({ queryKey: ["billing-data"] }) });
  const subscriptions: Subscription[] = (query.data ?? []).map(({ subscription: sub, variant, product }) => {
    const date = (value: string) => new Intl.DateTimeFormat("en-US", { dateStyle: "medium" }).format(new Date(value));
    const status = ({ active: "Active", trialing: "Active", past_due: "Confirmed", canceled: "Closed", expired: "Closed" } as const)[sub.status as "active" | "trialing" | "past_due" | "canceled" | "expired"] ?? "Active";
    return { id: String(sub.id), number: `SUB-${String(sub.id).padStart(6, "0")}`, customer: `Customer #${sub.customerId}`, initials: `C${sub.customerId}`,
      plan: variant ? `${product?.name ?? `Product #${variant.productId}`} · ${variant.name}` : `Variant #${sub.variantId}`,
      startDate: date(sub.currentPeriodStart), expiryDate: date(sub.currentPeriodEnd), paymentTerms: "Due on receipt", status,
      lines: [{ product: variant?.name ?? `Variant #${sub.variantId}`, quantity: sub.quantity, unitPrice: sub.unitPriceCents / 100, taxRate: 0 }] };
  });
  if (query.isLoading) return <main className={styles.subscriptions} aria-busy="true" aria-label="Loading subscriptions">
    <div className={styles.loadingHeading}><span/><span/><span/></div><div className={styles.loadingSummary}>{[1, 2, 3, 4].map(item => <div key={item} className={styles.loadingCard}><span/><span/><span/></div>)}</div>
    <div className={styles.loadingPanel}><div className={styles.loadingPanelHeader}><span/><span/></div>{[1, 2, 3, 4, 5].map(item => <div className={styles.loadingRow} key={item}><span/><span/><span/><span/><span/></div>)}</div>
  </main>;
  if (query.isError && !query.data) return <main className={styles.subscriptions}><section className={styles.errorState} role="alert"><span className={styles.errorIcon}>!</span><h1>Subscriptions couldn’t load</h1><p>{query.error.message || "Check your connection and try again."}</p><button type="button" onClick={() => query.refetch()} disabled={query.isFetching}>{query.isFetching ? "Trying again…" : "Try again"}</button></section></main>;
  return <main className={styles.subscriptions}>
    <SubscriptionsHeader />
    {query.isError && <div className={styles.errorBanner} role="status">Couldn’t refresh subscriptions. Showing the last loaded data. <button onClick={() => query.refetch()}>Retry</button></div>}
    {issue.isError && <div className={styles.errorBanner} role="alert">Could not issue invoice: {issue.error.message}<button onClick={() => issue.reset()}>Dismiss</button></div>}
    <SubscriptionSummary subscriptions={subscriptions} />
    <SubscriptionTable subscriptions={subscriptions} onIssueInvoice={id => issue.mutate(id)} issuing={issue.isPending} />
  </main>;
}
