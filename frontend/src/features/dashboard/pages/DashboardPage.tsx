import { Plus } from "lucide-react";
import styles from "./DashboardPage.module.css";

export function DashboardPage() {
  return (
    <div className={styles.dashboard}>
      <div className={styles.headingRow}>
        <div>
          <p className={styles.eyebrow}>Financial pulse · Live synced</p>
          <h1>Subscription Performance</h1>
          <p className={styles.subtitle}>
            Recurring revenue, churn velocity, and subscriber health at a
            glance.
          </p>
        </div>
        <div className={styles.headingActions}>
          <button className={styles.secondaryButton}>Export report</button>
          <button className={styles.primaryButton}>
            <Plus size={15} /> New subscription
          </button>
        </div>
      </div>
    </div>
  );
}
