CREATE TABLE users (
    id UUID PRIMARY KEY,
    display_name TEXT NOT NULL,
    registered_at TIMESTAMPTZ NOT NULL,
    region TEXT NOT NULL,
    timezone TEXT NOT NULL,
    is_test_profile BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_test_profiles
    ON users (display_name, id)
    WHERE is_test_profile;

CREATE TABLE categories (
    id UUID PRIMARY KEY,
    code TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL,
    parent_id UUID REFERENCES categories (id),

    CHECK (parent_id IS NULL OR parent_id <> id)
);

CREATE INDEX idx_categories_parent_id
    ON categories (parent_id)
    WHERE parent_id IS NOT NULL;

CREATE TABLE listings (
    id UUID PRIMARY KEY,
    seller_id UUID NOT NULL REFERENCES users (id),
    category_id UUID NOT NULL REFERENCES categories (id),
    region TEXT NOT NULL,
    price_band TEXT NOT NULL,
    published_at TIMESTAMPTZ NOT NULL,
    closed_at TIMESTAMPTZ,
    delivery_available BOOLEAN NOT NULL DEFAULT FALSE,
    photo_count INTEGER NOT NULL DEFAULT 0,
    description_complete BOOLEAN NOT NULL DEFAULT FALSE,

    CHECK (photo_count >= 0),
    CHECK (closed_at IS NULL OR closed_at >= published_at)
);

CREATE INDEX idx_listings_seller_published_at
    ON listings (seller_id, published_at DESC);

CREATE INDEX idx_listings_seller_closed_at
    ON listings (seller_id, closed_at DESC)
    WHERE closed_at IS NOT NULL;

CREATE INDEX idx_listings_category_published_at
    ON listings (category_id, published_at DESC);
