ALTER TABLE domain_counts
DROP CONSTRAINT IF EXISTS domain_user_unique;

ALTER TABLE domain_counts
DROP COLUMN IF EXISTS user_id;

-- restore original uniqueness (global)
ALTER TABLE domain_counts
ADD CONSTRAINT domain_counts_domain_key UNIQUE (domain);
