-- Drop old PK (domain only)
ALTER TABLE domain_counts
DROP CONSTRAINT IF EXISTS domain_counts_pkey;

-- Create new composite PK (domain, user_id)
ALTER TABLE domain_counts
ADD PRIMARY KEY (domain, user_id);
