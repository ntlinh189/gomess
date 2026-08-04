package database

import "database/sql"

type TxFunc func(*sql.Tx) error

func (d *Database) WithTransaction(fn TxFunc) error {
	tx, err := d.db.Begin()

	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
	}()

	if err := fn(tx); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}