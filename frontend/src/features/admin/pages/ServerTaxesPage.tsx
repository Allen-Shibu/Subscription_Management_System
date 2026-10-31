import { useState, type FormEvent } from "react";
import { useAuth } from "@clerk/react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { Pencil, Plus, Trash2 } from "lucide-react";
import { request } from "../catalog/api";
import { ConfirmDelete } from "../catalog/ListControls";
import styles from "../catalog/Catalog.module.css";

type Tax = { id: number; name: string; rateBps: number; active: boolean };
type Fields = Omit<Tax, "id">;
const empty: Fields = { name: "", rateBps: 0, active: true };

export function ServerTaxesPage() {
  const { getToken } = useAuth();
  const queryClient = useQueryClient();
  const [editing, setEditing] = useState<(Fields & { id?: number }) | null>(null);
  const [deleting, setDeleting] = useState<Tax | null>(null);
  const [error, setError] = useState("");
  const taxes = useQuery({ queryKey: ["server-taxes"], queryFn: () => request<Tax[]>("/tax-rules", getToken), retry: false });
  const save = useMutation({
    mutationFn: (input: Fields & { id?: number }) => request<Tax>(`/admin/tax-rules${input.id ? `/${input.id}` : ""}`, getToken, { method: input.id ? "PUT" : "POST", body: JSON.stringify(input) }),
    onSuccess: async () => { setEditing(null); await queryClient.invalidateQueries({ queryKey: ["server-taxes"] }); },
    onError: (failure: Error) => setError(failure.message),
  });
  const remove = useMutation({
    mutationFn: (id: number) => request<void>(`/admin/tax-rules/${id}`, getToken, { method: "DELETE" }),
    onSuccess: async () => { setDeleting(null); await queryClient.invalidateQueries({ queryKey: ["server-taxes"] }); },
    onError: (failure: Error) => { setDeleting(null); setError(failure.message); },
  });
  const submit = (event: FormEvent) => {
    event.preventDefault(); setError("");
    if (!editing || !editing.name.trim() || !Number.isInteger(editing.rateBps) || editing.rateBps < 0 || editing.rateBps > 10000) { setError("Enter a name and a rate between 0% and 100%."); return; }
    save.mutate({ ...editing, name: editing.name.trim() });
  };
  return <div className={styles.page}>
    <header className={styles.header}><div><p className={styles.eyebrow}>Live data</p><h1>Tax rules</h1><p>Rules saved to the Go API. Rate is stored in basis points.</p></div>
      <button className={styles.primary} onClick={() => setEditing({ ...empty })}><Plus size={16} /> Add tax rule</button></header>
    {error && <div className={styles.error} role="alert">{error}</div>}
    <div className={styles.tableWrap}>{taxes.isPending ? <div className={styles.skeleton} role="status" aria-label="Loading tax rules"><span /><span /></div> :
      taxes.isError ? <div className={styles.empty} role="alert"><h2>Could not load tax rules</h2><p>{taxes.error.message}</p><button onClick={() => taxes.refetch()}>Retry</button></div> :
      !taxes.data.length ? <div className={styles.empty}><h2>No tax rules yet</h2></div> :
      <table className={styles.table}><thead><tr><th>Name</th><th>Rate</th><th>Status</th><th>Actions</th></tr></thead><tbody>{taxes.data.map(item =>
        <tr key={item.id}><td className={styles.name}>{item.name}</td><td>{(item.rateBps / 100).toFixed(2)}%</td><td><span className={styles.badge}>{item.active ? "Active" : "Inactive"}</span></td>
          <td><div className={styles.actions}><button aria-label={`Edit ${item.name}`} onClick={() => setEditing({ ...item })}><Pencil size={16} /></button><button aria-label={`Delete ${item.name}`} onClick={() => setDeleting(item)}><Trash2 size={16} /></button></div></td></tr>)}</tbody></table>}
    </div>
    {editing && <form className={styles.form} onSubmit={submit}><section className={styles.formCard}><h2>{editing.id ? "Edit tax rule" : "Add tax rule"}</h2><div className={styles.formGrid}>
      <label>Name<input required value={editing.name} onChange={event => setEditing({ ...editing, name: event.target.value })} /></label>
      <label>Rate (%)<input type="number" min="0" max="100" step="0.01" required value={editing.rateBps / 100} onChange={event => setEditing({ ...editing, rateBps: Math.round(Number(event.target.value) * 100) })} /></label>
      <label><span>Active</span><input type="checkbox" checked={editing.active} onChange={event => setEditing({ ...editing, active: event.target.checked })} /></label>
    </div></section><div className={styles.formActions}><button type="button" onClick={() => setEditing(null)}>Cancel</button><button className={styles.primary} disabled={save.isPending}>{save.isPending ? "Saving..." : "Save tax rule"}</button></div></form>}
    {deleting && <ConfirmDelete name={deleting.name} pending={remove.isPending} onCancel={() => setDeleting(null)} onConfirm={() => remove.mutate(deleting.id)} />}
  </div>;
}
