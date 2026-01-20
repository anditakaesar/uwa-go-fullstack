CREATE UNIQUE INDEX gift_ratings_active
ON gift_ratings(user_id, gift_id)
WHERE deleted_at IS NULL;