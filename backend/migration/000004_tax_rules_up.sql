CREATE TABLE tax_rules (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    rate_bps INTEGER NOT NULL CHECK (rate_bps BETWEEN 0 AND 10000), -- basis points, 1800 = 18%
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_by TEXT NOT NULL, -- Clerk user ID of the admin who created it
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

ALTER TABLE subscriptions
    ADD COLUMN tax_rule_id BIGINT REFERENCES tax_rules (id) ON DELETE SET NULL;

ALTER TABLE invoices
    ADD COLUMN tax_rule_id BIGINT REFERENCES tax_rules (id) ON DELETE SET NULL,
    ADD COLUMN tax_rate_bps INTEGER NOT NULL DEFAULT 0 CHECK (tax_rate_bps BETWEEN 0 AND 10000),
    ADD COLUMN tax_cents BIGINT NOT NULL DEFAULT 0 CHECK (tax_cents >= 0);

-- Invoice totals now include tax, so replace the old total check.
DO $$
DECLARE constraint_name TEXT;
BEGIN
    SELECT conname INTO constraint_name
    FROM pg_constraint
    WHERE conrelid = 'invoices'::regclass
      AND contype = 'c'
      AND pg_get_constraintdef(oid) LIKE '%subtotal_cents%'
      AND pg_get_constraintdef(oid) LIKE '%discount_cents%';

    IF constraint_name IS NOT NULL THEN
        EXECUTE format('ALTER TABLE invoices DROP CONSTRAINT %I', constraint_name);
    END IF;
END $$;

ALTER TABLE invoices
    ADD CONSTRAINT invoices_total_check
    CHECK (total_cents = subtotal_cents - discount_cents + tax_cents);
