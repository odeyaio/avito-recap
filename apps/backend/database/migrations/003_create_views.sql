CREATE TABLE views (
    id BIGSERIAL PRIMARY KEY,

    user_id BIGINT NOT NULL
        REFERENCES users(id),

    listing_id BIGINT NOT NULL
        REFERENCES listings(id),

    category TEXT NOT NULL,

    viewed_at TIMESTAMPTZ NOT NULL,

    duration_seconds INTEGER NOT NULL
        CHECK (duration_seconds >= 0),

    is_repeat BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE INDEX idx_views_user_id
    ON views (user_id);

CREATE INDEX idx_views_listing_id
    ON views (listing_id);

CREATE INDEX idx_views_viewed_at
    ON views (viewed_at);

CREATE INDEX idx_views_user_category
    ON views (user_id, category);
