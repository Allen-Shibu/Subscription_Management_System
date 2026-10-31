import styles from "./BillingsPage.module.css";

export function BillingsPage() {
  return (
    <div className={styles.billings}>
      <header className={styles.header}>
        <div>
          <p className={styles.eyebrow}>Billing & Invoices</p>
          <h1>Invoices & Payments</h1>
          <p className={styles.subtitle}>
            Track upcoming invoices, collected payments, and overdue balances.
          </p>
        </div>
      </header>
      <section className={styles.content}>
        <div className={styles.panel}>
          <h2>Overview</h2>
          <p>Billing data will appear here.</p>
        </div>
      </section>
    </div>
  );
}
