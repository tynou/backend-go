-- name: CreateWallet :exec
INSERT INTO wallets (user_id) VALUES ($1);

-- name: DeductBalance :execrows
UPDATE wallets
SET balance = balance - $1
WHERE user_id = $2 AND balance >= $1;