-- Rollback: drop new PK and restore old one
ALTER TABLE domain_counts
DROP CONSTRAINT IF EXISTS domain_counts_pkey;

ALTER TABLE domain_counts
ADD PRIMARY KEY (domain);
