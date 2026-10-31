package repository

import (
	"context"

	"recur_backend/internals/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductRepository struct {
	db *pgxpool.Pool
}

func NewProductRepository(db *pgxpool.Pool) *ProductRepository {
	return &ProductRepository{db: db}
}

const productColumns = `id, name, description, is_active, created_at, updated_at`

const variantColumns = `id, product_id, name, description, price_cents, currency, billing_interval, is_active, created_at, updated_at`

func scanProduct(row pgx.Row) (models.Product, error) {
	var p models.Product
	err := row.Scan(&p.ID, &p.Name, &p.Description, &p.Active, &p.CreatedAt, &p.UpdatedAt)
	return p, err
}

func scanVariant(row pgx.Row) (models.ProductVariant, error) {
	var v models.ProductVariant
	err := row.Scan(&v.ID, &v.ProductID, &v.Name, &v.Description, &v.PriceCents,
		&v.Currency, &v.BillingInterval, &v.Active, &v.CreatedAt, &v.UpdatedAt)
	return v, err
}

func (r *ProductRepository) CreateProduct(ctx context.Context, p models.Product) (models.Product, error) {
	row := r.db.QueryRow(ctx, `
		INSERT INTO products (name, description, is_active)
		VALUES ($1, $2, $3)
		RETURNING `+productColumns,
		p.Name, p.Description, p.Active)
	return scanProduct(row)
}

func (r *ProductRepository) ListProducts(ctx context.Context) ([]models.Product, error) {
	rows, err := r.db.Query(ctx, `SELECT `+productColumns+` FROM products ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := []models.Product{}
	for rows.Next() {
		p, err := scanProduct(rows)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, rows.Err()
}

func (r *ProductRepository) GetProduct(ctx context.Context, id int64) (models.Product, error) {
	row := r.db.QueryRow(ctx, `SELECT `+productColumns+` FROM products WHERE id = $1`, id)
	return scanProduct(row)
}

func (r *ProductRepository) UpdateProduct(ctx context.Context, p models.Product) (models.Product, error) {
	row := r.db.QueryRow(ctx, `
		UPDATE products
		SET name = $2, description = $3, is_active = $4, updated_at = NOW()
		WHERE id = $1
		RETURNING `+productColumns,
		p.ID, p.Name, p.Description, p.Active)
	return scanProduct(row)
}

func (r *ProductRepository) DeleteProduct(ctx context.Context, id int64) (bool, error) {
	tag, err := r.db.Exec(ctx, `DELETE FROM products WHERE id = $1`, id)
	if err != nil {
		return false, err
	}
	return tag.RowsAffected() > 0, nil
}

func (r *ProductRepository) CreateVariant(ctx context.Context, v models.ProductVariant) (models.ProductVariant, error) {
	row := r.db.QueryRow(ctx, `
		INSERT INTO product_variants (product_id, name, description, price_cents, currency, billing_interval, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING `+variantColumns,
		v.ProductID, v.Name, v.Description, v.PriceCents, v.Currency, v.BillingInterval, v.Active)
	return scanVariant(row)
}

func (r *ProductRepository) ListVariants(ctx context.Context, productID int64) ([]models.ProductVariant, error) {
	rows, err := r.db.Query(ctx,
		`SELECT `+variantColumns+` FROM product_variants WHERE product_id = $1 ORDER BY id`, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	variants := []models.ProductVariant{}
	for rows.Next() {
		v, err := scanVariant(rows)
		if err != nil {
			return nil, err
		}
		variants = append(variants, v)
	}
	return variants, rows.Err()
}

func (r *ProductRepository) GetVariant(ctx context.Context, id int64) (models.ProductVariant, error) {
	row := r.db.QueryRow(ctx, `SELECT `+variantColumns+` FROM product_variants WHERE id = $1`, id)
	return scanVariant(row)
}

func (r *ProductRepository) UpdateVariant(ctx context.Context, v models.ProductVariant) (models.ProductVariant, error) {
	row := r.db.QueryRow(ctx, `
		UPDATE product_variants
		SET name = $2, description = $3, price_cents = $4, currency = $5,
			billing_interval = $6, is_active = $7, updated_at = NOW()
		WHERE id = $1
		RETURNING `+variantColumns,
		v.ID, v.Name, v.Description, v.PriceCents, v.Currency, v.BillingInterval, v.Active)
	return scanVariant(row)
}

func (r *ProductRepository) DeleteVariant(ctx context.Context, id int64) (bool, error) {
	tag, err := r.db.Exec(ctx, `DELETE FROM product_variants WHERE id = $1`, id)
	if err != nil {
		return false, err
	}
	return tag.RowsAffected() > 0, nil
}
