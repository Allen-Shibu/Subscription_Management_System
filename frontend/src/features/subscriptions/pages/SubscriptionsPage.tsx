import { subscriptions } from "../data";
import { SubscriptionSummary } from "../components/SubscriptionSummary";
import { SubscriptionsHeader } from "../components/SubscriptionsHeader";
import { SubscriptionTable } from "../components/SubscriptionTable";
import styles from "./SubscriptionsPage.module.css";

export function SubscriptionsPage() {
  return <main className={styles.subscriptions}>
    <SubscriptionsHeader />
    <SubscriptionSummary subscriptions={subscriptions} />
    <SubscriptionTable subscriptions={subscriptions} />
  </main>;
}
