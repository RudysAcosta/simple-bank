// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: transfer.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createTransfer = `-- name: CreateTransfer :one
INSERT INTO transfers (
  from_account_id,
  to_account_id,
  amount
) VALUES (
  $1, $2, $3
) RETURNING id, from_account_id, to_account_id, amount, created_at
`

type CreateTransferParams struct {
	FromAccountID pgtype.Int8 `json:"from_account_id"`
	ToAccountID   pgtype.Int8 `json:"to_account_id"`
	Amount        int64       `json:"amount"`
}

func (q *Queries) CreateTransfer(ctx context.Context, arg CreateTransferParams) (Transfer, error) {
	row := q.db.QueryRow(ctx, createTransfer, arg.FromAccountID, arg.ToAccountID, arg.Amount)
	var i Transfer
	err := row.Scan(
		&i.ID,
		&i.FromAccountID,
		&i.ToAccountID,
		&i.Amount,
		&i.CreatedAt,
	)
	return i, err
}

const getTransfer = `-- name: GetTransfer :one
SELECT id, from_account_id, to_account_id, amount, created_at FROM transfers
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetTransfer(ctx context.Context, id int64) (Transfer, error) {
	row := q.db.QueryRow(ctx, getTransfer, id)
	var i Transfer
	err := row.Scan(
		&i.ID,
		&i.FromAccountID,
		&i.ToAccountID,
		&i.Amount,
		&i.CreatedAt,
	)
	return i, err
}

const listTransfers = `-- name: ListTransfers :many
SELECT id, from_account_id, to_account_id, amount, created_at FROM transfers
WHERE from_account_id = $1 OR to_account_id = $1
ORDER BY id
LIMIT $2
OFFSET $3
`

type ListTransfersParams struct {
	FromAccountID pgtype.Int8 `json:"from_account_id"`
	Limit         int32       `json:"limit"`
	Offset        int32       `json:"offset"`
}

func (q *Queries) ListTransfers(ctx context.Context, arg ListTransfersParams) ([]Transfer, error) {
	rows, err := q.db.Query(ctx, listTransfers, arg.FromAccountID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Transfer
	for rows.Next() {
		var i Transfer
		if err := rows.Scan(
			&i.ID,
			&i.FromAccountID,
			&i.ToAccountID,
			&i.Amount,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listTransfersByFromAccount = `-- name: ListTransfersByFromAccount :many
SELECT id, from_account_id, to_account_id, amount, created_at FROM transfers
WHERE from_account_id = $1
ORDER BY id
LIMIT $2
OFFSET $3
`

type ListTransfersByFromAccountParams struct {
	FromAccountID pgtype.Int8 `json:"from_account_id"`
	Limit         int32       `json:"limit"`
	Offset        int32       `json:"offset"`
}

func (q *Queries) ListTransfersByFromAccount(ctx context.Context, arg ListTransfersByFromAccountParams) ([]Transfer, error) {
	rows, err := q.db.Query(ctx, listTransfersByFromAccount, arg.FromAccountID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Transfer
	for rows.Next() {
		var i Transfer
		if err := rows.Scan(
			&i.ID,
			&i.FromAccountID,
			&i.ToAccountID,
			&i.Amount,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listTransfersByToAccount = `-- name: ListTransfersByToAccount :many
SELECT id, from_account_id, to_account_id, amount, created_at FROM transfers
WHERE to_account_id = $1
ORDER BY id
LIMIT $2
OFFSET $3
`

type ListTransfersByToAccountParams struct {
	ToAccountID pgtype.Int8 `json:"to_account_id"`
	Limit       int32       `json:"limit"`
	Offset      int32       `json:"offset"`
}

func (q *Queries) ListTransfersByToAccount(ctx context.Context, arg ListTransfersByToAccountParams) ([]Transfer, error) {
	rows, err := q.db.Query(ctx, listTransfersByToAccount, arg.ToAccountID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Transfer
	for rows.Next() {
		var i Transfer
		if err := rows.Scan(
			&i.ID,
			&i.FromAccountID,
			&i.ToAccountID,
			&i.Amount,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const totalAmountReceivedByAccount = `-- name: TotalAmountReceivedByAccount :one
SELECT COALESCE(SUM(amount), 0) AS total_received
FROM transfers
WHERE to_account_id = $1
`

func (q *Queries) TotalAmountReceivedByAccount(ctx context.Context, toAccountID pgtype.Int8) (interface{}, error) {
	row := q.db.QueryRow(ctx, totalAmountReceivedByAccount, toAccountID)
	var total_received interface{}
	err := row.Scan(&total_received)
	return total_received, err
}

const totalAmountSentByAccount = `-- name: TotalAmountSentByAccount :one
SELECT COALESCE(SUM(amount), 0) AS total_sent
FROM transfers
WHERE from_account_id = $1
`

func (q *Queries) TotalAmountSentByAccount(ctx context.Context, fromAccountID pgtype.Int8) (interface{}, error) {
	row := q.db.QueryRow(ctx, totalAmountSentByAccount, fromAccountID)
	var total_sent interface{}
	err := row.Scan(&total_sent)
	return total_sent, err
}
