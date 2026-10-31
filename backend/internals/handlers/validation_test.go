package handlers

import "testing"

func TestProductInputValidate(t *testing.T) {
	tests := []struct {
		name    string
		input   productInput
		wantErr bool
		want    string
	}{
		{name: "valid", input: productInput{Name: "  Starter  "}, want: "Starter"},
		{name: "empty name", input: productInput{Name: "   "}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.input.validate()
			if (err != nil) != tt.wantErr {
				t.Fatalf("validate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.want != "" && tt.input.Name != tt.want {
				t.Fatalf("name = %q, want %q", tt.input.Name, tt.want)
			}
		})
	}
}

func TestVariantInputValidate(t *testing.T) {
	tests := []struct {
		name         string
		input        variantInput
		wantErr      bool
		wantCurrency string
		wantInterval string
	}{
		{
			name:         "defaults currency and normalizes",
			input:        variantInput{Name: "Pro", PriceCents: 1999, BillingInterval: " MONTH "},
			wantCurrency: "USD",
			wantInterval: "month",
		},
		{
			name:         "lowercase currency is uppercased",
			input:        variantInput{Name: "Pro", PriceCents: 0, Currency: "eur", BillingInterval: "year"},
			wantCurrency: "EUR",
			wantInterval: "year",
		},
		{name: "missing name", input: variantInput{PriceCents: 100, BillingInterval: "month"}, wantErr: true},
		{name: "negative price", input: variantInput{Name: "Pro", PriceCents: -1, BillingInterval: "month"}, wantErr: true},
		{name: "bad interval", input: variantInput{Name: "Pro", BillingInterval: "week"}, wantErr: true},
		{name: "bad currency", input: variantInput{Name: "Pro", Currency: "US", BillingInterval: "month"}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.input.validate()
			if (err != nil) != tt.wantErr {
				t.Fatalf("validate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil {
				return
			}
			if tt.input.Currency != tt.wantCurrency {
				t.Fatalf("currency = %q, want %q", tt.input.Currency, tt.wantCurrency)
			}
			if tt.input.BillingInterval != tt.wantInterval {
				t.Fatalf("billingInterval = %q, want %q", tt.input.BillingInterval, tt.wantInterval)
			}
		})
	}
}

func TestInputActiveDefaults(t *testing.T) {
	product := productInput{}
	if !product.active() {
		t.Fatal("product active should default to true")
	}

	variant := variantInput{}
	if !variant.active() {
		t.Fatal("variant active should default to true")
	}

	off := false
	if (productInput{Active: &off}).active() {
		t.Fatal("explicit false should be respected")
	}
}
