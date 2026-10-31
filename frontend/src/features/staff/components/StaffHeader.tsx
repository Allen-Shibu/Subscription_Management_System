import { Plus } from "lucide-react";
import styles from "../Staff.module.css";

type StaffHeaderProps = {
  onAddStaff: () => void;
};

export function StaffHeader({ onAddStaff }: StaffHeaderProps) {
  return (
    <header className={styles.header}>
      <div>
        <p className={styles.eyebrow}>Staff Management</p>
        <h1>Team Members</h1>
        <p className={styles.subtitle}>
          Preview staff members and roles.
        </p>
      </div>
      <div className={styles.headerActions}>
        <button className={styles.primaryButton} onClick={onAddStaff}>
          <Plus size={15} /> Add Staff
        </button>
      </div>
    </header>
  );
}
