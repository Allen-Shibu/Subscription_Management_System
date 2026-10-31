import type { Dispatch, FormEvent, SetStateAction } from "react";
import type { StaffForm } from "../types";
import styles from "../Staff.module.css";

type StaffModalProps = {
  mode: "create" | "edit";
  form: StaffForm;
  setForm: Dispatch<SetStateAction<StaffForm>>;
  onClose: () => void;
  onSubmit: (event: FormEvent<HTMLFormElement>) => void;
};

export function StaffModal({ mode, form, setForm, onClose, onSubmit }: StaffModalProps) {
  return (
    <div className={styles.modalBackdrop} onClick={onClose}>
      <div
        className={styles.modal}
        role="dialog"
        aria-modal="true"
        onClick={(event) => event.stopPropagation()}
      >
        <header className={styles.modalHeader}>
          <h2>{mode === "create" ? "Add New Staff" : "Edit Staff"}</h2>
          <button className={styles.closeButton} onClick={onClose} aria-label="Close">
            ×
          </button>
        </header>
        <form className={styles.modalForm} onSubmit={onSubmit}>
          <label className={styles.field}>
            <span>Staff Name</span>
            <input
              type="text"
              placeholder="John Doe"
              value={form.name}
              onChange={(event) => setForm((current) => ({ ...current, name: event.target.value }))}
              required
            />
          </label>
          <label className={styles.field}>
            <span>Email</span>
            <input
              type="email"
              placeholder="johndoe@odoo.com"
              value={form.email}
              onChange={(event) => setForm((current) => ({ ...current, email: event.target.value }))}
              required
            />
          </label>
          <label className={styles.checkboxField}>
            <input
              type="checkbox"
              checked={form.isAdmin}
              onChange={(event) => setForm((current) => ({ ...current, isAdmin: event.target.checked }))}
            />
            <span>Admin role (demo only)</span>
          </label>
          <div className={styles.modalActions}>
            <button type="button" className={styles.secondaryButton} onClick={onClose}>
              Cancel
            </button>
            <button type="submit" className={styles.primaryButton}>
              {mode === "create" ? "Add demo staff" : "Save demo changes"}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
