package services

import (
	"context"

	"github.com/shopspring/decimal"
	"github.com/tartoide/stori/stori-challenge/internal/domain"
)

// FileProcessor handles file processing operations
type FileProcessor interface {
	ProcessFile(ctx context.Context, filePath string) ([]domain.Transaction, error)
	WatchDirectory(ctx context.Context, dirPath string) (<-chan string, error)
}

// EmailService handles email operations
type EmailService interface {
	SendSummary(ctx context.Context, recipient string, summary *domain.Summary) error
	RenderTemplate(summary *domain.Summary) (string, error)
}

// DataStore handles data persistence operations
type DataStore interface {
	SaveTransactions(ctx context.Context, transactions []domain.Transaction) error
	GetAccountBalance(ctx context.Context, accountID string) (decimal.Decimal, error)
	SaveAccount(ctx context.Context, email string) (string, error)
}

// SummaryCalculator handles summary calculation operations
type SummaryCalculator interface {
	Calculate(transactions []domain.Transaction) *domain.Summary
}
