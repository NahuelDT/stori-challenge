package domain

import (
	"github.com/shopspring/decimal"
)

type Summary struct {
	TotalBalance        decimal.Decimal            `json:"total_balance"`
	MonthlyTransactions map[string]int             `json:"monthly_transactions"`
	MonthlyCredits      map[string]decimal.Decimal `json:"monthly_credits"`
	MonthlyDebits       map[string]decimal.Decimal `json:"monthly_debits"`
	AverageCredit       decimal.Decimal            `json:"average_credit"`
	AverageDebit        decimal.Decimal            `json:"average_debit"`
}

// NewSummary creates a new summary from a list of transactions
func NewSummary(transactions []Transaction) *Summary {
	summary := &Summary{
		TotalBalance:        decimal.Zero,
		MonthlyTransactions: make(map[string]int),
		MonthlyCredits:      make(map[string]decimal.Decimal),
		MonthlyDebits:       make(map[string]decimal.Decimal),
		AverageCredit:       decimal.Zero,
		AverageDebit:        decimal.Zero,
	}

	if len(transactions) == 0 {
		return summary
	}

	var totalCredits, totalDebits decimal.Decimal
	var creditCount, debitCount int

	for _, transaction := range transactions {
		monthYear := transaction.MonthYear()

		summary.MonthlyTransactions[monthYear]++

		if transaction.IsCredit() {
			summary.TotalBalance = summary.TotalBalance.Add(transaction.Amount)
			summary.MonthlyCredits[monthYear] = summary.MonthlyCredits[monthYear].Add(transaction.Amount)
			totalCredits = totalCredits.Add(transaction.Amount)
			creditCount++
		} else {
			summary.TotalBalance = summary.TotalBalance.Sub(transaction.Amount)
			summary.MonthlyDebits[monthYear] = summary.MonthlyDebits[monthYear].Add(transaction.Amount)
			totalDebits = totalDebits.Add(transaction.Amount)
			debitCount++
		}
	}

	if creditCount > 0 {
		summary.AverageCredit = totalCredits.Div(decimal.NewFromInt(int64(creditCount)))
	}
	if debitCount > 0 {
		summary.AverageDebit = totalDebits.Div(decimal.NewFromInt(int64(debitCount)))
	}

	return summary
}

// HasTransactions returns true if there are any transactions in the summary
func (s *Summary) HasTransactions() bool {
	total := 0
	for _, count := range s.MonthlyTransactions {
		total += count
	}
	return total > 0
}
