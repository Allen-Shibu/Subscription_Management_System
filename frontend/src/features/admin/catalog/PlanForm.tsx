import { useEffect, useState } from "react";
import { useAuth } from "@clerk/react";
import { zodResolver } from "@hookform/resolvers/zod";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useForm, type FieldPath } from "react-hook-form";
import { Link, useNavigate, useParams } from "react-router-dom";
import { ArrowLeft } from "lucide-react";
import { z } from "zod";
import { catalog } from "./api";
import { CatalogError, type PlanInput } from "./types";
import { planSchema } from "./validation";
import styles from "./Catalog.module.css";

type Form = z.infer<typeof planSchema>;
const empty: Form = {
  name: "", price: 0, billingPeriod: "Monthly", minQuantity: 1, availableFrom: "", availableTo: "",
  autoClose: false, earlyClose: false, allowPause: false, allowRenew: false,
};

export function PlanForm() {
  const { id } = useParams();
  const { getToken } = useAuth();
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const [saveError, setSaveError] = useState("");
  const plan = useQuery({ queryKey: ["plan", id], queryFn: () => catalog.getPlan(id!, getToken), enabled: Boolean(id), retry: false });
  const { register, handleSubmit, reset, setError, formState: { errors } } = useForm<Form>({
    resolver: zodResolver(planSchema), defaultValues: empty,
  });
  useEffect(() => { if (plan.data) reset(plan.data); }, [plan.data, reset]);
  const save = useMutation({
    mutationFn: (input: PlanInput) => catalog.savePlan(id, input, getToken),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ["plans"] });
      navigate("/admin/plans", { state: { toast: id ? "Plan updated" : "Plan created" } });
    },
    onError: (error: Error) => {
      setSaveError(error.message);
      if (error instanceof CatalogError && error.status === 400) {
        Object.entries(error.fields).forEach(([field, message]) => setError(field as FieldPath<Form>, { type: "server", message }));
      }
    },
  });

  if (id && plan.isPending) return <div className={styles.skeleton} role="status" aria-label="Loading plan"><span /><span /><span /></div>;
  if (id && plan.isError) return <div className={styles.empty} role="alert"><h2>Could not load plan</h2><p>{plan.error.message}</p><button onClick={() => plan.refetch()}>Retry</button></div>;
  return <div className={styles.page}>
    <header className={styles.header}><div><Link className={styles.back} to="/admin/plans"><ArrowLeft size={15} /> Plans</Link>
      <h1>{id ? "Edit plan" : "Add plan"}</h1><p>Configure pricing, availability, and lifecycle controls.</p></div></header>
    {saveError && <div className={styles.error} role="alert">{saveError}</div>}
    <form className={styles.form} onSubmit={handleSubmit(data => { setSaveError(""); save.mutate(data); })} noValidate>
      <section className={styles.formCard}><h2>Plan details</h2><div className={styles.formGrid}>
        <label>Name<input {...register("name")} maxLength={100} aria-invalid={!!errors.name} /><small>{errors.name?.message}</small></label>
        <label>Price<input type="number" min="0" step="0.01" {...register("price", { valueAsNumber: true })} aria-invalid={!!errors.price} /><small>{errors.price?.message}</small></label>
        <label>Billing period<select {...register("billingPeriod")}>{["Daily", "Weekly", "Monthly", "Yearly"].map(period => <option key={period}>{period}</option>)}</select><small>{errors.billingPeriod?.message}</small></label>
        <label>Minimum quantity<input type="number" min="1" step="1" {...register("minQuantity", { valueAsNumber: true })} aria-invalid={!!errors.minQuantity} /><small>{errors.minQuantity?.message}</small></label>
        <label>Available from<input type="date" {...register("availableFrom")} aria-invalid={!!errors.availableFrom} /><small>{errors.availableFrom?.message}</small></label>
        <label>Available to<input type="date" {...register("availableTo")} aria-invalid={!!errors.availableTo} /><small>{errors.availableTo?.message}</small></label>
      </div></section>
      <section className={styles.formCard}><h2>Lifecycle options</h2><div className={styles.checkList}>
        {[
          { key: "autoClose" as const, label: "Auto-close when the plan ends", help: "Close subscriptions when this plan reaches its end date." },
          { key: "earlyClose" as const, label: "Allow early close", help: "Let staff close a subscription before its scheduled end." },
          { key: "allowPause" as const, label: "Allow pause", help: "Let staff temporarily pause subscriptions on this plan." },
          { key: "allowRenew" as const, label: "Allow renew", help: "Let subscriptions renew at the end of a billing period." },
        ].map(option => <label key={option.key}><input type="checkbox" {...register(option.key)} /><span><strong>{option.label}</strong><small>{option.help}</small></span></label>)}
      </div></section>
      <div className={styles.formActions}><Link to="/admin/plans">Cancel</Link><button className={styles.primary} disabled={save.isPending}>{save.isPending ? "Saving..." : "Save plan"}</button></div>
    </form>
  </div>;
}
