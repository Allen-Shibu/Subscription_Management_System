import { useEffect, useState } from "react";
import styles from "./Catalog.module.css";

export function ConfirmDelete({ name, pending, onCancel, onConfirm }: { name: string; pending: boolean; onCancel: () => void; onConfirm: () => void }) {
  const [dialog, setDialog] = useState<HTMLDialogElement | null>(null);
  useEffect(() => {
    if (dialog) dialog.showModal();
    return () => { if (dialog?.open) dialog.close(); };
  }, [dialog]);
  return <dialog className={styles.dialog} ref={setDialog} onCancel={event => { if (pending) event.preventDefault(); else onCancel(); }} aria-labelledby="delete-title">
    <h2 id="delete-title">Delete {name}?</h2>
    <p>This action cannot be undone.</p>
    <div className={styles.dialogActions}><button disabled={pending} onClick={onCancel}>Cancel</button><button disabled={pending} className={styles.dangerButton} onClick={onConfirm}>{pending ? "Deleting..." : "Delete"}</button></div>
  </dialog>;
}
