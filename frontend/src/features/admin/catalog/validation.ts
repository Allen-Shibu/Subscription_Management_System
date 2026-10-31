import { z } from "zod";

const money = z.number().finite().multipleOf(0.01).nonnegative();
export const productSchema = z.object({
  name: z.string().trim().min(1, "Name is required").max(100, "Use 100 characters or fewer"),
  type: z.enum(["Service", "Physical", "Digital"]),
  salesPrice: money.positive("Sales price must be greater than zero"),
  costPrice: money,
  notes: z.string(),
  variants: z.array(z.object({
    id: z.string().optional(),
    attribute: z.string().trim().min(1, "Attribute is required"),
    value: z.string().trim().min(1, "Value is required"),
    extraPrice: money,
  })),
}).superRefine((product, context) => {
  const seen = new Set<string>();
  product.variants.forEach((variant, index) => {
    const key = `${variant.attribute.trim().toLocaleLowerCase()}\0${variant.value.trim().toLocaleLowerCase()}`;
    if (seen.has(key)) context.addIssue({
      code: "custom", path: ["variants", index, "value"], message: "This attribute and value pair already exists",
    });
    seen.add(key);
  });
});

export const planSchema = z.object({
  name: z.string().trim().min(1, "Name is required").max(100, "Use 100 characters or fewer"),
  price: money.positive("Price must be greater than zero"),
  billingPeriod: z.enum(["Daily", "Weekly", "Monthly", "Yearly"]),
  minQuantity: z.number().int("Use a whole number").min(1, "Minimum quantity must be at least 1"),
  availableFrom: z.string(),
  availableTo: z.string(),
  autoClose: z.boolean(),
  earlyClose: z.boolean(),
  allowPause: z.boolean(),
  allowRenew: z.boolean(),
}).superRefine((plan, context) => {
  if (plan.availableFrom && plan.availableTo && plan.availableTo <= plan.availableFrom) {
    context.addIssue({ code: "custom", path: ["availableTo"], message: "End date must be after start date" });
  }
});

export function planAvailability(plan: { availableFrom: string; availableTo: string }, today: string) {
  if (plan.availableFrom && today < plan.availableFrom) return "Scheduled";
  if (plan.availableTo && today > plan.availableTo) return "Ended";
  return "Available";
}
