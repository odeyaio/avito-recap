CREATE TABLE listings (
    id BIGSERIAL PRIMARY KEY,

    seller_id BIGINT NOT NULL
        REFERENCES users(id),

    category TEXT NOT NULL,

    price NUMERIC(12, 2) NOT NULL
        CHECK (price >= 0),

    published_at TIMESTAMPTZ NOT NULL,
    closed_at TIMESTAMPTZ,

    delivery_available BOOLEAN NOT NULL DEFAULT FALSE,

    CHECK (closed_at IS NULL OR closed_at >= published_at)
);

CREATE INDEX idx_listings_seller_id
    ON listings (seller_id);

CREATE INDEX idx_listings_category
    ON listings (category);

CREATE INDEX idx_listings_published_at
    ON listings (published_at);
