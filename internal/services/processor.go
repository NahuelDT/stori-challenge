package services

import (
	"context"
	"fmt"
	"log/slog"
)

type TransactionProcessor struct {
	fileProcessor FileProcessor
	emailService  EmailService
	calculator    SummaryCalculator
	dataStore     DataStore
	logger        *slog.Logger
}

// NewTransactionProcessor creates a new transaction processor
func NewTransactionProcessor(
	fileProcessor FileProcessor,
	emailService EmailService,
	calculator SummaryCalculator,
	dataStore DataStore,
	logger *slog.Logger,
) *TransactionProcessor {
	return &TransactionProcessor{
		fileProcessor: fileProcessor,
		emailService:  emailService,
		calculator:    calculator,
		dataStore:     dataStore,
		logger:        logger,
	}
}

// ProcessTransactionFile processes a single transaction file
func (p *TransactionProcessor) ProcessTransactionFile(ctx context.Context, filePath, recipientEmail string) error {
	p.logger.Info("processing transaction file", "file", filePath, "recipient", recipientEmail)

	// Process the file
	transactions, err := p.fileProcessor.ProcessFile(ctx, filePath)
	if err != nil {
		p.logger.Error("failed to process file", "error", err, "file", filePath)
		return fmt.Errorf("processing file %s: %w", filePath, err)
	}

	if len(transactions) == 0 {
		p.logger.Warn("no transactions found in file", "file", filePath)
		return fmt.Errorf("no transactions found in file %s", filePath)
	}

	p.logger.Info("transactions processed", "count", len(transactions), "file", filePath)

	// Calculate summary
	summary := p.calculator.Calculate(transactions)

	// Save to database if datastore is available
	if p.dataStore != nil {
		if err := p.dataStore.SaveTransactions(ctx, transactions); err != nil {
			p.logger.Error("failed to save transactions to database", "error", err)
			// Don't fail the entire process if database save fails
		} else {
			p.logger.Info("transactions saved to database", "count", len(transactions))
		}
	}

	// Send email summary
	if err := p.emailService.SendSummary(ctx, recipientEmail, summary); err != nil {
		p.logger.Error("failed to send email summary", "error", err, "recipient", recipientEmail)
		return fmt.Errorf("sending email to %s: %w", recipientEmail, err)
	}

	p.logger.Info("file succesfully processed", "recipient", recipientEmail)
	return nil
}

// WatchAndProcess watches a directory for new files and processes them
func (p *TransactionProcessor) WatchAndProcess(ctx context.Context, dirPath, recipientEmail string) error {
	p.logger.Info("starting directory watch", "directory", dirPath, "recipient", recipientEmail)

	fileChan, err := p.fileProcessor.WatchDirectory(ctx, dirPath)
	if err != nil {
		return fmt.Errorf("watching directory %s: %w", dirPath, err)
	}

	for {
		select {
		case <-ctx.Done():
			p.logger.Info("stopping directory watch due to context cancellation")
			return ctx.Err()
		case filePath := <-fileChan:
			p.logger.Info("new file detected", "file", filePath)
			if err := p.ProcessTransactionFile(ctx, filePath, recipientEmail); err != nil {
				p.logger.Error("failed to process detected file", "error", err, "file", filePath)
			}
		}
	}
}
