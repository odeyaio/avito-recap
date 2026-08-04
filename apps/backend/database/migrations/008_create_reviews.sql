CREATE TABLE reviews (
    id BIGSERIAL PRIMARY KEY,

    reviewer_id BIGINT NOT NULL
        REFERENCES users(id),

    reviewed_user_id BIGINT NOT NULL
        REFERENCES users(id),

    deal_id BIGINT
        REFERENCES deals(id),

    rating SMALLINT NOT NULL
        CHECK (rating BETWEEN 1 AND 5),

    created_at TIMESTAMPTZ NOT NULL,

    CHECK (reviewer_id <> reviewed_user_id)
);

CREATE INDEX idx_reviews_reviewer_id
    ON reviews (reviewer_id);

CREATE INDEX idx_reviews_reviewed_user_id
    ON reviews (reviewed_user_id);

CREATE INDEX idx_reviews_deal_id
    ON reviews (deal_id);
