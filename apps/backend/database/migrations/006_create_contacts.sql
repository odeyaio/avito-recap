CREATE TABLE contacts (
    id BIGSERIAL PRIMARY KEY,

    user_id BIGINT NOT NULL
        REFERENCES users(id),

    listing_id BIGINT NOT NULL
        REFERENCES listings(id),

    contact_type TEXT NOT NULL
        CHECK (contact_type IN ('chat', 'call', 'offer')),

    occurred_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX idx_contacts_user_id
    ON contacts (user_id);

CREATE INDEX idx_contacts_listing_id
    ON contacts (listing_id);

CREATE INDEX idx_contacts_occurred_at
    ON contacts (occurred_at);
