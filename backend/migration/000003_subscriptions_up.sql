-- Recurring subscriptions. Price/interval are snapshotted from the variant so
-- later variant edits don't change what an existing customer is billed.
CREATE TABLE subscriptions (
    id BIGSERIAL PRIMARY KEY,
    customer_id BIGINT NOT NULL REFERENCES users (id),
    variant_id BIGINT NOT NULL REFERENCES product_variants (id),
    status TEXT NOT NULL DEFAULT 'active'
        CHECK (status IN ('active', 'trialing', 'past_due', 'canceled', 'expired')),
    quantity INTEGER NOT NULL DEFAULT 1 CHECK (quantity > 0),
    unit_price_cents BIGINT NOT NULL CHECK (unit_price_cents >= 0),
    currency TEXT NOT NULL DEFAULT 'USD' CHECK (currency ~ '^[A-Z]{3}$'),
    billing_interval TEXT NOT NULL CHECK (billing_interval IN ('month', 'year')),
    current_period_start TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    current_period_end TIMESTAMPTZ NOT NULL,
    cancel_at_period_end BOOLEAN NOT NULL DEFAULT FALSE,
    canceled_at TIMESTAMPTZ,
    ended_at TIMESTAMPTZ,
    created_by TEXT NOT NULL, -- Clerk user ID of the staff/admin who entered it
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CHECK (current_period_end > current_period_start)
);

CREATE INDEX idx_subscriptions_customer_id ON subscriptions (customer_id);
CREATE INDEX idx_subscriptions_variant_id ON subscriptions (variant_id);
CREATE INDEX idx_subscriptions_status ON subscriptions (status);

CREATE SEQUENCE invoice_number_seq;

CREATE TABLE invoices (
    id BIGSERIAL PRIMARY KEY,
    invoice_number TEXT NOT NULL UNIQUE
        DEFAULT 'INV-' || lpad(nextval('invoice_number_seq')::text, 6, '0'),
    subscription_id BIGINT NOT NULL REFERENCES subscriptions (id),
    status TEXT NOT NULL DEFAULT 'draft'
        CHECK (status IN ('draft', 'open', 'paid', 'void', 'uncollectible')),
    period_start TIMESTAMPTZ NOT NULL,
    period_end TIMESTAMPTZ NOT NULL,
    subtotal_cents BIGINT NOT NULL CHECK (subtotal_cents >= 0),
    discount_cents BIGINT NOT NULL DEFAULT 0 CHECK (discount_cents >= 0),
    total_cents BIGINT NOT NULL CHECK (total_cents >= 0),
    issued_at TIMESTAMPTZ,
    due_at TIMESTAMPTZ,
    paid_at TIMESTAMPTZ,
    email_sent_at TIMESTAMPTZ,
    created_by TEXT NOT NULL, -- Clerk user ID of the staff/admin who issued it
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CHECK (period_end > period_start),
    CHECK (total_cents = subtotal_cents - discount_cents),
    UNIQUE (subscription_id, period_start) -- one invoice per billing period
);

CREATE INDEX idx_invoices_subscription_id ON invoices (subscription_id);
CREATE INDEX idx_invoices_status ON invoices (status);

-- One-time discounts; applied_invoice_id stays NULL until consumed by an invoice.
CREATE TABLE subscription_discounts (
    id BIGSERIAL PRIMARY KEY,
    subscription_id BIGINT NOT NULL REFERENCES subscriptions (id) ON DELETE CASCADE,
    amount_cents BIGINT NOT NULL CHECK (amount_cents > 0),
    reason TEXT NOT NULL,
    applied_invoice_id BIGINT REFERENCES invoices (id) ON DELETE SET NULL,
    created_by TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_subscription_discounts_subscription_id ON subscription_discounts (subscription_id);

-- Refunds are recorded against a paid invoice; no payment provider call is made here.
CREATE TABLE refunds (
    id BIGSERIAL PRIMARY KEY,
    invoice_id BIGINT NOT NULL REFERENCES invoices (id),
    amount_cents BIGINT NOT NULL CHECK (amount_cents > 0),
    reason TEXT NOT NULL,
    refunded_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_refunds_invoice_id ON refunds (invoice_id);
