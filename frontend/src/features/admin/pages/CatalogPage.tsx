import { useState, type FormEvent } from "react";
import { Pencil, Plus, Trash2 } from "lucide-react";
import styles from "./CatalogPage.module.css";
import { RevenueRulesPage } from "./RevenueRulesPage";

type Field = { key: string; label: string; type?: "number" | "date" | "checkbox" | "textarea" | "select"; options?: string[]; required?: boolean };
type RecordValue = string | boolean;
type Row = { id: string; [key: string]: RecordValue };
type Config = { title: string; singular: string; description: string; fields: Field[]; rows: Row[] };

const configs: Record<string, Config> = {
  taxes: {
    title: "Taxes", singular: "Tax", description: "Rules applied to invoices.",
    fields: [
      { key: "name", label: "Rule name", required: true },
      { key: "rate", label: "Rate (%)", type: "number", required: true },
      { key: "appliesTo", label: "Applies to", type: "select", options: ["All invoices", "Products", "Subscriptions"] },
      { key: "region", label: "Region or jurisdiction", required: true },
      { key: "notes", label: "Rule notes", type: "textarea" },
    ],
    rows: [{ id: "tax-1", name: "Standard VAT", rate: "20", appliesTo: "All invoices", region: "Example region", notes: "Sample rule" }],
  },
  discounts: {
    title: "Discounts", singular: "Discount", description: "Eligibility, limits, and validity dates.",
    fields: [
      { key: "name", label: "Name", required: true },
      { key: "kind", label: "Discount type", type: "select", options: ["Percentage", "Fixed amount"] },
      { key: "value", label: "Value", type: "number", required: true },
      { key: "minPurchase", label: "Minimum purchase", type: "number" },
      { key: "minQuantity", label: "Minimum quantity", type: "number" },
      { key: "validFrom", label: "Valid from", type: "date" },
      { key: "validUntil", label: "Valid until", type: "date" },
      { key: "usageLimit", label: "Usage limit", type: "number" },
      { key: "appliesTo", label: "Applies to", type: "select", options: ["Products", "Subscriptions", "Both"] },
    ],
    rows: [{ id: "discount-1", name: "Welcome offer", kind: "Percentage", value: "10", minPurchase: "100", minQuantity: "1", validFrom: "2026-01-01", validUntil: "", usageLimit: "100", appliesTo: "Subscriptions" }],
  },
  "quotation-templates": {
    title: "Quotation templates", singular: "Quotation template", description: "Preset plan and product lines for new quotations.",
    fields: [
      { key: "name", label: "Template name", required: true },
      { key: "planLines", label: "Plan lines (one per line)", type: "textarea" },
      { key: "productLines", label: "Product lines (one per line)", type: "textarea" },
      { key: "validityDays", label: "Validity days", type: "number", required: true },
    ],
    rows: [{ id: "template-1", name: "Growth starter", planLines: "Growth × 1", productLines: "Cloud Suite × 1", validityDays: "30" }],
  },
};

const emptyRow = (fields: Field[]): Row => Object.fromEntries([
  ["id", ""], ...fields.map(field => [field.key, field.type === "checkbox" ? false : field.options?.[0] ?? ""]),
]) as Row;

export function CatalogPage({ section }: { section: keyof typeof configs }) {
  if (section === "taxes" || section === "discounts") return <RevenueRulesPage kind={section} />;
  const config = configs[section];
  const [rows, setRows] = useState(config.rows);
  const [form, setForm] = useState<Row | null>(null);
  const save = (event: FormEvent) => {
    event.preventDefault();
    if (!form) return;
    if (form.id) setRows(current => current.map(row => row.id === form.id ? form : row));
    else setRows(current => [{ ...form, id: crypto.randomUUID() }, ...current]);
    setForm(null);
  };

  return <div className={styles.page}>
    <header className={styles.header}><div><p className={styles.eyebrow}>Admin · Sample data</p><h1>{config.title}</h1><p>{config.description}</p></div>
      <button className={styles.primary} onClick={() => setForm(emptyRow(config.fields))}><Plus size={16} /> Add {config.singular}</button></header>
    <p className={styles.notice}>Demo only. Edits are held in this tab and are lost on refresh.</p>
    <section className={styles.list} aria-label={config.title}>
      {rows.map(row => <article className={styles.item} key={row.id}>
        <div><h2>{String(row.name)}</h2><dl>{config.fields.filter(field => field.key !== "name" && row[field.key] !== "").map(field =>
          <div key={field.key}><dt>{field.label}</dt><dd>{typeof row[field.key] === "boolean" ? row[field.key] ? "Yes" : "No" : String(row[field.key])}</dd></div>)}</dl></div>
        <div className={styles.actions}><button aria-label={`Edit ${row.name}`} onClick={() => setForm({ ...row })}><Pencil size={16} /></button>
          <button aria-label={`Delete ${row.name}`} onClick={() => setRows(current => current.filter(item => item.id !== row.id))}><Trash2 size={16} /></button></div>
      </article>)}
      {!rows.length && <p className={styles.empty}>No {config.title.toLowerCase()} yet.</p>}
    </section>
    {form && <div className={styles.backdrop}><div className={styles.dialog} role="dialog" aria-modal="true" aria-label={`Edit ${config.title}`}>
      <div className={styles.dialogHeading}><h2>{form.id ? "Edit" : "Add"} {config.singular}</h2><button onClick={() => setForm(null)} aria-label="Close">×</button></div>
      <form onSubmit={save} className={styles.form}>{config.fields.map(field =>
        <label className={field.type === "checkbox" ? styles.check : ""} key={field.key}>
          {field.type === "checkbox" ? <><input type="checkbox" checked={Boolean(form[field.key])} onChange={event => setForm({ ...form, [field.key]: event.target.checked })} />{field.label}</> :
            <>{field.label}{field.type === "textarea" ?
              <textarea value={String(form[field.key])} onChange={event => setForm({ ...form, [field.key]: event.target.value })} /> :
              field.type === "select" ? <select value={String(form[field.key])} onChange={event => setForm({ ...form, [field.key]: event.target.value })}>{field.options?.map(option => <option key={option}>{option}</option>)}</select> :
              <input type={field.type ?? "text"} min={field.type === "number" ? "0" : undefined} step={field.type === "number" ? "any" : undefined} required={field.required} value={String(form[field.key])} onChange={event => setForm({ ...form, [field.key]: event.target.value })} />}</>}
        </label>)}
        <div className={styles.formActions}><button type="button" onClick={() => setForm(null)}>Cancel</button><button className={styles.primary} type="submit">Save demo changes</button></div>
      </form></div></div>}
  </div>;
}
