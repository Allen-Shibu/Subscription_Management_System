import { Fragment, useState } from "react";
import { ChevronDown, ChevronRight, FileText, Search } from "lucide-react";
import type { Subscription, SubscriptionStatus } from "../types";
import styles from "../pages/SubscriptionsPage.module.css";

const statuses: Array<"All" | SubscriptionStatus> = ["All", "Draft", "Quotation", "Confirmed", "Active", "Closed"];
const money = (amount: number) => new Intl.NumberFormat("en-US", { style: "currency", currency: "USD", maximumFractionDigits: 2 }).format(amount);
const subtotal = (item: Subscription) => item.lines.reduce((sum, line) => sum + line.quantity * line.unitPrice, 0);
const total = (item: Subscription) => item.lines.reduce((sum, line) => sum + line.quantity * line.unitPrice * (1 + line.taxRate / 100), 0);

type Props = { subscriptions: Subscription[] };

export function SubscriptionTable({ subscriptions }: Props) {
  const [query, setQuery] = useState("");
  const [filter, setFilter] = useState<(typeof statuses)[number]>("All");
  const [expandedId, setExpandedId] = useState<string | null>(null);
  const filtered = subscriptions.filter(item => {
    const matchesStatus = filter === "All" || item.status === filter;
    const search = `${item.number} ${item.customer} ${item.plan}`.toLowerCase();
    return matchesStatus && search.includes(query.toLowerCase().trim());
  });

  return <section className={styles.panel} aria-label="All subscriptions">
    <div className={styles.panelTop}>
      <div><h2>All subscriptions</h2><p>Track offers, contracts, and renewals in one place.</p></div>
      <label className={styles.searchBox}><Search size={16} /><input value={query} onChange={event => setQuery(event.target.value)} placeholder="Search subscriptions" aria-label="Search subscriptions" /></label>
    </div>
    <div className={styles.filters} role="tablist" aria-label="Filter by status">
      {statuses.map(status => <button key={status} type="button" role="tab" aria-selected={filter === status} onClick={() => setFilter(status)} className={`${styles.filterButton} ${filter === status ? styles.filterActive : ""}`}>{status}</button>)}
    </div>
    <div className={styles.tableWrap}>
      <table className={styles.table}>
        <thead><tr><th>Subscription</th><th>Customer</th><th>Plan & period</th><th>Amount</th><th>Status</th><th aria-label="Expand details" /></tr></thead>
        <tbody>{filtered.map(item => <Fragment key={item.id}>
          <tr key={item.id} className={expandedId === item.id ? styles.rowExpanded : ""}>
            <td><div className={styles.numberCell}><span className={styles.fileIcon}><FileText size={16} /></span><div><strong>{item.number}</strong>{item.sourceTemplate && <small>From {item.sourceTemplate}</small>}</div></div></td>
            <td><div className={styles.customerCell}><span className={styles.avatar}>{item.initials}</span><span>{item.customer}</span></div></td>
            <td><div className={styles.planCell}><strong>{item.plan}</strong><small>{item.startDate} – {item.expiryDate}</small></div></td>
            <td><div className={styles.amountCell}><strong>{money(total(item))}</strong><small>incl. taxes</small></div></td>
            <td><span className={`${styles.statusBadge} ${styles[`status${item.status}`]}`}><i />{item.status}</span></td>
            <td><button className={styles.expandButton} type="button" aria-label={`${expandedId === item.id ? "Hide" : "Show"} ${item.number} lines`} aria-expanded={expandedId === item.id} onClick={() => setExpandedId(expandedId === item.id ? null : item.id)}>{expandedId === item.id ? <ChevronDown size={17} /> : <ChevronRight size={17} />}</button></td>
          </tr>
          {expandedId === item.id && <tr key={`${item.id}-details`} className={styles.detailsRow}><td colSpan={6}>
            <div className={styles.detailsGrid}>
              <div className={styles.lineDetails}><div className={styles.detailHeading}><strong>Product lines</strong><span>{item.paymentTerms} payment terms</span></div>
                <div className={styles.lineHeader}><span>Product</span><span>Qty</span><span>Unit price</span><span>Tax</span><span>Amount</span></div>
                {item.lines.map(line => <div className={styles.lineRow} key={line.product}><span>{line.product}</span><span>{line.quantity}</span><span>{money(line.unitPrice)}</span><span>{line.taxRate}%</span><strong>{money(line.quantity * line.unitPrice * (1 + line.taxRate / 100))}</strong></div>)}
              </div>
              <div className={styles.totals}><div><span>Subtotal</span><strong>{money(subtotal(item))}</strong></div><div><span>Taxes</span><strong>{money(total(item) - subtotal(item))}</strong></div><div className={styles.grandTotal}><span>Total</span><strong>{money(total(item))}</strong></div></div>
            </div>
          </td></tr>}
        </Fragment>)}</tbody>
      </table>
      {filtered.length === 0 && <div className={styles.emptyState}>No subscriptions match your search.</div>}
    </div>
    <footer className={styles.tableFooter}><span>Showing <strong>{filtered.length}</strong> of <strong>{subscriptions.length}</strong> subscriptions</span><span>Amounts shown in USD</span></footer>
  </section>;
}
