import styles from "../pages/SubscriptionsPage.module.css";

export function SubscriptionsHeader() {
  return (
    <header className={styles.header}>
      <div>
        <p className={styles.eyebrow}>Revenue workspace</p>
        <h1>Subscriptions</h1>
        <p className={styles.subtitle}>Review subscription status, billing periods, and issue invoices for active subscriptions.</p>
      </div>
    </header>
  );
}
