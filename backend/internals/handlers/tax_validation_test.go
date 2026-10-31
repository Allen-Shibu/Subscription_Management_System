package handlers

import "testing"

func TestTaxRuleInputValidate(t *testing.T) {
	tests := []struct {
		name    string
		input   taxRuleInput
		wantErr bool
	}{
		{name: "valid", input: taxRuleInput{Name: " VAT 18% ", RateBps: 1800}},
		{name: "zero rate allowed", input: taxRuleInput{Name: "Exempt", RateBps: 0}},
		{name: "missing name", input: taxRuleInput{Name: "  ", RateBps: 100}, wantErr: true},
		{name: "rate too low", input: taxRuleInput{Name: "Tax", RateBps: -1}, wantErr: true},
		{name: "rate too high", input: taxRuleInput{Name: "Tax", RateBps: 10001}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.input.validate()
			if (err != nil) != tt.wantErr {
				t.Fatalf("validate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil && tt.input.Name != "VAT 18%" && tt.input.Name != "Exempt" {
				t.Fatalf("name was not trimmed: %q", tt.input.Name)
			}
		})
	}
}

func TestStaffInviteInputValidate(t *testing.T) {
	tests := []struct {
		name    string
		email   string
		wantErr bool
	}{
		{name: "valid", email: "staff@example.com"},
		{name: "valid trimmed", email: " staff@example.com "},
		{name: "missing", email: "", wantErr: true},
		{name: "not an email", email: "nope", wantErr: true},
		{name: "display name rejected", email: "Staff <staff@example.com>", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := staffInviteInput{Email: tt.email}
			err := in.validate()
			if (err != nil) != tt.wantErr {
				t.Fatalf("validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
