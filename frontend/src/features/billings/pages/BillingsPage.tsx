import { useMemo, useState } from "react";
import { useAuth } from "@clerk/react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { ArrowUpRight, Check, CircleCheck, Clock3, FilePlus2, MoreHorizontal, Printer, Search, X } from "lucide-react";
import { billingRequest, loadBilling, type ApiInvoice } from "../../billing/api";
import styles from "./BillingsPage.module.css";

type Invoice = { id: number; number: string; customer: string; email: string; plan: string; date: string; due: string; amount: number; status: string; subscriptionId: number; currency: string };
const dateLabel = (value?: string) => value ? new Intl.DateTimeFormat("en-US", { dateStyle: "medium" }).format(new Date(value)) : "—";
const money = (n: number, currency: string) => new Intl.NumberFormat("en-US", { style: "currency", currency }).format(n);
const tabs = ["All invoices", "Open", "Paid", "Draft", "Void"];

export function BillingsPage() {
  const { getToken } = useAuth();
  const client = useQueryClient();
  const [tab, setTab] = useState("All invoices");
  const [query, setQuery] = useState("");
  const [notice, setNotice] = useState("");
  const [menu, setMenu] = useState<number | null>(null);
  const data = useQuery({ queryKey: ["billing-data"], queryFn: () => loadBilling(getToken) });
  const refresh = () => client.invalidateQueries({ queryKey: ["billing-data"] });
  const pay = useMutation({ mutationFn: (id: number) => billingRequest<ApiInvoice>(`/invoices/${id}/pay`, getToken, { method: "POST" }), onSuccess: invoice => { setNotice(`${invoice.invoiceNumber} marked paid.`); refresh(); } });
  const invoices: Invoice[] = (data.data ?? []).flatMap(({ subscription, invoices, variant, product }) => invoices.map(invoice => ({
    id: invoice.id, number: invoice.invoiceNumber, customer: `Customer #${subscription.customerId}`, email: "", plan: variant ? `${product?.name ?? "Product"} · ${variant.name}` : `Variant #${subscription.variantId}`,
    date: dateLabel(invoice.issuedAt ?? invoice.periodStart), due: dateLabel(invoice.dueAt), amount: invoice.totalCents / 100,
    status: invoice.status, subscriptionId: subscription.id, currency: subscription.currency || "USD",
  })));
  const visible = useMemo(() => invoices.filter(invoice => (tab === "All invoices" || invoice.status.toLowerCase() === tab.toLowerCase()) && `${invoice.number} ${invoice.customer} ${invoice.subscriptionId}`.toLowerCase().includes(query.toLowerCase())), [invoices, tab, query]);
  const billed = invoices.filter(item => item.status !== "void").reduce((sum, item) => sum + item.amount, 0);
  const paid = invoices.filter(item => item.status === "paid").reduce((sum, item) => sum + item.amount, 0);
  if (data.isLoading) return <main className={styles.page} aria-busy="true" aria-label="Loading invoices">
    <div className={styles.loadingHeader}><span/><span/><span/></div><div className={styles.loadingRun}/><div className={styles.metrics}>{[1, 2, 3, 4].map(item => <div className={styles.loadingMetric} key={item}><span/><span/><span/></div>)}</div>
    <div className={styles.loadingInvoicePanel}><div className={styles.loadingPanelHeading}><span/><span/></div><div className={styles.loadingToolbar}><span/><span/></div>{[1, 2, 3, 4, 5].map(item => <div className={styles.loadingInvoiceRow} key={item}>{[1, 2, 3, 4, 5, 6].map(cell => <span key={cell}/>)}</div>)}</div>
  </main>;
  if (data.isError && !data.data) return <main className={styles.page}><section className={styles.errorState} role="alert"><span className={styles.errorIcon}>!</span><h1>Invoices couldn’t load</h1><p>{data.error.message || "Check your connection and try again."}</p><button type="button" onClick={() => data.refetch()} disabled={data.isFetching}>{data.isFetching ? "Trying again…" : "Try again"}</button></section></main>;
  return <main className={styles.page}>
    <header className={styles.header}><div><div className={styles.breadcrumb}>Workspace <span>/</span> Billing</div><h1>Billing &amp; invoices</h1><p>Track invoices generated from your subscriptions.</p></div>
      <div className={styles.headerActions}><button className={styles.primaryButton} disabled={data.isFetching} onClick={() => refresh()}>{data.isFetching ? <span className={styles.spinner}/> : <FilePlus2 size={17}/>} {data.isFetching ? "Refreshing…" : "Refresh invoices"}</button></div></header>
    {data.isError && <div className={styles.errorBanner} role="status">Couldn’t refresh invoices. Showing the last loaded data. <button onClick={() => data.refetch()}>Retry</button></div>}
    <section className={styles.runCard}><div className={styles.runIcon}><Clock3 size={19}/></div><div className={styles.runCopy}><div className={styles.runTitle}>Billing records <span className={styles.live}><i/> Connected</span></div><p>Invoices and payment status are loaded from the billing API. Issue invoices from a subscription when its current billing period is ready.</p></div><div className={styles.runAction}><span>{invoices.length} <small>invoices</small></span></div></section>
    <section className={styles.metrics}>
      <article className={styles.metric}><div><span>Total invoiced</span><span className={styles.metricIcon}><ArrowUpRight size={16}/></span></div><strong>{money(billed, invoices[0]?.currency ?? "USD")}</strong><small>All non-void invoices</small></article>
      <article className={styles.metric}><div><span>Collected</span><span className={styles.metricIcon}><Check size={16}/></span></div><strong>{money(paid, invoices[0]?.currency ?? "USD")}</strong><small>{invoices.filter(item => item.status === "paid").length} paid invoices</small></article>
      <article className={styles.metric}><div><span>Awaiting payment</span><span className={styles.metricIcon}><Clock3 size={16}/></span></div><strong>{money(invoices.filter(item => item.status === "open").reduce((sum, item) => sum + item.amount, 0), invoices[0]?.currency ?? "USD")}</strong><small>{invoices.filter(item => item.status === "open").length} open invoices</small></article>
      <article className={styles.metric}><div><span>Draft invoices</span><span className={styles.metricIcon}><FilePlus2 size={16}/></span></div><strong>{invoices.filter(item => item.status === "draft").length}</strong><small>Ready for your review</small></article>
    </section>
    <section className={styles.invoicePanel}><div className={styles.panelHeader}><div><h2>Invoices <span className={styles.count}>{invoices.length}</span></h2><p>Every charge, tied back to its subscription.</p></div></div>
      <div className={styles.toolbar}><div className={styles.tabs}>{tabs.map(item => <button key={item} onClick={() => setTab(item)} className={tab === item ? styles.activeTab : ""}>{item}</button>)}</div><label className={styles.search}><Search size={16}/><input value={query} onChange={event => setQuery(event.target.value)} placeholder="Search invoices…"/></label></div>
      {(notice || pay.isError) && <div className={`${styles.notice} ${pay.isError ? styles.errorNotice : ""}`} role={pay.isError ? "alert" : "status"}>{pay.isError ? <X size={15}/> : <CircleCheck size={15}/>} {notice || pay.error?.message}<button onClick={() => { setNotice(""); pay.reset(); }} aria-label="Dismiss"><X size={14}/></button></div>}
      <div className={styles.tableScroll}><table><thead><tr><th>INVOICE</th><th>CUSTOMER</th><th>SUBSCRIPTION</th><th>ISSUED</th><th>AMOUNT</th><th>STATUS</th><th/></tr></thead><tbody>
        {visible.map(invoice => <tr key={invoice.id}><td><span className={styles.invoiceId}>{invoice.number}</span><span className={styles.due}>Due {invoice.due}</span></td><td><div className={styles.customer}><span className={styles.avatar}>C</span><span><b>{invoice.customer}</b><small>{invoice.email || "Customer details unavailable"}</small></span></div></td><td><span className={styles.planName}>{invoice.plan}</span><small className={styles.subId}>SUB-{invoice.subscriptionId}</small></td><td className={styles.date}>{invoice.date}</td><td><b className={styles.amount}>{money(invoice.amount, invoice.currency)}</b><small className={styles.taxLine}>Tax breakdown unavailable</small></td><td><span className={`${styles.status} ${styles[invoice.status === "open" ? "confirmed" : invoice.status]}`}><i/>{invoice.status[0]?.toUpperCase()}{invoice.status.slice(1)}</span></td><td className={styles.actionsCell}><div className={styles.actions}>{invoice.status === "open" && <button className={styles.confirmButton} disabled={pay.isPending} onClick={() => pay.mutate(invoice.id)}><Check size={14}/>{pay.isPending ? "Saving…" : "Mark paid"}</button>}<button className={styles.moreButton} aria-label={`More actions for ${invoice.number}`} onClick={() => setMenu(menu === invoice.id ? null : invoice.id)}><MoreHorizontal size={18}/></button>{menu === invoice.id && <div className={styles.menu}><button onClick={() => window.print()}><Printer size={14}/> Print invoice</button></div>}</div></td></tr>)}
        {visible.length === 0 && <tr><td colSpan={7} className={styles.empty}>No invoices match this view.</td></tr>}
      </tbody></table></div><footer className={styles.tableFooter}><span>Showing <b>{visible.length}</b> of <b>{invoices.length}</b> invoices</span></footer>
    </section>
    <div className={styles.taxNote}><span><CircleCheck size={15}/></span><div><b>Invoice actions reflect backend support.</b> Invoice issuing and payment recording are connected. Email delivery, tax calculation, draft confirmation, and cancellation are not available in the API.</div></div>
  </main>;
}
