CREATE TABLE catalog_versions (
    id UUID PRIMARY KEY,
    kind TEXT NOT NULL,
    version TEXT NOT NULL,
    content_hash TEXT NOT NULL,
    imported_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CHECK (kind IN ('achievement', 'behavior')),
    CHECK (content_hash ~ '^[0-9a-f]{64}$'),
    UNIQUE (kind, version)
);

CREATE TABLE behavior_type_definitions (
    id UUID PRIMARY KEY,
    code TEXT NOT NULL,
    name TEXT NOT NULL,
    rule_description TEXT NOT NULL,
    rule JSONB NOT NULL,
    catalog_version_id UUID NOT NULL REFERENCES catalog_versions (id),
    default_action JSONB NOT NULL,
    enabled BOOLEAN NOT NULL DEFAULT TRUE,
    sort_order INTEGER NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    UNIQUE (catalog_version_id, code)
);

CREATE TABLE achievement_definitions (
    id UUID PRIMARY KEY,
    code TEXT NOT NULL,
    name TEXT NOT NULL,
    rule_description TEXT NOT NULL,
    rule JSONB NOT NULL,
    catalog_version_id UUID NOT NULL REFERENCES catalog_versions (id),
    icon_key TEXT NOT NULL,
    enabled BOOLEAN NOT NULL DEFAULT TRUE,
    shareable_by_default BOOLEAN NOT NULL DEFAULT FALSE,
    sort_order INTEGER NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    UNIQUE (catalog_version_id, code)
);

CREATE TABLE recap_behavior_types (
    recap_id UUID NOT NULL
        REFERENCES recap_snapshots (id) ON DELETE CASCADE,
    behavior_type_definition_id UUID NOT NULL
        REFERENCES behavior_type_definitions (id),
    is_primary BOOLEAN NOT NULL DEFAULT FALSE,
    position INTEGER NOT NULL,
    score DOUBLE PRECISION,
    evidence JSONB NOT NULL DEFAULT '{}'::JSONB,

    PRIMARY KEY (recap_id, behavior_type_definition_id)
);

CREATE UNIQUE INDEX idx_recap_behavior_types_one_primary
    ON recap_behavior_types (recap_id)
    WHERE is_primary;

CREATE INDEX idx_recap_behavior_types_order
    ON recap_behavior_types (recap_id, position);

CREATE TABLE recap_achievements (
    recap_id UUID NOT NULL
        REFERENCES recap_snapshots (id) ON DELETE CASCADE,
    achievement_definition_id UUID NOT NULL
        REFERENCES achievement_definitions (id),
    position INTEGER NOT NULL,
    achieved_at TIMESTAMPTZ,
    evidence JSONB NOT NULL DEFAULT '{}'::JSONB,
    is_shareable BOOLEAN NOT NULL DEFAULT FALSE,

    PRIMARY KEY (recap_id, achievement_definition_id)
);

CREATE INDEX idx_recap_achievements_order
    ON recap_achievements (recap_id, position);

CREATE TABLE recap_next_action (
    recap_id UUID PRIMARY KEY
        REFERENCES recap_snapshots (id) ON DELETE CASCADE,
    code TEXT NOT NULL,
    href TEXT NOT NULL,
    target JSONB NOT NULL,
    evidence JSONB NOT NULL DEFAULT '{}'::JSONB,
    resolved_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE recap_presentations (
    id UUID PRIMARY KEY,
    recap_id UUID NOT NULL
        REFERENCES recap_snapshots (id) ON DELETE CASCADE,
    locale TEXT NOT NULL,
    prompt_version TEXT NOT NULL,
    model_name TEXT NOT NULL,
    input_hash TEXT NOT NULL,
    content JSONB NOT NULL,
    generated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    UNIQUE (
        recap_id,
        locale,
        prompt_version,
        model_name,
        input_hash
    )
);

CREATE INDEX idx_recap_presentations_recap_generated_at
    ON recap_presentations (recap_id, generated_at DESC);
