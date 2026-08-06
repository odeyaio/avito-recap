CREATE TABLE recap_snapshots (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users (id),
    period_start DATE NOT NULL,
    period_end DATE NOT NULL,
    ruleset_version TEXT NOT NULL,
    dataset_version TEXT NOT NULL,
    data_cutoff_at TIMESTAMPTZ NOT NULL,
    metrics JSONB NOT NULL,
    generated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CHECK (period_end > period_start),
    UNIQUE (
        user_id,
        period_start,
        period_end,
        ruleset_version,
        dataset_version
    )
);

CREATE INDEX idx_recap_snapshots_user_generated_at
    ON recap_snapshots (user_id, generated_at DESC);
