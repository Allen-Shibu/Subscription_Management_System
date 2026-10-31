import { useState, type FormEvent } from "react";
import { useAuth } from "@clerk/react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { Link, useNavigate, useParams } from "react-router-dom";
import { ArrowLeft, Plus, Trash2 } from "lucide-react";
import { ConfirmDelete } from "./ListControls";
import { serverProducts, type ProductFields, type ServerVariant, type VariantFields } from "./serverProducts";
import styles from "./Catalog.module.css";

const emptyProduct: ProductFields = { name: "", description: "", active: true };
const emptyVariant: VariantFields = { name: "", description: "", priceCents: 0, currency: "USD", billingInterval: "month", active: true };
const money = (cents: number, currency: string) => new Intl.NumberFormat("en-US", { style: "currency", currency }).format(cents / 100);

export function ServerProductForm({ readonly = false }: { readonly?: boolean }) {
  const { id } = useParams();
  const { getToken } = useAuth();
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const [fields, setFields] = useState<ProductFields | null>(null);
  const [tier, setTier] = useState<(VariantFields & { id?: number }) | null>(null);
  const [deleting, setDeleting] = useState<ServerVariant | null>(null);
  const [error, setError] = useState("");
  const product = useQuery({ queryKey: ["server-product", id], queryFn: () => serverProducts.get(id!, getToken), enabled: Boolean(id), retry: false });
  const variants = useQuery({ queryKey: ["server-variants", id], queryFn: () => serverProducts.variants(id!, getToken), enabled: Boolean(id), retry: false });
  const currentFields = fields ?? (product.data ? { name: product.data.name, description: product.data.description, active: product.data.active } : emptyProduct);
  const save = useMutation({
    mutationFn: (input: ProductFields) => serverProducts.save(id, input, getToken),
    onSuccess: async result => { await queryClient.invalidateQueries({ queryKey: ["server-products"] }); navigate(`/admin/products/${result.id}/edit`); },
    onError: (failure: Error) => setError(failure.message),
  });
  const saveTier = useMutation({
    mutationFn: (input: VariantFields & { id?: number }) => serverProducts.saveVariant(id!, input.id, input, getToken),
    onSuccess: async () => { setTier(null); await queryClient.invalidateQueries({ queryKey: ["server-variants", id] }); },
    onError: (failure: Error) => setError(failure.message),
  });
  const deleteTier = useMutation({
    mutationFn: (variantId: number) => serverProducts.deleteVariant(variantId, getToken),
    onSuccess: async () => { setDeleting(null); await queryClient.invalidateQueries({ queryKey: ["server-variants", id] }); },
    onError: (failure: Error) => { setDeleting(null); setError(failure.message); },
  });
  const submit = (event: FormEvent) => { event.preventDefault(); setError(""); if (!currentFields.name.trim()) { setError("Name is required"); return; } save.mutate({ ...currentFields, name: currentFields.name.trim() }); };
  const submitTier = (event: FormEvent) => {
    event.preventDefault(); setError("");
    if (!tier || !tier.name.trim() || !Number.isInteger(tier.priceCents) || tier.priceCents < 0 || !/^[A-Z]{3}$/.test(tier.currency)) { setError("Enter a tier name, a nonnegative whole cent price, and a three-letter currency."); return; }
    saveTier.mutate({ ...tier, name: tier.name.trim() });
  };

  if (id && (product.isPending || variants.isPending)) return <div className={styles.skeleton} role="status" aria-label="Loading product"><span /><span /><span /></div>;
  if (id && (product.isError || variants.isError)) return <div className={styles.empty} role="alert"><h2>Could not load product</h2><p>{product.error?.message ?? variants.error?.message}</p><button onClick={() => { product.refetch(); variants.refetch(); }}>Retry</button></div>;
  if (readonly && !product.data?.active) return <div className={styles.empty}><h2>Product unavailable</h2></div>;
  const back = readonly ? "/products" : "/admin/products";
  return <div className={styles.page}>
    <header className={styles.header}><div><Link className={styles.back} to={back}><ArrowLeft size={15} /> Products</Link>
      <h1>{readonly ? product.data?.name : id ? "Edit product" : "Add product"}</h1>
      <p>Live Go catalog · Billing tiers have their own prices and billing periods.</p></div></header>
    {error && <div className={styles.error} role="alert">{error}</div>}
    {readonly ? <section className={styles.formCard}><h2>Product details</h2><p>{product.data?.description || "No description"}</p><p>{product.data?.active ? "Active" : "Inactive"}</p></section> :
      <form className={styles.form} onSubmit={submit}><section className={styles.formCard}><h2>Product details</h2><div className={styles.formGrid}>
        <label>Name<input required maxLength={100} value={currentFields.name} onChange={event => setFields({ ...currentFields, name: event.target.value })} /></label>
        <label className={styles.wide}>Description<textarea rows={4} value={currentFields.description} onChange={event => setFields({ ...currentFields, description: event.target.value })} /></label>
        <label><span>Active</span><input type="checkbox" checked={currentFields.active} onChange={event => setFields({ ...currentFields, active: event.target.checked })} /></label>
      </div></section><div className={styles.formActions}><Link to={back}>Cancel</Link><button className={styles.primary} disabled={save.isPending}>{save.isPending ? "Saving..." : "Save product"}</button></div></form>}
    {id && <section className={styles.formCard}><div className={styles.sectionHead}><div><h2>Billing tiers</h2><p>Each tier is saved separately by the Go API.</p></div>
      {!readonly && <button type="button" onClick={() => setTier({ ...emptyVariant })}><Plus size={15} /> Add tier</button>}</div>
      {variants.data?.length ? <div className={styles.variantList}>{variants.data.map(item => <div className={styles.variant} key={item.id}>
        <strong>{item.name}</strong> · {money(item.priceCents, item.currency)} / {item.billingInterval} · {item.active ? "Active" : "Inactive"}
        {item.description && <p className={styles.preview}>{item.description}</p>}
        {!readonly && <div className={styles.actions}><button type="button" aria-label={`Edit ${item.name}`} onClick={() => setTier({ ...item })}>Edit</button>
          <button type="button" aria-label={`Delete ${item.name}`} onClick={() => setDeleting(item)}><Trash2 size={16} /></button></div>}
      </div>)}</div> : <p className={styles.muted}>No billing tiers yet.</p>}
      {tier && <form className={styles.form} onSubmit={submitTier}><div className={styles.formGrid}>
        <label>Tier name<input required value={tier.name} onChange={event => setTier({ ...tier, name: event.target.value })} /></label>
        <label>Description<input value={tier.description} onChange={event => setTier({ ...tier, description: event.target.value })} /></label>
        <label>Price in cents<input required type="number" min="0" step="1" value={tier.priceCents} onChange={event => setTier({ ...tier, priceCents: Number(event.target.value) })} /></label>
        <label>Currency<input required maxLength={3} value={tier.currency} onChange={event => setTier({ ...tier, currency: event.target.value.toUpperCase() })} /></label>
        <label>Billing period<select value={tier.billingInterval} onChange={event => setTier({ ...tier, billingInterval: event.target.value as "month" | "year" })}><option value="month">Monthly</option><option value="year">Yearly</option></select></label>
        <label><span>Active</span><input type="checkbox" checked={tier.active} onChange={event => setTier({ ...tier, active: event.target.checked })} /></label>
      </div><div className={styles.formActions}><button type="button" onClick={() => setTier(null)}>Cancel</button><button className={styles.primary} disabled={saveTier.isPending}>{saveTier.isPending ? "Saving..." : "Save tier"}</button></div></form>}
    </section>}
    {deleting && <ConfirmDelete name={deleting.name} pending={deleteTier.isPending} onCancel={() => setDeleting(null)} onConfirm={() => deleteTier.mutate(deleting.id)} />}
  </div>;
}
