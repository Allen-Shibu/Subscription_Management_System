import { useState, type FormEvent } from "react";
import { useAuth } from "@clerk/react";
import { useMutation } from "@tanstack/react-query";
import { request } from "../../admin/catalog/api";
import styles from "../../admin/catalog/Catalog.module.css";

type Invitation = { id: string; emailAddress: string; status: string };

export function ServerStaffPage() {
  const { getToken } = useAuth();
  const [email, setEmail] = useState("");
  const [success, setSuccess] = useState("");
  const invite = useMutation({
    mutationFn: (address: string) => request<Invitation>("/admin/staff", getToken, { method: "POST", body: JSON.stringify({ email: address }) }),
    onSuccess: result => { setSuccess(`Invitation sent to ${result.emailAddress}`); setEmail(""); },
  });
  const submit = (event: FormEvent) => { event.preventDefault(); setSuccess(""); invite.mutate(email.trim()); };
  return <div className={styles.page}>
    <header className={styles.header}><div><p className={styles.eyebrow}>Live Clerk invitation</p><h1>Invite staff</h1>
      <p>The Go API sends a Clerk invitation with the staff role. This API does not provide a staff list or role editing.</p></div></header>
    {success && <div role="status" className={styles.badge}>{success}</div>}
    {invite.isError && <div className={styles.error} role="alert">{invite.error.message}</div>}
    <form className={styles.form} onSubmit={submit}><section className={styles.formCard}><div className={styles.formGrid}>
      <label>Email address<input type="email" required autoComplete="email" value={email} onChange={event => setEmail(event.target.value)} /></label>
    </div></section><div className={styles.formActions}><button className={styles.primary} disabled={invite.isPending}>{invite.isPending ? "Sending..." : "Send invitation"}</button></div></form>
  </div>;
}
