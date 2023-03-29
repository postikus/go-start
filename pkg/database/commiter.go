package database

import "github.com/jmoiron/sqlx"

type Commiter struct {
	tx *sqlx.Tx
}

func NewCommitter() *Commiter {
	return new(Commiter)
}

func (c Commiter) Close(tx *sqlx.Tx) func(error) {
	c.tx = tx

	return func(err error) {
		if err != nil {
			c.tx.Rollback()
			return
		}

		c.tx.Commit()
	}
}
