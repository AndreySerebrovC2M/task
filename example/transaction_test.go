package example

import (
	"testing"

	"github.com/google/uuid"
)

const (
	TestDBConn   = "user=dev password=dev host=localhost port=5432 dbname=example_db sslmode=disable connect_timeout=5"
	TestDBDriver = "postgres"
)

var store *Store

func init() {
	st, err := NewStore(TestDBDriver, TestDBConn)
	if err != nil {
		panic("failed to connect to db")
	}
	store = st
}

func TestStore_InsertTransaction(t *testing.T) {
	tr := Transaction{
		ID:     uuid.New().String(),
		Amount: 0,
		State:  Win,
	}
	op := NewTransactionOperation(store)
	if err := op.Run(tr); err == nil {
		return
	}
	t.Fail()
}

func TestStore_DuplicatedTransaction(t *testing.T) {
	tr := Transaction{
		ID:     uuid.New().String(),
		Amount: 0,
		State:  Win,
	}
	op := NewTransactionOperation(store)
	if err := op.Run(tr); err == nil {
		if err := op.Run(tr); err != nil && err.Error() == "transaction is exist " {
			return
		}
	}
	t.Fail()
}

func TestStore_NegativeTransaction(t *testing.T) {
	tr := Transaction{
		ID:     uuid.New().String(),
		Amount: 10.0,
		State:  Lost,
	}
	op := NewTransactionOperation(store)
	if err := op.Run(tr); err != nil && err.Error() == "total amount is negative " {
		return
	}
	t.Fail()
}
