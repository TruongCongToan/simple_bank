-- name: CreateTransfer :one
INSERT INTO tranfers (
	from_account_id,
	to_account_id,
	amount
) VALUES (
	$1, $2, $3
) RETURNING *;

-- name: GetTranfers :one
SELECT * FROM tranfers WHERE id = $1 LIMIT 1;

-- name: ListTranfers :many
SELECT * FROM tranfers ORDER BY id LIMIT $1 OFFSET $2; 

-- name: UpdateTranfers :one
UPDATE tranfers SET amount = $2 WHERE id = $1 RETURNING *;

-- name: DeleteTranfers :exec
DELETE FROM tranfers WHERE id = $1;