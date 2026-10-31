import { useEffect, useState } from "react";
import { useAuth } from "@clerk/react";
import { zodResolver } from "@hookform/resolvers/zod";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useFieldArray, useForm, useWatch, type FieldPath } from "react-hook-form";
import { Link, useNavigate, useParams } from "react-router-dom";
import { ArrowLeft, Plus, Trash2 } from "lucide-react";
import { z } from "zod";
import { catalog } from "./api";
import { CatalogError, type ProductInput } from "./types";
import { productSchema } from "./validation";
import styles from "./Catalog.module.css";

type Form = z.infer<typeof productSchema>;
const empty: Form = { name: "", type: "Service", salesPrice: 0, costPrice: 0, notes: "", variants: [] };
const money = (value: number) => new Intl.NumberFormat("en-US", { style: "currency", currency: "USD" }).format(value);

export function ProductForm() {
  const { id } = useParams();
  const { getToken } = useAuth();
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const [saveError, setSaveError] = useState("");
  const product = useQuery({ queryKey: ["product", id], queryFn: () => catalog.getProduct(id!, getToken), enabled: Boolean(id), retry: false });
  const { register, control, handleSubmit, reset, setError, formState: { errors } } = useForm<Form>({
    resolver: zodResolver(productSchema), defaultValues: empty,
  });
  const { fields, append, remove } = useFieldArray({ control, name: "variants" });
  const values = useWatch({ control });
  useEffect(() => { if (product.data) reset(product.data); }, [product.data, reset]);
  const save = useMutation({
    mutationFn: (input: ProductInput) => catalog.saveProduct(id, input, getToken),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ["products"] });
      navigate("/admin/products", { state: { toast: id ? "Product updated" : "Product created" } });
    },
    onError: (error: Error) => {
      setSaveError(error.message);
      if (error instanceof CatalogError && error.status === 400) {
        Object.entries(error.fields).forEach(([field, message]) => setError(field as FieldPath<Form>, { type: "server", message }));
      }
    },
  });
  const submit = handleSubmit(data => { setSaveError(""); save.mutate(data); });

  if (id && product.isPending) return <div className={styles.skeleton} role="status" aria-label="Loading product"><span /><span /><span /></div>;
  if (id && product.isError) return <div className={styles.empty} role="alert"><h2>Could not load product</h2><p>{product.error.message}</p><button onClick={() => product.refetch()}>Retry</button></div>;
  const base = Number(values.salesPrice) || 0;
  const cost = Number(values.costPrice) || 0;
  return <div className={styles.page}>
    <header className={styles.header}><div><Link className={styles.back} to="/admin/products"><ArrowLeft size={15} /> Products</Link>
      <h1>{id ? "Edit product" : "Add product"}</h1><p>Set pricing and optional variant adjustments.</p></div></header>
    {saveError && <div className={styles.error} role="alert">{saveError}</div>}
    <form className={styles.form} onSubmit={submit} noValidate>
      <section className={styles.formCard}><h2>Product details</h2><div className={styles.formGrid}>
        <label>Name<input {...register("name")} maxLength={100} aria-invalid={!!errors.name} /><small>{errors.name?.message}</small></label>
        <label>Type<select {...register("type")}>{["Service", "Physical", "Digital"].map(type => <option key={type}>{type}</option>)}</select><small>{errors.type?.message}</small></label>
        <label>Sales price<input type="number" min="0" step="0.01" {...register("salesPrice", { valueAsNumber: true })} aria-invalid={!!errors.salesPrice} /><small>{errors.salesPrice?.message}</small></label>
        <label>Cost price<input type="number" min="0" step="0.01" {...register("costPrice", { valueAsNumber: true })} aria-invalid={!!errors.costPrice} /><small>{errors.costPrice?.message}</small></label>
        {cost > base && <p className={styles.warning}>Cost price is higher than sales price; this product has a negative margin.</p>}
        <label className={styles.wide}>Notes<textarea rows={4} {...register("notes")} /><small>{errors.notes?.message}</small></label>
      </div></section>
      <section className={styles.formCard}><div className={styles.sectionHead}><div><h2>Variants</h2><p>Each variant adds to the base sales price.</p></div>
        <button type="button" onClick={() => append({ attribute: "", value: "", extraPrice: 0 })}><Plus size={15} /> Add variant</button></div>
        {fields.length ? <div className={styles.variantList}>{fields.map((field, index) => {
          const variant = values.variants?.[index];
          const extra = Number(variant?.extraPrice) || 0;
          return <div className={styles.variant} key={field.id}>
            <div className={styles.variantFields}>
              <label>Attribute<input placeholder="Brand" {...register(`variants.${index}.attribute`)} aria-invalid={!!errors.variants?.[index]?.attribute} /><small>{errors.variants?.[index]?.attribute?.message}</small></label>
              <label>Value<input placeholder="Odoo" {...register(`variants.${index}.value`)} aria-invalid={!!errors.variants?.[index]?.value} /><small>{errors.variants?.[index]?.value?.message}</small></label>
              <label>Extra price<input type="number" min="0" step="0.01" {...register(`variants.${index}.extraPrice`, { valueAsNumber: true })} aria-invalid={!!errors.variants?.[index]?.extraPrice} /><small>{errors.variants?.[index]?.extraPrice?.message}</small></label>
              <button type="button" className={styles.iconButton} aria-label={`Remove variant ${index + 1}`} onClick={() => remove(index)}><Trash2 size={16} /></button>
            </div>
            <p className={styles.preview}>Base {money(base)}, {variant?.value || "Variant"} +{money(extra)} = <strong>{money(base + extra)}</strong></p>
          </div>;
        })}</div> : <p className={styles.muted}>No variants added.</p>}
      </section>
      <div className={styles.formActions}><Link to="/admin/products">Cancel</Link><button className={styles.primary} disabled={save.isPending}>{save.isPending ? "Saving..." : "Save product"}</button></div>
    </form>
  </div>;
}
