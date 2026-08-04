CREATE TABLE user_events (
    id BIGSERIAL PRIMARY KEY,

    user_id BIGINT NOT NULL
        REFERENCES users(id),

    event_type TEXT NOT NULL,

    occurred_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX idx_user_events_user_time
    ON user_events (user_id, occurred_at);

CREATE INDEX idx_user_events_time
    ON user_events (occurred_at);
