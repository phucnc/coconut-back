BEGIN;

ALTER TABLE collectible ADD COLUMN IF NOT EXISTS quote_token_id int NOT NULL DEFAULT 1 ;

CREATE TABLE IF NOT EXISTS token (
    id int GENERATED ALWAYS AS IDENTITY PRIMARY KEY ,
    name text NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now() ,
    updated_at timestamptz NOT NULL DEFAULT now() ,
    deleted_at timestamptz DEFAULT NULL
);

INSERT INTO token
    (name)
VALUES
    ('BNB'),
    ('BUSD'),
    ('CONT');

COMMIT;
