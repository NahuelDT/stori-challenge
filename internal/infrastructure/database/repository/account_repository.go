package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
)

type Account struct {
	ID        string
	Email     string
	CreatedAt sql.NullTime
}

type AccountRepository struct {
	db     *sql.DB
	loader *QueryLoader
	logger *slog.Logger
}

// NewAccountRepository creates a new account repository
func NewAccountRepository(db *sql.DB, loader *QueryLoader, logger *slog.Logger) *AccountRepository {
	return &AccountRepository{
		db:     db,
		loader: loader,
		logger: logger,
	}
}

// Create creates a new account and returns the account ID
func (r *AccountRepository) Create(ctx context.Context, email string) (string, error) {
	r.logger.Debug("creating new account", "email", email)

	query, err := r.loader.GetQuery("CreateAccount")
	if err != nil {
		return "", fmt.Errorf("getting create account query: %w", err)
	}

	accountID := uuid.New().String()
	var returnedID string

	err = r.db.QueryRowContext(ctx, query, accountID, email).Scan(&returnedID)
	if err != nil {
		return "", fmt.Errorf("creating account: %w", err)
	}

	r.logger.Info("account created", "account_id", returnedID, "email", email)
	return returnedID, nil
}

// GetByEmail retrieves an account by email
func (r *AccountRepository) GetByEmail(ctx context.Context, email string) (*Account, error) {
	r.logger.Debug("retrieving account by email", "email", email)

	query, err := r.loader.GetQuery("GetAccountByEmail")
	if err != nil {
		return nil, fmt.Errorf("getting account by email query: %w", err)
	}

	var account Account
	err = r.db.QueryRowContext(ctx, query, email).Scan(&account.ID, &account.Email, &account.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Account not found
		}
		return nil, fmt.Errorf("querying account by email: %w", err)
	}

	return &account, nil
}

// GetByID retrieves an account by ID
func (r *AccountRepository) GetByID(ctx context.Context, accountID string) (*Account, error) {
	r.logger.Debug("retrieving account by ID", "account_id", accountID)

	query, err := r.loader.GetQuery("GetAccountByID")
	if err != nil {
		return nil, fmt.Errorf("getting account by ID query: %w", err)
	}

	var account Account
	err = r.db.QueryRowContext(ctx, query, accountID).Scan(&account.ID, &account.Email, &account.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Account not found
		}
		return nil, fmt.Errorf("querying account by ID: %w", err)
	}

	return &account, nil
}
