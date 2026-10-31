import styles from "./SubscriptionsPage.module.css";

export function SubscriptionsPage() {
  return (
    <div className={styles.subscriptions}>
      <header className={styles.header}>
        <div>
          <p className={styles.eyebrow}>Subscriptions</p>
          <h1>Subscription Plans & Contracts</h1>
          <p className={styles.subtitle}>
            View active, trialing, paused, and canceled subscriptions in one place.
          </p>
        </div>
      </header>
      <section className={styles.content}>
        <div className={styles.panel}>
          <h2>All Subscriptions</h2>
          <p>Subscription list will appear here.</p>
        </div>
      </section>
    </div>
  );
}
