package services

import "github.com/tartoide/stori/stori-challenge/internal/domain"

type summaryCalculator struct{}

// NewSummaryCalculator creates a new summary calculator
func NewSummaryCalculator() SummaryCalculator {
	return &summaryCalculator{}
}

// Calculate computes summary statistics from transactions
func (c *summaryCalculator) Calculate(transactions []domain.Transaction) *domain.Summary {
	return domain.NewSummary(transactions)
}
