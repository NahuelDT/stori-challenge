package database

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"log/slog"

	"github.com/NahuelDT/stori-challenge/internal/domain"
	"github.com/NahuelDT/stori-challenge/internal/infrastructure/database/repository"
	"github.com/NahuelDT/stori-challenge/internal/services"
	_ "github.com/lib/pq"
	"github.com/shopspring/decimal"
)

//go:embed migrations/*.sql
var migrationFiles embed.FS

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type postgresDataStore struct {
	db              *sql.DB
	accountRepo     *repository.AccountRepository
	transactionRepo *repository.TransactionRepository
	logger          *slog.Logger
}

// NewPostgresDataStore creates a new PostgreSQL data store
func NewPostgresDataStore(config PostgresConfig, logger *slog.Logger) (services.DataStore, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("opening database connection: %w", err)
	}

	// Test
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("pinging database: %w", err)
	}

	queryLoader, err := repository.NewQueryLoader()
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("creating query loader: %w", err)
	}

	accountRepo := repository.NewAccountRepository(db, queryLoader, logger)
	transactionRepo := repository.NewTransactionRepository(db, queryLoader, logger)

	store := &postgresDataStore{
		db:              db,
		accountRepo:     accountRepo,
		transactionRepo: transactionRepo,
		logger:          logger,
	}

	if err := store.runMigrations(); err != nil {
		db.Close()
		return nil, fmt.Errorf("running migrations: %w", err)
	}

	logger.Info("PostgreSQL connection established successfully")
	return store, nil
}

// SaveTransactions saves transactions to the database
func (p *postgresDataStore) SaveTransactions(ctx context.Context, transactions []domain.Transaction) error {
	if len(transactions) == 0 {
		return nil
	}

	// default
	defaultEmail := "default@stori.com"

	account, err := p.accountRepo.GetByEmail(ctx, defaultEmail)
	if err != nil {
		return fmt.Errorf("checking account existence: %w", err)
	}

	var accountID string
	if account == nil {
		accountID, err = p.accountRepo.Create(ctx, defaultEmail)
		if err != nil {
			return fmt.Errorf("creating default account: %w", err)
		}
	} else {
		accountID = account.ID
	}

	return p.transactionRepo.SaveBatch(ctx, accountID, transactions)
}

// GetAccountBalance returns the account balance for a given account ID
func (p *postgresDataStore) GetAccountBalance(ctx context.Context, accountID string) (decimal.Decimal, error) {
	return p.transactionRepo.GetBalance(ctx, accountID)
}

// SaveAccount creates a new account and returns the account ID
func (p *postgresDataStore) SaveAccount(ctx context.Context, email string) (string, error) {
	return p.accountRepo.Create(ctx, email)
}

// Close closes the database connection
func (p *postgresDataStore) Close() error {
	return p.db.Close()
}

func (p *postgresDataStore) runMigrations() error {
	p.logger.Debug("running database migrations")

	migrations, err := migrationFiles.ReadDir("migrations")
	if err != nil {
		return fmt.Errorf("reading migrations directory: %w", err)
	}

	for _, migration := range migrations {
		if !migration.IsDir() {
			migrationPath := "migrations/" + migration.Name()
			p.logger.Debug("applying migration", "file", migrationPath)

			content, err := migrationFiles.ReadFile(migrationPath)
			if err != nil {
				return fmt.Errorf("reading migration file %s: %w", migrationPath, err)
			}

			_, err = p.db.Exec(string(content))
			if err != nil {
				return fmt.Errorf("executing migration %s: %w", migrationPath, err)
			}

			p.logger.Info("migration applied successfully", "file", migrationPath)
		}
	}

	p.logger.Info("all database migrations completed successfully")
	return nil
}
