import type { Subscription } from "./types";

export const subscriptions: Subscription[] = [
  {
    id: "sub-1048", number: "SUB-2025-1048", customer: "Northstar Health Group", initials: "NH",
    plan: "Enterprise Annual", startDate: "Oct 01, 2025", expiryDate: "Sep 30, 2026", paymentTerms: "Net 30",
    status: "Active", sourceTemplate: "Enterprise rollout · 30 days",
    lines: [
      { product: "Care Platform — Enterprise", quantity: 120, unitPrice: 89, taxRate: 8.25 },
      { product: "Priority Support", quantity: 1, unitPrice: 1200, taxRate: 8.25 },
      { product: "Onboarding & data migration", quantity: 1, unitPrice: 4500, taxRate: 8.25 },
    ],
  },
  {
    id: "sub-1047", number: "SUB-2025-1047", customer: "Lumen Field Services", initials: "LF",
    plan: "Business Plus", startDate: "Nov 01, 2025", expiryDate: "Oct 31, 2026", paymentTerms: "Net 15",
    status: "Quotation", sourceTemplate: "Business standard · 21 days",
    lines: [
      { product: "Operations Suite", quantity: 48, unitPrice: 64, taxRate: 7.5 },
      { product: "Advanced analytics", quantity: 8, unitPrice: 18, taxRate: 7.5 },
    ],
  },
  {
    id: "sub-1046", number: "SUB-2025-1046", customer: "Cedar & Finch Retail", initials: "CF",
    plan: "Growth Monthly", startDate: "Oct 15, 2025", expiryDate: "Nov 14, 2025", paymentTerms: "Due on receipt",
    status: "Confirmed",
    lines: [{ product: "Retail Workspace", quantity: 32, unitPrice: 42, taxRate: 6.25 }, { product: "Inventory sync", quantity: 1, unitPrice: 240, taxRate: 6.25 }],
  },
  {
    id: "sub-1045", number: "SUB-2025-1045", customer: "Meridian Learning Co.", initials: "ML",
    plan: "Team Annual", startDate: "Dec 01, 2025", expiryDate: "Nov 30, 2026", paymentTerms: "Net 30",
    status: "Draft", sourceTemplate: "Education team · 14 days",
    lines: [{ product: "Learning Hub", quantity: 75, unitPrice: 26, taxRate: 5 }, { product: "Educator seats", quantity: 10, unitPrice: 9, taxRate: 5 }],
  },
  {
    id: "sub-1044", number: "SUB-2025-1044", customer: "Atlas Works Studio", initials: "AW",
    plan: "Business Plus", startDate: "Apr 01, 2025", expiryDate: "Sep 30, 2025", paymentTerms: "Net 30",
    status: "Closed",
    lines: [{ product: "Operations Suite", quantity: 16, unitPrice: 64, taxRate: 7.5 }, { product: "Priority Support", quantity: 1, unitPrice: 600, taxRate: 7.5 }],
  },
];
