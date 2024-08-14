-- name: CreateAccount :execresult
INSERT INTO accounts (
    id, owner, balance, currency, created_at
) VALUES (
    DEFAULT, ?, ?, ?, DEFAULT
);

-- name: GetAccount :one
SELECT * FROM accounts
WHERE id = ? LIMIT 1;

-- name: GetAccountForUpdate :one
SELECT * FROM accounts
WHERE id = ? LIMIT 1
FOR UPDATE;

-- name: ListAccounts :many
SELECT * FROM accounts
WHERE owner = ?
ORDER BY id
LIMIT ?
OFFSET ?;

-- name: UpdateAccount :execresult
UPDATE accounts
SET balance = ?
WHERE id = ?;

-- name: AddAccountBalance :execresult
UPDATE accounts
SET balance = balance + ?
WHERE id = ?;

-- name: DeleteAccount :exec
DELETE FROM accounts
WHERE id = ?;