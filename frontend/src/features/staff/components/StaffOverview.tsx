import { ShieldCheck, UserRound, Users } from "lucide-react";
import styles from "../Staff.module.css";

type StaffOverviewProps = {
  total: number;
  admins: number;
  staff: number;
};

export function StaffOverview({ total, admins, staff }: StaffOverviewProps) {
  return (
    <section className={styles.summary} aria-label="Staff overview">
      <article className={`${styles.summaryCard} ${styles.summaryTotal}`}>
        <span className={styles.summaryIcon}><Users size={18} /></span>
        <div><span>Team members</span><strong>{total}</strong></div>
      </article>
      <article className={`${styles.summaryCard} ${styles.summaryAdmins}`}>
        <span className={styles.summaryIcon}><ShieldCheck size={18} /></span>
        <div><span>Admins</span><strong>{admins}</strong></div>
      </article>
      <article className={`${styles.summaryCard} ${styles.summaryStaff}`}>
        <span className={styles.summaryIcon}><UserRound size={18} /></span>
        <div><span>Staff</span><strong>{staff}</strong></div>
      </article>
    </section>
  );
}
