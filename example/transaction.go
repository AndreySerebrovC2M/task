package example

import (
	"database/sql"
	"errors"
	_ "github.com/lib/pq"
)

const (
	Win  TransactionState = "win"
	Lost                  = "lost"
)

type TransactionState string

type Transaction struct {
	ID     string           `json:"transactionId"`
	Source string
	State  TransactionState `json:"state"`
	Amount float64          `json:"amount"`
}

type TransactionOperation struct {
	tx    *sql.Tx
	store *Store
}

func NewTransactionOperation(store *Store) *TransactionOperation {
	return &TransactionOperation{store: store}
}

func (op *TransactionOperation) Run(tr Transaction) (err error) {
	var total float64
	op.tx, err = op.store.dbConn.Begin()
	if err != nil {
		return
	}
	defer func() {
		if err == nil {
			_ = op.tx.Commit()
		} else {
			_ = op.tx.Rollback()
		}
	}()
	if total, err = op.total(); err != nil {
		return
	}
	if total = op.calcTotal(total, tr); total < 0 {
		return errors.New("total amount is negative")
	}
	exs, err := op.exists(tr.ID)
	if err != nil {
		return
	}
	if exs {
		return errors.New("transaction is exist")
	}
	if err = op.insert(tr); err != nil {
		return
	}
	return op.updateTotal(total)
}

func (op *TransactionOperation) total() (total float64, err error) {
	err = op.tx.QueryRow("SELECT total FROM account WHERE id = 1 FOR UPDATE").Scan(&total)
	return
}

func (op *TransactionOperation) updateTotal(total float64) (err error) {
	_, err = op.tx.Exec("UPDATE account SET total = $1 WHERE id = 1", total)
	return
}

func (op *TransactionOperation) calcTotal(total float64, tr Transaction) float64 {
	switch tr.State {
	case Win:
		total = total + tr.Amount
	case Lost:
		total = total - tr.Amount
	}
	return total
}

func (op *TransactionOperation) exists(id string) (bool, error) {
	res, err := op.tx.Exec("SELECT 1 FROM account_transaction WHERE id = $1", id)
	if err != nil {
		return false, err
	}
	if c, err := res.RowsAffected(); err != nil || c == 0 {
		return false, err
	}
	return true, nil
}

func (op *TransactionOperation) insert(tr Transaction) error {
	_, err := op.tx.Exec("INSERT INTO account_transaction (id, state, amount, source) VALUES ($1, $2, $3, $4)",
		tr.ID, tr.State, tr.Amount, tr.Source,
	)
	if err != nil {
		return err
	}
	return nil
}
