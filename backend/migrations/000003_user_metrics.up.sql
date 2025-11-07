ALTER TABLE domain_counts
ADD COLUMN IF NOT EXISTS user_id INT REFERENCES users(id) ON DELETE CASCADE;

-- Drop any old uniqueness constraint on 'domain' only
ALTER TABLE domain_counts
DROP CONSTRAINT IF EXISTS domain_counts_domain_key;

-- Ensure uniqueness per user + domain
ALTER TABLE domain_counts
ADD CONSTRAINT domain_user_unique UNIQUE (domain, user_id);
