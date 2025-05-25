package domain

import (
	"testing"
	"time"

	"github.com/shopspring/decimal"
)

func TestNewTransaction(t *testing.T) {
	tests := []struct {
		name        string
		id          string
		date        string
		amount      string
		wantID      int
		wantType    TransactionType
		wantAmount  string
		expectError bool
	}{
		{
			name:        "valid credit transaction",
			id:          "1",
			date:        "7/15",
			amount:      "+60.5",
			wantID:      1,
			wantType:    Credit,
			wantAmount:  "60.5",
			expectError: false,
		},
		{
			name:        "valid debit transaction",
			id:          "2",
			date:        "8/2",
			amount:      "-20.46",
			wantID:      2,
			wantType:    Debit,
			wantAmount:  "20.46",
			expectError: false,
		},
		{
			name:        "invalid transaction ID",
			id:          "invalid",
			date:        "7/15",
			amount:      "+60.5",
			expectError: true,
		},
		{
			name:        "invalid date format",
			id:          "1",
			date:        "invalid-date",
			amount:      "+60.5",
			expectError: true,
		},
		{
			name:        "invalid amount format",
			id:          "1",
			date:        "7/15",
			amount:      "invalid",
			expectError: true,
		},
		{
			name:        "missing amount sign",
			id:          "1",
			date:        "7/15",
			amount:      "60.5",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transaction, err := NewTransaction(tt.id, tt.date, tt.amount)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if transaction.ID != tt.wantID {
				t.Errorf("ID = %d, want %d", transaction.ID, tt.wantID)
			}

			if transaction.Type != tt.wantType {
				t.Errorf("Type = %v, want %v", transaction.Type, tt.wantType)
			}

			wantAmount := decimal.RequireFromString(tt.wantAmount)
			if !transaction.Amount.Equal(wantAmount) {
				t.Errorf("Amount = %s, want %s", transaction.Amount.String(), wantAmount.String())
			}
		})
	}
}

func TestTransactionMethods(t *testing.T) {
	creditTx, _ := NewTransaction("1", "7/15", "+60.5")
	debitTx, _ := NewTransaction("2", "8/2", "-20.46")

	if !creditTx.IsCredit() {
		t.Error("credit transaction should return true for IsCredit()")
	}
	if creditTx.IsDebit() {
		t.Error("credit transaction should return false for IsDebit()")
	}

	if debitTx.IsCredit() {
		t.Error("debit transaction should return false for IsCredit()")
	}
	if !debitTx.IsDebit() {
		t.Error("debit transaction should return true for IsDebit()")
	}

	// Test AbsoluteAmount
	expectedAmount := decimal.RequireFromString("60.5")
	if !creditTx.AbsoluteAmount().Equal(expectedAmount) {
		t.Errorf("AbsoluteAmount() = %s, want %s", creditTx.AbsoluteAmount().String(), expectedAmount.String())
	}

	expectedAmount = decimal.RequireFromString("20.46")
	if !debitTx.AbsoluteAmount().Equal(expectedAmount) {
		t.Errorf("AbsoluteAmount() = %s, want %s", debitTx.AbsoluteAmount().String(), expectedAmount.String())
	}

	currentYear := time.Now().Year()
	expectedMonthYear := "July " + string(rune(currentYear/1000)+'0') + string(rune((currentYear/100)%10)+'0') + string(rune((currentYear/10)%10)+'0') + string(rune(currentYear%10)+'0')
	if creditTx.MonthYear() != expectedMonthYear {
		t.Errorf("MonthYear() = %s, want %s", creditTx.MonthYear(), expectedMonthYear)
	}
}

func TestTransactionTypeString(t *testing.T) {
	tests := []struct {
		txType TransactionType
		want   string
	}{
		{Credit, "credit"},
		{Debit, "debit"},
		{TransactionType(999), "unknown"},
	}

	for _, tt := range tests {
		if got := tt.txType.String(); got != tt.want {
			t.Errorf("TransactionType.String() = %s, want %s", got, tt.want)
		}
	}
}
