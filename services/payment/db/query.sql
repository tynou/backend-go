-- name: CreatePayment :one
INSERT INTO payments (user_id, amount, status)
VALUES ($1, $2, $3)
RETURNING *;

-- name: UpdatePaymentStatus :exec
UPDATE payments
SET status = $2
WHERE id = $1
RETURNING *;

-- name: GetPayment :one
SELECT * FROM payments
WHERE id = $1 LIMIT 1;