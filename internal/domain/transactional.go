package domain

import (
	"context"

	"cloud.google.com/go/firestore"
)

type Transactional struct {
	client *firestore.Client
}

func NewTransactional(client *firestore.Client) *Transactional {
	return &Transactional{
		client: client,
	}
}

func (t *Transactional) RunFirestoreTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	err := t.client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		txCtx := context.WithValue(ctx, "firestoreTransaction", tx)
		return fn(txCtx)
	})

	return err
}

func GetFirestoreTransaction(ctx context.Context) *firestore.Transaction {
	tx, ok := ctx.Value("firestoreTransaction").(*firestore.Transaction)
	if !ok {
		return nil
	}
	return tx
}
