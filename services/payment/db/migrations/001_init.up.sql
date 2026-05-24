CREATE TYPE payment_status_type AS ENUM ('pending', 'success', 'failure');

CREATE TABLE IF NOT EXISTS payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id INT UNIQUE NOT NULL,
    amount NUMERIC(10, 2) NOT NULL,
    status payment_status_type NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);