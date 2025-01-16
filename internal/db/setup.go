package db

import (
	"context"
	"errors"
)

// DefaultServer relies on this being initialized (and not overwritten)
var DB = new(Queries)

func (q *Queries) Connected() bool {
	return q.db != nil || q.tx != nil
}

func Setup(ctx context.Context) error {
	if !Connected() {
		return errors.New("lsp DB: no database connection, bailing")
	}

	qs, err := Prepare(ctx, Conn)
	if err != nil {
		return err
	}
	*DB = *qs
	return nil
}
