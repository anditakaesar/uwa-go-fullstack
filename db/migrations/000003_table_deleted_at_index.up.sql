CREATE INDEX IF NOT EXISTS idx_users_active ON users(id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_gifts_active ON gifts(id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_gift_ratings_active ON gift_ratings(id) WHERE deleted_at IS NULL;