CREATE TABLE IF NOT EXISTS transactions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    merchant_id UUID NOT NULL REFERENCES merchants(id) ON DELETE CASCADE,
    order_id VARCHAR(255) NOT NULL UNIQUE,
    provider VARCHAR(255),
    amount NUMERIC(15, 2) NOT NULL,
    status VARCHAR(50) NOT NULL,
    external_ref VARCHAR(255) NOT NULL,
    redirect_url VARCHAR(255),
    raw_response TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);