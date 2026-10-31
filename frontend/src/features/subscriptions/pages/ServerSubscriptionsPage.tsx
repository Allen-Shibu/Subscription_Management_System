import { useState } from "react";
import { useAuth } from "@clerk/react";
import { useQuery } from "@tanstack/react-query";
import { request } from "../../admin/catalog/api";
import styles from "../../admin/catalog/Catalog.module.css";

type Subscription = { id: number; customerId: number; variantId: number; status: string; quantity: number; unitPriceCents: number; currency: string; billingInterval: string; currentPeriodEnd: string };
const money = (cents: number, currency: string) => new Intl.NumberFormat("en-US", { style: "currency", currency }).format(cents / 100);

export function ServerSubscriptionsPage() {
  const { getToken } = useAuth();
  const [status, setStatus] = useState("");
  const subscriptions = useQuery({ queryKey: ["server-subscriptions", status], queryFn: () => request<Subscription[]>(`/staff/subscriptions${status ? `?status=${status}` : ""}`, getToken), retry: false });
  return <div className={styles.page}>
    <header className={styles.header}><div><p className={styles.eyebrow}>Live data</p><h1>Subscriptions</h1><p>Subscriptions available to staff and admins.</p></div></header>
    <div className={styles.toolbar}><label className={styles.pageSize}>Status <select value={status} onChange={event => setStatus(event.target.value)}><option value="">All statuses</option>{["active", "trialing", "past_due", "canceled", "expired"].map(value => <option key={value} value={value}>{value.replaceAll("_", " ")}</option>)}</select></label></div>
    <div className={styles.tableWrap}>{subscriptions.isPending ? <div className={styles.skeleton} role="status" aria-label="Loading subscriptions"><span /><span /><span /></div> :
      subscriptions.isError ? <div className={styles.empty} role="alert"><h2>Could not load subscriptions</h2><p>{subscriptions.error.message}</p><button onClick={() => subscriptions.refetch()}>Retry</button></div> :
      !subscriptions.data.length ? <div className={styles.empty}><h2>No subscriptions yet</h2></div> :
      <table className={styles.table}><thead><tr><th>ID</th><th>Customer ID</th><th>Variant ID</th><th>Status</th><th>Quantity</th><th>Price</th><th>Period ends</th></tr></thead><tbody>{subscriptions.data.map(item =>
        <tr key={item.id}><td>{item.id}</td><td>{item.customerId}</td><td>{item.variantId}</td><td><span className={styles.badge}>{item.status.replaceAll("_", " ")}</span></td><td>{item.quantity}</td><td>{money(item.unitPriceCents, item.currency)} / {item.billingInterval}</td><td>{new Date(item.currentPeriodEnd).toLocaleDateString()}</td></tr>)}</tbody></table>}
    </div>
  </div>;
}
