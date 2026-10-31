import styles from "./CustomersPage.module.css";

export function CustomersPage() {
  return (
    <div className={styles.customers}>
      <header className={styles.header}>
        <div>
          <p className={styles.eyebrow}>Customers</p>
          <h1>Customer Directory</h1>
          <p className={styles.subtitle}>
            Manage accounts, contact details, and subscription history.
          </p>
        </div>
      </header>
      <section className={styles.content}>
        <div className={styles.panel}>
          <h2>Directory</h2>
          <p>Customer records will appear here.</p>
        </div>
      </section>
    </div>
  );
}
