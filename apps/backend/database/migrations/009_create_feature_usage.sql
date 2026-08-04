CREATE TABLE feature_usage (
    id BIGSERIAL PRIMARY KEY,

    user_id BIGINT NOT NULL
        REFERENCES users(id),

    feature TEXT NOT NULL
        CHECK (
            feature IN (
                'search',
                'notifications',
                'delivery',
                'promotion'
            )
        ),

    action TEXT NOT NULL
        CHECK (
            action IN (
                'used',
                'enabled',
                'disabled'
            )
        ),

    occurred_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX idx_feature_usage_user_id
    ON feature_usage (user_id);

CREATE INDEX idx_feature_usage_feature
    ON feature_usage (feature);

CREATE INDEX idx_feature_usage_occurred_at
    ON feature_usage (occurred_at);
