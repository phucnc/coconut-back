BEGIN;

ALTER TABLE exchange_event_buy_token ADD COLUMN IF NOT EXISTS type int DEFAULT -1;
ALTER TABLE exchange_event_buy_token ADD COLUMN IF NOT EXISTS account text DEFAULT NULL; 
ALTER TABLE exchange_event_buy_token RENAME TO exchange_event_token;

COMMIT;
