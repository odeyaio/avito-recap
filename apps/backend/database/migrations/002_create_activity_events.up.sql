CREATE TABLE activity_events (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users (id),
    event_type TEXT NOT NULL,
    occurred_at TIMESTAMPTZ NOT NULL,
    listing_id UUID REFERENCES listings (id),
    category_id UUID REFERENCES categories (id),
    duration_seconds INTEGER,
    result_count INTEGER,
    filter_count INTEGER,
    topic_key TEXT,
    source_type TEXT,
    properties JSONB NOT NULL DEFAULT '{}'::JSONB,
    ingested_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CHECK (duration_seconds IS NULL OR duration_seconds >= 0),
    CHECK (result_count IS NULL OR result_count >= 0),
    CHECK (filter_count IS NULL OR filter_count >= 0)
);

CREATE INDEX idx_activity_events_user_time
    ON activity_events (user_id, occurred_at DESC);

CREATE INDEX idx_activity_events_user_type_time
    ON activity_events (user_id, event_type, occurred_at DESC);

CREATE INDEX idx_activity_events_user_category_time
    ON activity_events (user_id, category_id, occurred_at DESC)
    WHERE category_id IS NOT NULL;

CREATE INDEX idx_activity_events_listing_time
    ON activity_events (listing_id, occurred_at DESC)
    WHERE listing_id IS NOT NULL;
