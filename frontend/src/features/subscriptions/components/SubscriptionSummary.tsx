import { FileText, CircleCheck, Radio, Clock3 } from "lucide-react";
import type { Subscription } from "../types";
import styles from "../pages/SubscriptionsPage.module.css";

type Props = { subscriptions: Subscription[] };

export function SubscriptionSummary({ subscriptions }: Props) {
  const metrics = [
    { label: "All subscriptions", value: subscriptions.length, Icon: FileText, tone: "purple" },
    { label: "Awaiting response", value: subscriptions.filter(item => item.status === "Quotation").length, Icon: Clock3, tone: "amber" },
    { label: "Active", value: subscriptions.filter(item => item.status === "Active").length, Icon: Radio, tone: "green" },
    { label: "Confirmed", value: subscriptions.filter(item => item.status === "Confirmed").length, Icon: CircleCheck, tone: "blue" },
  ];
  return <section className={styles.summary} aria-label="Subscription overview">
    {metrics.map(({ label, value, Icon, tone }) => <article className={styles.summaryCard} key={label}>
      <span className={`${styles.summaryIcon} ${styles[tone]}`}><Icon size={17} /></span>
      <div><span>{label}</span><strong>{value}</strong></div>
    </article>)}
  </section>;
}
