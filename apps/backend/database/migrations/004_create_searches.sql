CREATE TABLE searches (
    id BIGSERIAL PRIMARY KEY,

    user_id BIGINT NOT NULL
        REFERENCES users(id),

    searched_at TIMESTAMPTZ NOT NULL,

    topic TEXT NOT NULL,

    category TEXT,

    filters JSONB,

    result_count INTEGER NOT NULL
        CHECK (result_count >= 0)
);

CREATE INDEX idx_searches_user_id
    ON searches (user_id);

CREATE INDEX idx_searches_searched_at
    ON searches (searched_at);

CREATE INDEX idx_searches_category
    ON searches (category);

CREATE INDEX idx_searches_filters
    ON searches USING GIN (filters);
