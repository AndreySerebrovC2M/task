package example

import (
	"database/sql"
	_ "github.com/lib/pq"
)

type Store struct {
	dbConn *sql.DB
}

func NewStore(driver string, conn string) (*Store, error) {
	dbConn, err := raiseDBConnection(driver, conn)
	if err != nil {
		return nil, err
	}
	return &Store{dbConn: dbConn}, nil
}

func raiseDBConnection(driver string, connParams string) (dbConn *sql.DB, err error) {
	return sql.Open(driver, connParams)
}

func (store *Store) Close() (err error) {
	err = store.dbConn.Close()
	return err
}

func (store *Store) InsertTransaction(tr Transaction) (err error) {
	return NewTransactionOperation(store).Run(tr)
}
