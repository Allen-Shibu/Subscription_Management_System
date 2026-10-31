import { Plus } from "lucide-react";
import styles from "../pages/SubscriptionsPage.module.css";

export function SubscriptionsHeader() {
  return (
    <header className={styles.header}>
      <div>
        <p className={styles.eyebrow}>Revenue workspace</p>
        <h1>Subscriptions</h1>
        <p className={styles.subtitle}>Prepare offers from quotation templates, then follow every subscription through to close.</p>
      </div>
      <button className={styles.primaryButton} type="button"><Plus size={16} /> New subscription</button>
    </header>
  );
}
