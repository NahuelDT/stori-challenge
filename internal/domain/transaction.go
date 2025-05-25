package domain

import (
	"strconv"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

type TransactionType int

const (
	Credit TransactionType = iota
	Debit
)

func (t TransactionType) String() string {
	switch t {
	case Credit:
		return "credit"
	case Debit:
		return "debit"
	default:
		return "unknown"
	}
}

type Transaction struct {
	ID     int             `json:"id"`
	Date   time.Time       `json:"date"`
	Amount decimal.Decimal `json:"amount"`
	Type   TransactionType `json:"type"`
}

// NewTransaction creates a new transaction from CSV data
func NewTransaction(id, date, amount string) (*Transaction, error) {
	transactionID, err := strconv.Atoi(id)
	if err != nil {
		return nil, ErrInvalidTransactionID
	}

	parsedDate, err := parseDate(date)
	if err != nil {
		return nil, ErrInvalidDate
	}

	parsedAmount, transactionType, err := parseAmount(amount)
	if err != nil {
		return nil, ErrInvalidAmount
	}

	return &Transaction{
		ID:     transactionID,
		Date:   parsedDate,
		Amount: parsedAmount,
		Type:   transactionType,
	}, nil
}

// IsCredit returns true if the transaction is a credit
func (t *Transaction) IsCredit() bool {
	return t.Type == Credit
}

// IsDebit returns true if the transaction is a debit
func (t *Transaction) IsDebit() bool {
	return t.Type == Debit
}

// AbsoluteAmount returns the absolute value of the transaction amount
func (t *Transaction) AbsoluteAmount() decimal.Decimal {
	return t.Amount.Abs()
}

// MonthYear returns the month and year as a string (e.g., "July 2024")
func (t *Transaction) MonthYear() string {
	return t.Date.Format("January 2006")
}

func parseDate(dateStr string) (time.Time, error) {
	formats := []string{
		"1/2",        // M/D
		"1/2/06",     // M/D/YY
		"1/2/2006",   // M/D/YYYY
		"2006-01-02", // ISO format
	}

	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			// If year is missing, assume current year
			if t.Year() == 0 {
				now := time.Now()
				t = time.Date(now.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
			}
			return t, nil
		}
	}

	return time.Time{}, ErrInvalidDate
}

func parseAmount(amountStr string) (decimal.Decimal, TransactionType, error) {
	amountStr = strings.TrimSpace(amountStr)
	if amountStr == "" {
		return decimal.Zero, Credit, ErrInvalidAmount
	}

	var transactionType TransactionType
	if strings.HasPrefix(amountStr, "+") {
		transactionType = Credit
		amountStr = strings.TrimPrefix(amountStr, "+")
	} else if strings.HasPrefix(amountStr, "-") {
		transactionType = Debit
		amountStr = strings.TrimPrefix(amountStr, "-")
	} else {
		return decimal.Zero, Credit, ErrInvalidAmount
	}

	amount, err := decimal.NewFromString(amountStr)
	if err != nil {
		return decimal.Zero, Credit, ErrInvalidAmount
	}

	return amount.Abs(), transactionType, nil
}
