import { useEffect, useState } from "react";
import { useAuth } from "@clerk/react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { Link, useLocation } from "react-router-dom";
import { Pencil, Plus, Search, Trash2 } from "lucide-react";
import { catalog } from "./api";
import { ConfirmDelete } from "./ListControls";
import { useListControls } from "./useListControls";
import { planAvailability } from "./validation";
import type { Page, Plan, Product } from "./types";
import styles from "./Catalog.module.css";

const currency = (value: number) => new Intl.NumberFormat("en-US", { style: "currency", currency: "USD" }).format(value);
const today = () => {
  const value = new Date();
  return `${value.getFullYear()}-${String(value.getMonth() + 1).padStart(2, "0")}-${String(value.getDate()).padStart(2, "0")}`;
};

export function CatalogList({ kind }: { kind: "products" | "plans" }) {
  const { getToken } = useAuth();
  const queryClient = useQueryClient();
  const { params, search, setSearch, setPage, setPageSize } = useListControls();
  const [deleting, setDeleting] = useState<Product | Plan | null>(null);
  const [actionError, setActionError] = useState("");
  const location = useLocation();
  const [toast, setToast] = useState<string | null>((location.state as { toast?: string } | null)?.toast ?? null);
  useEffect(() => {
    if (!toast) return;
    const timer = setTimeout(() => setToast(null), 4000);
    return () => clearTimeout(timer);
  }, [toast]);
  const isProduct = kind === "products";
  const singular = isProduct ? "product" : "plan";
  const list = useQuery<Page<Product | Plan>>({
    queryKey: [kind, params.q, params.page, params.pageSize],
    queryFn: ({ signal }) => isProduct ? catalog.listProducts(params, getToken, signal) : catalog.listPlans(params, getToken, signal),
    retry: false,
  });
  const remove = useMutation({
    mutationFn: (id: string) => isProduct ? catalog.deleteProduct(id, getToken) : catalog.deletePlan(id, getToken),
    onSuccess: async () => {
      setDeleting(null);
      setActionError("");
      await queryClient.invalidateQueries({ queryKey: [kind] });
      if (list.data?.items.length === 1 && params.page > 1) setPage(params.page - 1);
      setToast(`${isProduct ? "Product" : "Plan"} deleted`);
    },
    onError: (error: Error) => { setDeleting(null); setActionError(error.message); },
  });
  const totalPages = Math.max(1, Math.ceil((list.data?.total ?? 0) / params.pageSize));

  return <div className={styles.page}>
    <header className={styles.header}><div><p className={styles.eyebrow}>Admin</p><h1>{isProduct ? "Products" : "Plans"}</h1>
      <p>Manage {kind} and their billing details.</p></div>
      <Link className={styles.primary} to={`/admin/${kind}/new`}><Plus size={16} /> Add {singular}</Link></header>
    {toast && <div className={styles.toast} role="status"><span>{toast}</span><button onClick={() => setToast(null)} aria-label="Dismiss notification">×</button></div>}
    {actionError && <div className={styles.error} role="alert">{actionError}</div>}
    <div className={styles.toolbar}>
      <label className={styles.search}><Search size={16} /><input aria-label={`Search ${kind}`} placeholder={`Search ${kind}...`} value={search} onChange={event => setSearch(event.target.value)} /></label>
      <label className={styles.pageSize}>Rows per page <select value={params.pageSize} onChange={event => setPageSize(Number(event.target.value))}>{[5, 10, 20].map(size => <option key={size} value={size}>{size}</option>)}</select></label>
    </div>
    <div className={styles.tableWrap}>
      {list.isPending ? <div className={styles.skeleton} aria-label="Loading" role="status">{Array.from({ length: 5 }, (_, index) => <span key={index} />)}</div> :
        list.isError ? <div className={styles.empty} role="alert"><h2>Could not load {kind}</h2><p>{list.error.message}</p><button onClick={() => list.refetch()}>Retry</button></div> :
        !list.data.items.length ? <div className={styles.empty}><h2>{params.q ? `No ${kind} match your search` : `No ${kind} yet`}</h2><p>{params.q ? "Try a different search." : `Add your first ${singular} to get started.`}</p></div> :
        <table className={styles.table}><thead><tr>{(isProduct ?
          ["Name", "Type", "Sales price", "Cost price", "Margin", "Variants", "Actions"] :
          ["Name", "Price", "Billing period", "Min qty", "Available from", "Available to", "Availability", "Flags", "Actions"]).map(label => <th key={label}>{label}</th>)}</tr></thead>
          <tbody>{list.data.items.map(item => <tr key={item.id}>
            {isProduct && "salesPrice" in item ? <>
              <td className={styles.name}>{item.name}</td><td><span className={styles.badge}>{item.type}</span></td>
              <td>{currency(item.salesPrice)}</td><td>{currency(item.costPrice)}</td><td>{currency(item.salesPrice - item.costPrice)}</td>
              <td>{item.variants.length}</td>
            </> : "billingPeriod" in item ? <>
              <td className={styles.name}>{item.name}</td><td>{currency(item.price)}</td><td><span className={styles.badge}>{item.billingPeriod}</span></td>
              <td>{item.minQuantity}</td><td>{item.availableFrom || "—"}</td><td>{item.availableTo || "—"}</td>
              <td><span className={styles.badge}>{planAvailability(item, today())}</span></td>
              <td><div className={styles.flags}>{[
                ["Auto-close", item.autoClose], ["Early close", item.earlyClose], ["Pause", item.allowPause], ["Renew", item.allowRenew],
              ].filter(([, enabled]) => enabled).map(([label]) => <span className={styles.badge} key={String(label)}>{label}</span>)}</div></td>
            </> : null}
            <td><div className={styles.actions}><Link aria-label={`Edit ${item.name}`} to={`/admin/${kind}/${item.id}/edit`}><Pencil size={16} /></Link>
              <button aria-label={`Delete ${item.name}`} onClick={() => setDeleting(item)}><Trash2 size={16} /></button></div></td>
          </tr>)}</tbody></table>}
    </div>
    <div className={styles.pagination}><span>{list.data?.total ?? 0} total</span><div>
      <button disabled={params.page <= 1} onClick={() => setPage(params.page - 1)}>Previous</button>
      <span>Page {params.page} of {totalPages}</span>
      <button disabled={params.page >= totalPages} onClick={() => setPage(params.page + 1)}>Next</button>
    </div></div>
    {deleting && <ConfirmDelete name={deleting.name} pending={remove.isPending} onCancel={() => setDeleting(null)} onConfirm={() => remove.mutate(deleting.id)} />}
  </div>;
}
