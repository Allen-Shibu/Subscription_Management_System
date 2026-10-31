**Core idea:** Build a subscription billing platform for a fictional company, Brightloop Software. It replaces spreadsheets and manual work with one system covering plans, subscriptions, invoices, payments, and reports.


1. **Auth and roles**
   - Signup with a unique email and a strong password (more than 8 characters, with upper, lower, and special characters).
   - Password reset via emailed link.
   - Three roles: **Admin** (everything, and only they can add staff), **Internal staff** (limited access), **Customer** (portal showing only their own subscriptions and invoices).

2. **Products and plans** (admin)
   - Product: name, type, sales price, cost price, notes.
   - Variants, e.g. brand = Odoo adds 560 to the price.
   - Plan: price, billing period (daily, weekly, monthly, yearly), min quantity, and availability dates.
   - Per-plan flags for auto-close on end, early close, pause, and renew.

3. **Quotation templates and subscriptions**
   - Reusable templates with preset plan and product lines, valid for N days.
   - A subscription has a number, customer, plan, start and expiry dates, payment terms, and order lines (product, qty, unit price, tax, amount).
   - Status flow: **Draft → Quotation → Confirmed → Active → Closed**, visible to the team.

4. **Invoicing**
   - Auto-generated from subscriptions, with customer details, lines, taxes, and totals prefilled.
   - Taxes come from admin-defined rules.
   - Actions: confirm, cancel, send, print.
   - Status flow: **Draft → Confirmed → Paid**.
   - **Idempotent billing:** no duplicate invoices and no skipped ones.

5. **Discounts** (admin only)
   - Fixed or percentage, with a minimum purchase or quantity, start and end dates, and a usage limit.
   - Applies to specific products or whole subscriptions.
   - Applied automatically on the next invoice, and stops on expiry or when the usage limit is hit.

6. **Payments**
   - Record method, amount, date, and note against an invoice.
   - Clear view of paid vs outstanding, so overdue invoices can be chased.

7. **Finance dashboard**
   - Active subscriptions, revenue, payments received, amount owed, and overdue invoices.
   - Filters and a summary view.
   - **Performance:** pages must load in under 2 seconds with thousands of subscriptions.

**Hidden requirements (what judges will likely check)**
- Role-based access control enforced on the backend, not just hidden in the UI.
- A scheduled job or cron for recurring invoice generation, with duplicate protection (e.g. a unique constraint on subscription + billing period).
- Strict state machines for subscription and invoice status.
- Discount usage counters and expiry logic.
- Pagination, indexes, and aggregate queries for the dashboard to hit the 2s target.

**Rough data model:** User, Product, Variant, Plan, QuotationTemplate, Subscription, SubscriptionLine, Tax, Discount, Invoice, InvoiceLine, Payment.
