package domain

import (
	"testing"
	"time"

	"github.com/shopspring/decimal"
)

func TestNewSummary(t *testing.T) {
	transactions := []Transaction{
		{ID: 1, Date: time.Date(2024, 7, 15, 0, 0, 0, 0, time.UTC), Amount: decimal.RequireFromString("60.5"), Type: Credit},
		{ID: 2, Date: time.Date(2024, 7, 28, 0, 0, 0, 0, time.UTC), Amount: decimal.RequireFromString("10.3"), Type: Debit},
		{ID: 3, Date: time.Date(2024, 8, 2, 0, 0, 0, 0, time.UTC), Amount: decimal.RequireFromString("20.46"), Type: Debit},
		{ID: 4, Date: time.Date(2024, 8, 13, 0, 0, 0, 0, time.UTC), Amount: decimal.RequireFromString("10"), Type: Credit},
		{ID: 5, Date: time.Date(2024, 7, 18, 0, 0, 0, 0, time.UTC), Amount: decimal.RequireFromString("15.25"), Type: Credit},
	}

	summary := NewSummary(transactions)

	expectedBalance := decimal.RequireFromString("54.99")
	if !summary.TotalBalance.Equal(expectedBalance) {
		t.Errorf("TotalBalance = %s, want %s", summary.TotalBalance.String(), expectedBalance.String())
	}

	expectedJuly := 3
	expectedAugust := 2

	if summary.MonthlyTransactions["July 2024"] != expectedJuly {
		t.Errorf("July transactions = %d, want %d", summary.MonthlyTransactions["July 2024"], expectedJuly)
	}

	if summary.MonthlyTransactions["August 2024"] != expectedAugust {
		t.Errorf("August transactions = %d, want %d", summary.MonthlyTransactions["August 2024"], expectedAugust)
	}

	// Credits: 60.5 + 10 + 15.25 = 85.75 / 3 = 28.583333...
	expectedAvgCredit := decimal.RequireFromString("85.75").Div(decimal.NewFromInt(3))
	if !summary.AverageCredit.Equal(expectedAvgCredit) {
		t.Errorf("AverageCredit = %s, want %s", summary.AverageCredit.String(), expectedAvgCredit.String())
	}

	// Debits: 10.3 + 20.46 = 30.76 / 2 = 15.38
	expectedAvgDebit := decimal.RequireFromString("30.76").Div(decimal.NewFromInt(2))
	if !summary.AverageDebit.Equal(expectedAvgDebit) {
		t.Errorf("AverageDebit = %s, want %s", summary.AverageDebit.String(), expectedAvgDebit.String())
	}

	if !summary.HasTransactions() {
		t.Error("HasTransactions() should return true when transactions exist")
	}
}

func TestNewSummaryEmpty(t *testing.T) {
	var transactions []Transaction
	summary := NewSummary(transactions)

	if !summary.TotalBalance.IsZero() {
		t.Errorf("TotalBalance should be zero for empty transactions, got %s", summary.TotalBalance.String())
	}

	if len(summary.MonthlyTransactions) != 0 {
		t.Errorf("MonthlyTransactions should be empty, got %d entries", len(summary.MonthlyTransactions))
	}

	if !summary.AverageCredit.IsZero() {
		t.Errorf("AverageCredit should be zero for empty transactions, got %s", summary.AverageCredit.String())
	}

	if !summary.AverageDebit.IsZero() {
		t.Errorf("AverageDebit should be zero for empty transactions, got %s", summary.AverageDebit.String())
	}

	if summary.HasTransactions() {
		t.Error("HasTransactions() should return false for empty transactions")
	}
}

func TestNewSummaryOnlyCredits(t *testing.T) {
	transactions := []Transaction{
		{ID: 1, Date: time.Date(2024, 7, 15, 0, 0, 0, 0, time.UTC), Amount: decimal.RequireFromString("50"), Type: Credit},
		{ID: 2, Date: time.Date(2024, 7, 20, 0, 0, 0, 0, time.UTC), Amount: decimal.RequireFromString("30"), Type: Credit},
	}

	summary := NewSummary(transactions)

	// Balance should be sum of credits
	expectedBalance := decimal.RequireFromString("80")
	if !summary.TotalBalance.Equal(expectedBalance) {
		t.Errorf("TotalBalance = %s, want %s", summary.TotalBalance.String(), expectedBalance.String())
	}

	// Average debit should be zero
	if !summary.AverageDebit.IsZero() {
		t.Errorf("AverageDebit should be zero when no debits exist, got %s", summary.AverageDebit.String())
	}

	// Average credit should be calculated
	expectedAvgCredit := decimal.RequireFromString("40") // (50 + 30) / 2
	if !summary.AverageCredit.Equal(expectedAvgCredit) {
		t.Errorf("AverageCredit = %s, want %s", summary.AverageCredit.String(), expectedAvgCredit.String())
	}
}

func TestNewSummaryOnlyDebits(t *testing.T) {
	transactions := []Transaction{
		{ID: 1, Date: time.Date(2024, 7, 15, 0, 0, 0, 0, time.UTC), Amount: decimal.RequireFromString("25"), Type: Debit},
		{ID: 2, Date: time.Date(2024, 7, 20, 0, 0, 0, 0, time.UTC), Amount: decimal.RequireFromString("15"), Type: Debit},
	}

	summary := NewSummary(transactions)

	// Balance should be negative sum of debits
	expectedBalance := decimal.RequireFromString("-40")
	if !summary.TotalBalance.Equal(expectedBalance) {
		t.Errorf("TotalBalance = %s, want %s", summary.TotalBalance.String(), expectedBalance.String())
	}

	// Average credit should be zero
	if !summary.AverageCredit.IsZero() {
		t.Errorf("AverageCredit should be zero when no credits exist, got %s", summary.AverageCredit.String())
	}

	// Average debit should be calculated
	expectedAvgDebit := decimal.RequireFromString("20") // (25 + 15) / 2
	if !summary.AverageDebit.Equal(expectedAvgDebit) {
		t.Errorf("AverageDebit = %s, want %s", summary.AverageDebit.String(), expectedAvgDebit.String())
	}
}
