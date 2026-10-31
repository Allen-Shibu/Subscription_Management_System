package repository

import (
	"context"

	"recur_backend/internals/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TaxRuleRepository struct {
	db *pgxpool.Pool
}

func NewTaxRuleRepository(db *pgxpool.Pool) *TaxRuleRepository {
	return &TaxRuleRepository{db: db}
}

const taxRuleColumns = `id, name, rate_bps, is_active, created_by, created_at, updated_at`

func scanTaxRule(row pgx.Row) (models.TaxRule, error) {
	var t models.TaxRule
	err := row.Scan(&t.ID, &t.Name, &t.RateBps, &t.Active, &t.CreatedBy, &t.CreatedAt, &t.UpdatedAt)
	return t, err
}

func (r *TaxRuleRepository) CreateTaxRule(ctx context.Context, t models.TaxRule) (models.TaxRule, error) {
	row := r.db.QueryRow(ctx, `
		INSERT INTO tax_rules (name, rate_bps, is_active, created_by)
		VALUES ($1, $2, $3, $4)
		RETURNING `+taxRuleColumns,
		t.Name, t.RateBps, t.Active, t.CreatedBy)
	return scanTaxRule(row)
}

func (r *TaxRuleRepository) ListTaxRules(ctx context.Context) ([]models.TaxRule, error) {
	rows, err := r.db.Query(ctx, `SELECT `+taxRuleColumns+` FROM tax_rules ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	rules := []models.TaxRule{}
	for rows.Next() {
		t, err := scanTaxRule(rows)
		if err != nil {
			return nil, err
		}
		rules = append(rules, t)
	}
	return rules, rows.Err()
}

func (r *TaxRuleRepository) GetTaxRule(ctx context.Context, id int64) (models.TaxRule, error) {
	row := r.db.QueryRow(ctx, `SELECT `+taxRuleColumns+` FROM tax_rules WHERE id = $1`, id)
	return scanTaxRule(row)
}

func (r *TaxRuleRepository) UpdateTaxRule(ctx context.Context, t models.TaxRule) (models.TaxRule, error) {
	row := r.db.QueryRow(ctx, `
		UPDATE tax_rules
		SET name = $2, rate_bps = $3, is_active = $4, updated_at = NOW()
		WHERE id = $1
		RETURNING `+taxRuleColumns,
		t.ID, t.Name, t.RateBps, t.Active)
	return scanTaxRule(row)
}

func (r *TaxRuleRepository) DeleteTaxRule(ctx context.Context, id int64) (bool, error) {
	tag, err := r.db.Exec(ctx, `DELETE FROM tax_rules WHERE id = $1`, id)
	if err != nil {
		return false, err
	}
	return tag.RowsAffected() > 0, nil
}
