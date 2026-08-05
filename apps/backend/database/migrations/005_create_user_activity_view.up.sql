CREATE VIEW user_activity AS
SELECT
    events.user_id,
    events.occurred_at,
    events.event_type,
    events.listing_id,
    COALESCE(events.category_id, listings.category_id) AS category_id,
    'actor'::TEXT AS user_role
FROM activity_events AS events
LEFT JOIN listings
    ON listings.id = events.listing_id

UNION ALL

SELECT
    listings.seller_id AS user_id,
    listings.published_at AS occurred_at,
    'listing_published'::TEXT AS event_type,
    listings.id AS listing_id,
    listings.category_id,
    'seller'::TEXT AS user_role
FROM listings

UNION ALL

SELECT
    listings.seller_id AS user_id,
    listings.closed_at AS occurred_at,
    'listing_closed'::TEXT AS event_type,
    listings.id AS listing_id,
    listings.category_id,
    'seller'::TEXT AS user_role
FROM listings
WHERE listings.closed_at IS NOT NULL

UNION ALL

SELECT
    deals.buyer_id AS user_id,
    deals.completed_at AS occurred_at,
    'deal_completed'::TEXT AS event_type,
    deals.listing_id,
    listings.category_id,
    'buyer'::TEXT AS user_role
FROM deals
JOIN listings
    ON listings.id = deals.listing_id
WHERE deals.status = 'completed'
  AND deals.completed_at IS NOT NULL

UNION ALL

SELECT
    listings.seller_id AS user_id,
    deals.completed_at AS occurred_at,
    'deal_completed'::TEXT AS event_type,
    deals.listing_id,
    listings.category_id,
    'seller'::TEXT AS user_role
FROM deals
JOIN listings
    ON listings.id = deals.listing_id
WHERE deals.status = 'completed'
  AND deals.completed_at IS NOT NULL

UNION ALL

SELECT
    reviews.author_id AS user_id,
    reviews.created_at AS occurred_at,
    'review_left'::TEXT AS event_type,
    deals.listing_id,
    listings.category_id,
    'author'::TEXT AS user_role
FROM reviews
LEFT JOIN deals
    ON deals.id = reviews.deal_id
LEFT JOIN listings
    ON listings.id = deals.listing_id;
