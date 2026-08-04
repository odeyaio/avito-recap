CREATE TABLE deals (
    id BIGSERIAL PRIMARY KEY,

    buyer_id BIGINT NOT NULL
        REFERENCES users(id),

    seller_id BIGINT NOT NULL
        REFERENCES users(id),

    listing_id BIGINT NOT NULL
        REFERENCES listings(id),

    category TEXT NOT NULL,

    price NUMERIC(12, 2) NOT NULL
        CHECK (price >= 0),

    delivery BOOLEAN NOT NULL DEFAULT FALSE,

    completed_at TIMESTAMPTZ NOT NULL,

    status TEXT NOT NULL DEFAULT 'completed'
        CHECK (status IN ('completed', 'cancelled'))
);

CREATE INDEX idx_deals_buyer_id
    ON deals (buyer_id);

CREATE INDEX idx_deals_seller_id
    ON deals (seller_id);

CREATE INDEX idx_deals_listing_id
    ON deals (listing_id);

CREATE INDEX idx_deals_completed_at
    ON deals (completed_at);
