// Package pdf renders invoice documents.
package pdf

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"time"

	"recur_backend/internals/models"

	"github.com/jung-kurt/gofpdf"
)

const (
	companyName = "Brightloop Software"
	pageMargin  = 15.0
	contentW    = 180.0
)

// RenderInvoice draws a single-page A4 invoice and returns the PDF bytes.
func RenderInvoice(d models.InvoiceDetails) ([]byte, error) {
	doc := gofpdf.New("P", "mm", "A4", "")
	doc.SetMargins(pageMargin, pageMargin, pageMargin)
	doc.SetAutoPageBreak(true, 20)
	doc.SetTitle("Invoice "+d.InvoiceNumber, false)
	doc.SetAuthor(companyName, false)
	doc.SetCreator("Recur", false)
	doc.AddPage()

	renderHeader(doc, d)
	renderParties(doc, d)
	renderLineItems(doc, d)
	renderTotals(doc, d)
	renderFooter(doc, d)

	if doc.Error() != nil {
		return nil, fmt.Errorf("render invoice pdf: %w", doc.Error())
	}

	var buf bytes.Buffer
	if err := doc.Output(&buf); err != nil {
		return nil, fmt.Errorf("render invoice pdf: %w", err)
	}
	return buf.Bytes(), nil
}

func renderHeader(doc *gofpdf.Fpdf, d models.InvoiceDetails) {
	doc.SetFont("Helvetica", "B", 20)
	doc.SetTextColor(24, 24, 27)
	doc.CellFormat(120, 10, companyName, "", 0, "L", false, 0, "")

	doc.SetFont("Helvetica", "B", 24)
	doc.CellFormat(60, 10, "INVOICE", "", 1, "R", false, 0, "")

	doc.SetFont("Helvetica", "", 10)
	doc.SetTextColor(113, 113, 122)
	doc.CellFormat(120, 6, "Subscription billing", "", 0, "L", false, 0, "")
	doc.SetFont("Helvetica", "B", 11)
	doc.SetTextColor(24, 24, 27)
	doc.CellFormat(60, 6, d.InvoiceNumber, "", 1, "R", false, 0, "")

	doc.Ln(2)
	doc.SetDrawColor(228, 228, 231)
	doc.SetLineWidth(0.4)
	doc.Line(pageMargin, doc.GetY(), pageMargin+contentW, doc.GetY())
	doc.Ln(6)
}

func renderParties(doc *gofpdf.Fpdf, d models.InvoiceDetails) {
	left := pageMargin
	top := doc.GetY()

	doc.SetFont("Helvetica", "B", 9)
	doc.SetTextColor(113, 113, 122)
	doc.SetXY(left, top)
	doc.CellFormat(80, 5, "BILLED TO", "", 1, "L", false, 0, "")

	doc.SetFont("Helvetica", "B", 11)
	doc.SetTextColor(24, 24, 27)
	doc.MultiCell(85, 5, d.CustomerName, "", "L", false)
	doc.SetFont("Helvetica", "", 10)
	doc.SetTextColor(63, 63, 70)
	doc.MultiCell(85, 5, d.CustomerEmail, "", "L", false)
	leftBottom := doc.GetY()

	doc.SetXY(pageMargin+95, top)
	doc.SetFont("Helvetica", "B", 9)
	doc.SetTextColor(113, 113, 122)
	doc.CellFormat(85, 5, "INVOICE DETAILS", "", 1, "L", false, 0, "")

	meta := [][2]string{
		{"Status", strings.ToUpper(d.Status)},
		{"Issued", formatDatePtr(d.IssuedAt)},
		{"Due", formatDatePtr(d.DueAt)},
		{"Period", formatDate(d.PeriodStart) + " to " + formatDate(d.PeriodEnd)},
	}
	if d.PaidAt != nil {
		meta = append(meta, [2]string{"Paid", formatDatePtr(d.PaidAt)})
	}
	for _, row := range meta {
		doc.SetX(pageMargin + 95)
		doc.SetFont("Helvetica", "", 9)
		doc.SetTextColor(113, 113, 122)
		doc.CellFormat(22, 5, row[0], "", 0, "L", false, 0, "")
		doc.SetFont("Helvetica", "B", 9)
		doc.SetTextColor(24, 24, 27)
		doc.CellFormat(63, 5, row[1], "", 1, "L", false, 0, "")
	}
	rightBottom := doc.GetY()

	nextY := leftBottom
	if rightBottom > nextY {
		nextY = rightBottom
	}
	doc.SetY(nextY + 6)
}

func renderLineItems(doc *gofpdf.Fpdf, d models.InvoiceDetails) {
	widths := []float64{90, 20, 35, 35}
	headers := []string{"Description", "Qty", "Unit price", "Amount"}

	doc.SetFillColor(244, 244, 245)
	doc.SetFont("Helvetica", "B", 9)
	doc.SetTextColor(63, 63, 70)
	for i, header := range headers {
		align := "L"
		if i > 0 {
			align = "R"
		}
		doc.CellFormat(widths[i], 8, header, "", 0, align, true, 0, "")
	}
	doc.Ln(8)

	description := d.ProductName
	if d.VariantName != "" && d.VariantName != d.ProductName {
		description += " - " + d.VariantName
	}
	switch strings.ToLower(d.BillingInterval) {
	case "month":
		description += " (monthly)"
	case "year":
		description += " (yearly)"
	default:
		if d.BillingInterval != "" {
			description += " (" + d.BillingInterval + ")"
		}
	}

	lineAmount := d.UnitPriceCents * int64(d.Quantity)
	values := []string{
		description,
		strconv.Itoa(d.Quantity),
		formatMoney(d.UnitPriceCents, d.Currency),
		formatMoney(lineAmount, d.Currency),
	}

	doc.SetFont("Helvetica", "", 9)
	doc.SetTextColor(24, 24, 27)
	doc.SetDrawColor(228, 228, 231)
	for i, value := range values {
		align := "L"
		if i > 0 {
			align = "R"
		}
		doc.CellFormat(widths[i], 9, value, "B", 0, align, false, 0, "")
	}
	doc.Ln(12)
}

func renderTotals(doc *gofpdf.Fpdf, d models.InvoiceDetails) {
	row := func(label, value string, bold bool) {
		doc.SetX(pageMargin + 90)
		style := ""
		if bold {
			style = "B"
		}
		doc.SetFont("Helvetica", style, 10)
		doc.SetTextColor(63, 63, 70)
		doc.CellFormat(45, 7, label, "", 0, "R", false, 0, "")
		doc.SetTextColor(24, 24, 27)
		doc.CellFormat(45, 7, value, "", 1, "R", false, 0, "")
	}

	row("Subtotal", formatMoney(d.SubtotalCents, d.Currency), false)
	if d.DiscountCents > 0 {
		row("Discount", "-"+formatMoney(d.DiscountCents, d.Currency), false)
	}
	if d.TaxCents > 0 {
		taxLabel := "Tax"
		if d.TaxRateBps > 0 {
			if d.TaxRateBps%100 == 0 {
				taxLabel = fmt.Sprintf("Tax (%d%%)", d.TaxRateBps/100)
			} else {
				taxLabel = fmt.Sprintf("Tax (%.2f%%)", float64(d.TaxRateBps)/100.0)
			}
		}
		row(taxLabel, formatMoney(d.TaxCents, d.Currency), false)
	}
	doc.SetDrawColor(228, 228, 231)
	doc.Line(pageMargin+90, doc.GetY(), pageMargin+180, doc.GetY())
	doc.Ln(1)
	totalLabel := "Total due"
	if d.Status == models.InvoiceStatusPaid {
		totalLabel = "Total paid"
	}
	row(totalLabel, formatMoney(d.TotalCents, d.Currency), true)
}

func renderFooter(doc *gofpdf.Fpdf, d models.InvoiceDetails) {
	if doc.GetY() > 260 {
		doc.AddPage()
	}
	doc.SetY(265)
	doc.SetDrawColor(228, 228, 231)
	doc.Line(pageMargin, doc.GetY(), pageMargin+contentW, doc.GetY())
	doc.Ln(3)

	doc.SetFont("Helvetica", "", 8)
	doc.SetTextColor(113, 113, 122)
	doc.CellFormat(contentW, 4, "Thank you for your business. Please contact billing@brightloop.example with any questions.", "", 1, "L", false, 0, "")
	doc.CellFormat(contentW, 4, "Invoice "+d.InvoiceNumber+" - generated "+time.Now().UTC().Format("02 Jan 2006 15:04 UTC"), "", 1, "L", false, 0, "")
}

func formatMoney(cents int64, currency string) string {
	sign := ""
	if cents < 0 {
		sign = "-"
		cents = -cents
	}
	whole := strconv.FormatInt(cents/100, 10)
	for i := len(whole) - 3; i > 0; i -= 3 {
		whole = whole[:i] + "," + whole[i:]
	}
	return fmt.Sprintf("%s%s %s.%02d", sign, currency, whole, cents%100)
}

func formatDate(t time.Time) string {
	if t.IsZero() {
		return "-"
	}
	return t.Format("02 Jan 2006")
}

func formatDatePtr(t *time.Time) string {
	if t == nil {
		return "-"
	}
	return formatDate(*t)
}
