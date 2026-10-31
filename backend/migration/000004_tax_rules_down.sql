ALTER TABLE invoices DROP CONSTRAINT IF EXISTS invoices_total_check;

ALTER TABLE invoices
    DROP COLUMN tax_cents,
    DROP COLUMN tax_rate_bps,
    DROP COLUMN tax_rule_id;

ALTER TABLE subscriptions DROP COLUMN tax_rule_id;

ALTER TABLE invoices
    ADD CONSTRAINT invoices_total_check
    CHECK (total_cents = subtotal_cents - discount_cents) NOT VALID;

DROP TABLE tax_rules;
