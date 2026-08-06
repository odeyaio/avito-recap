CREATE TABLE deals (
    id UUID PRIMARY KEY,
    listing_id UUID NOT NULL REFERENCES listings (id),
    buyer_id UUID NOT NULL REFERENCES users (id),
    created_at TIMESTAMPTZ NOT NULL,
    completed_at TIMESTAMPTZ,
    status TEXT NOT NULL,
    delivery_used BOOLEAN NOT NULL DEFAULT FALSE,
    price_band TEXT NOT NULL
);

CREATE INDEX idx_deals_buyer_completed_at
    ON deals (buyer_id, completed_at DESC);

CREATE INDEX idx_deals_listing_completed_at
    ON deals (listing_id, completed_at DESC);

CREATE TABLE reviews (
    id UUID PRIMARY KEY,
    deal_id UUID REFERENCES deals (id),
    author_id UUID NOT NULL REFERENCES users (id),
    recipient_id UUID NOT NULL REFERENCES users (id),
    rating SMALLINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,

    CHECK (author_id <> recipient_id),
    CHECK (rating BETWEEN 1 AND 5)
);

CREATE UNIQUE INDEX idx_reviews_deal_author_unique
    ON reviews (deal_id, author_id)
    WHERE deal_id IS NOT NULL;

CREATE INDEX idx_reviews_author_created_at
    ON reviews (author_id, created_at DESC);

CREATE INDEX idx_reviews_recipient_created_at
    ON reviews (recipient_id, created_at DESC);
