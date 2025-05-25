package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"github.com/shopspring/decimal"
	"github.com/tartoide/stori/stori-challenge/internal/domain"
)

type TransactionRepository struct {
	db     *sql.DB
	loader *QueryLoader
	logger *slog.Logger
}

// NewTransactionRepository creates a new transaction repository
func NewTransactionRepository(db *sql.DB, loader *QueryLoader, logger *slog.Logger) *TransactionRepository {
	return &TransactionRepository{
		db:     db,
		loader: loader,
		logger: logger,
	}
}

// SaveBatch saves multiple transactions in a single database transaction
func (r *TransactionRepository) SaveBatch(ctx context.Context, accountID string, transactions []domain.Transaction) error {
	if len(transactions) == 0 {
		return nil
	}

	r.logger.Debug("saving transaction batch", "count", len(transactions), "account_id", accountID)

	query, err := r.loader.GetQuery("InsertTransaction")
	if err != nil {
		return fmt.Errorf("getting insert transaction query: %w", err)
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("beginning transaction: %w", err)
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("preparing insert statement: %w", err)
	}
	defer stmt.Close()

	for _, transaction := range transactions {
		amount := transaction.Amount
		if transaction.Type == domain.Debit {
			amount = amount.Neg()
		}

		_, err := stmt.ExecContext(ctx,
			transaction.ID,
			accountID,
			transaction.Date,
			amount,
			transaction.Type.String(),
		)
		if err != nil {
			return fmt.Errorf("inserting transaction %d: %w", transaction.ID, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("committing transaction: %w", err)
	}

	r.logger.Info("transaction batch saved", "count", len(transactions), "account_id", accountID)
	return nil
}

// GetBalance retrieves the current balance for an account
func (r *TransactionRepository) GetBalance(ctx context.Context, accountID string) (decimal.Decimal, error) {
	r.logger.Debug("retrieving account balance", "account_id", accountID)

	query, err := r.loader.GetQuery("GetAccountBalance")
	if err != nil {
		return decimal.Zero, fmt.Errorf("getting balance query: %w", err)
	}

	var balance decimal.Decimal
	err = r.db.QueryRowContext(ctx, query, accountID).Scan(&balance)
	if err != nil {
		if err == sql.ErrNoRows {
			return decimal.Zero, nil
		}
		return decimal.Zero, fmt.Errorf("querying account balance: %w", err)
	}

	return balance, nil
}

// GetByAccount retrieves all transactions for an account
func (r *TransactionRepository) GetByAccount(ctx context.Context, accountID string) ([]domain.Transaction, error) {
	r.logger.Debug("retrieving transactions by account", "account_id", accountID)

	query, err := r.loader.GetQuery("GetTransactionsByAccount")
	if err != nil {
		return nil, fmt.Errorf("getting transactions by account query: %w", err)
	}

	rows, err := r.db.QueryContext(ctx, query, accountID)
	if err != nil {
		return nil, fmt.Errorf("querying transactions by account: %w", err)
	}
	defer rows.Close()

	return r.scanTransactions(rows)
}

// GetByDateRange retrieves transactions for an account within a date range
func (r *TransactionRepository) GetByDateRange(ctx context.Context, accountID string, startDate, endDate time.Time) ([]domain.Transaction, error) {
	r.logger.Debug("retrieving transactions by date range", 
		"account_id", accountID, 
		"start_date", startDate.Format("2006-01-02"), 
		"end_date", endDate.Format("2006-01-02"))

	query, err := r.loader.GetQuery("GetTransactionsByDateRange")
	if err != nil {
		return nil, fmt.Errorf("getting transactions by date range query: %w", err)
	}

	rows, err := r.db.QueryContext(ctx, query, accountID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("querying transactions by date range: %w", err)
	}
	defer rows.Close()

	return r.scanTransactions(rows)
}

func (r *TransactionRepository) scanTransactions(rows *sql.Rows) ([]domain.Transaction, error) {
	var transactions []domain.Transaction

	for rows.Next() {
		var (
			id              int
			accountID       string
			transactionDate time.Time
			amount          decimal.Decimal
			transactionType string
			processedAt     sql.NullTime
		)

		err := rows.Scan(&id, &accountID, &transactionDate, &amount, &transactionType, &processedAt)
		if err != nil {
			return nil, fmt.Errorf("scanning transaction row: %w", err)
		}

		var txType domain.TransactionType
		switch transactionType {
		case "credit":
			txType = domain.Credit
		case "debit":
			txType = domain.Debit
			amount = amount.Abs() // Store debits as positive amounts in domain
		default:
			r.logger.Warn("unknown transaction type", "type", transactionType, "id", id)
			continue
		}

		transaction := domain.Transaction{
			ID:     id,
			Date:   transactionDate,
			Amount: amount,
			Type:   txType,
		}

		transactions = append(transactions, transaction)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("scanning transaction rows: %w", err)
	}

	return transactions, nil
}