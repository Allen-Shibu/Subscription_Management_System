import { useMemo, useState } from "react";
import { Link } from "react-router-dom";
import { Area, AreaChart, CartesianGrid, ResponsiveContainer, Tooltip, XAxis, YAxis } from "recharts";
import { Activity, ArrowUpRight, CreditCard, FileWarning, Wallet } from "lucide-react";
import { summarizeDashboard } from "./dashboardSummary";
import styles from "./DashboardPage.module.css";

const iso = (value: Date) =>
  `${value.getFullYear()}-${String(value.getMonth() + 1).padStart(2, "0")}-${String(value.getDate()).padStart(2, "0")}`;
const day = (offset: number) => {
  const value = new Date();
  value.setDate(value.getDate() + offset);
  return iso(value);
};
const monthAgo = (months: number) => iso(new Date(new Date().getFullYear(), new Date().getMonth() - months, 15));
const today = day(0);
const subscriptions = [
  { customer: "Northstar Studio", plan: "Growth", status: "Active", started: day(-89) },
  { customer: "Brightpath Labs", plan: "Enterprise", status: "Active", started: day(-62) },
  { customer: "Meridian Co.", plan: "Starter", status: "Active", started: day(-38) },
  { customer: "Orbit Works", plan: "Growth", status: "Active", started: day(-18) },
  { customer: "Pine & Co.", plan: "Starter", status: "Paused", started: day(-95) },
  { customer: "Acme Digital", plan: "Enterprise", status: "Canceled", started: day(-115) },
];
const invoices = [
  { id: "INV-2041", customer: "Northstar Studio", plan: "Growth", amount: 2400, issued: day(-27), due: day(-10), paid: "" },
  { id: "INV-2042", customer: "Brightpath Labs", plan: "Enterprise", amount: 6200, issued: day(-20), due: day(-5), paid: "" },
  { id: "INV-2043", customer: "Meridian Co.", plan: "Starter", amount: 890, issued: day(-16), due: day(3), paid: "" },
  { id: "INV-2044", customer: "Orbit Works", plan: "Growth", amount: 2400, issued: day(-11), due: day(9), paid: "" },
  { id: "INV-2045", customer: "Northstar Studio", plan: "Growth", amount: 2400, issued: day(-8), due: day(-1), paid: day(-3) },
  { id: "INV-2046", customer: "Brightpath Labs", plan: "Enterprise", amount: 6200, issued: day(-5), due: day(11), paid: day(-2) },
  { id: "INV-2047", customer: "Meridian Co.", plan: "Starter", amount: 890, issued: day(-3), due: day(12), paid: day(-1) },
  { id: "INV-2035", customer: "Northstar Studio", plan: "Growth", amount: 2400, issued: monthAgo(1), due: monthAgo(1), paid: monthAgo(1) },
  { id: "INV-2029", customer: "Brightpath Labs", plan: "Enterprise", amount: 6200, issued: monthAgo(2), due: monthAgo(2), paid: monthAgo(2) },
  { id: "INV-2023", customer: "Meridian Co.", plan: "Starter", amount: 890, issued: monthAgo(3), due: monthAgo(3), paid: monthAgo(3) },
  { id: "INV-2017", customer: "Northstar Studio", plan: "Growth", amount: 2400, issued: monthAgo(4), due: monthAgo(4), paid: monthAgo(4) },
  { id: "INV-2011", customer: "Brightpath Labs", plan: "Enterprise", amount: 6200, issued: monthAgo(5), due: monthAgo(5), paid: monthAgo(5) },
];
const money = (value: number) => new Intl.NumberFormat("en-US", { style: "currency", currency: "USD", maximumFractionDigits: 0 }).format(value);

export function DashboardPage() {
  const [from, setFrom] = useState(day(-180));
  const [to, setTo] = useState(today);
  const [plan, setPlan] = useState("All plans");
  const [customer, setCustomer] = useState("All customers");
  const [status, setStatus] = useState("All statuses");
  const data = useMemo(() => summarizeDashboard(subscriptions, invoices, today, { from, to, plan, customer, status }), [from, to, plan, customer, status]);

  return (
    <div className={styles.dashboard}>
      <header className={styles.header}>
        <div>
          <p className={styles.eyebrow}>Dashboard · Sample data</p>
          <h1>Subscription overview</h1>
          <p className={styles.subtitle}>Revenue, outstanding invoices, and subscriptions at a glance.</p>
        </div>
        <Link className={styles.primaryButton} to="/subscriptions">View subscriptions <ArrowUpRight size={15} /></Link>
      </header>
      <section className={styles.filters} aria-label="Dashboard filters">
        <label>From<input type="date" value={from} max={to} onChange={event => setFrom(event.target.value)} /></label>
        <label>To<input type="date" value={to} min={from} onChange={event => setTo(event.target.value)} /></label>
        <label>Plan<select value={plan} onChange={event => setPlan(event.target.value)}>{["All plans", "Starter", "Growth", "Enterprise"].map(value => <option key={value}>{value}</option>)}</select></label>
        <label>Customer<select value={customer} onChange={event => setCustomer(event.target.value)}>{["All customers", ...new Set(subscriptions.map(item => item.customer))].map(value => <option key={value}>{value}</option>)}</select></label>
        <label>Status<select value={status} onChange={event => setStatus(event.target.value)}>{["All statuses", "Active", "Paused", "Canceled", "Paid", "Unpaid", "Overdue"].map(value => <option key={value}>{value}</option>)}</select></label>
      </section>
      <section className={styles.metricGrid} aria-label="Summary metrics">
        {[
          { label: "Active subscriptions", value: String(data.active), icon: Activity, detail: "Active at end date" },
          { label: "Revenue this month", value: money(data.revenue), icon: CreditCard, detail: "Payments received" },
          { label: "Outstanding amount", value: money(data.outstanding), icon: Wallet, detail: "Invoiced, not yet paid" },
          { label: "Overdue invoices", value: String(data.overdue.length), icon: FileWarning, detail: "Needs follow-up" },
        ].map(metric => <article className={styles.metricCard} key={metric.label}>
          <div className={styles.metricTop}><span>{metric.label}</span><span className={styles.metricIcon}><metric.icon size={18} /></span></div>
          <strong className={styles.metricValue}>{metric.value}</strong><span className={styles.metricDetail}>{metric.detail}</span>
        </article>)}
      </section>
      <section className={styles.contentGrid}>
        <div className={styles.card}>
          <div className={styles.cardHeading}><div><h2>Payments received</h2><p>Monthly total for invoices matching the filters</p></div></div>
          <div className={styles.chart} role="img" aria-label="Monthly payments received chart">
            {data.chart.length ? <ResponsiveContainer width="100%" height="100%"><AreaChart data={data.chart} margin={{ top: 16, right: 12, left: 0, bottom: 0 }}>
              <CartesianGrid stroke="var(--color-border-primary)" vertical={false} />
              <XAxis dataKey="month" tickLine={false} axisLine={false} />
              <YAxis tickFormatter={value => `$${value / 1000}k`} tickLine={false} axisLine={false} />
              <Tooltip formatter={value => money(Number(value))} />
              <Area type="monotone" dataKey="revenue" stroke="var(--color-primary)" fill="var(--color-primary-highlight)" strokeWidth={3} />
            </AreaChart></ResponsiveContainer> : <p className={styles.empty}>No payments match these filters.</p>}
          </div>
        </div>
        <div className={styles.card}>
          <div className={styles.cardHeading}><div><h2>Overdue to chase</h2><p>Unpaid invoices past their due date</p></div></div>
          {data.overdue.length ? <ul className={styles.overdueList}>{data.overdue.map(invoice => <li key={invoice.id}>
            <div><strong>{invoice.customer}</strong><span>{invoice.id} · Due {invoice.due}</span></div><b>{money(invoice.amount)}</b>
          </li>)}</ul> : <p className={styles.empty}>No overdue invoices match these filters.</p>}
        </div>
      </section>
      <p className={styles.sampleNote}>Sample data only. Figures are not connected to the backend.</p>
    </div>
  );
}
