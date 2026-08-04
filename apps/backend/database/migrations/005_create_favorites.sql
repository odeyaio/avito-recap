CREATE TABLE favorites (
    id BIGSERIAL PRIMARY KEY,

    user_id BIGINT NOT NULL
        REFERENCES users(id),

    listing_id BIGINT NOT NULL
        REFERENCES listings(id),

    category TEXT NOT NULL,

    action TEXT NOT NULL
        CHECK (action IN ('add', 'remove')),

    occurred_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX idx_favorites_user_id
    ON favorites (user_id);

CREATE INDEX idx_favorites_listing_id
    ON favorites (listing_id);

CREATE INDEX idx_favorites_occurred_at
    ON favorites (occurred_at);
