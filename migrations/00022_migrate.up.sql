BEGIN;

ALTER TABLE notice ADD COLUMN IF NOT EXISTS collectible_id bigint DEFAULT NULL;
ALTER TABLE notice ADD COLUMN IF NOT EXISTS type int  DEFAULT 0; 
ALTER TABLE notice ADD COLUMN IF NOT EXISTS from_account bigint DEFAULT NULL;

ALTER TABLE notice ADD FOREIGN KEY (collectible_id) REFERENCES collectible(id);
ALTER TABLE notice ADD FOREIGN KEY (from_account) REFERENCES account(id);

COMMIT;
