BEGIN;

ALTER TABLE collectible ADD COLUMN IF NOT EXISTS status int NOT NULL DEFAULT 0 ;

CREATE TABLE IF NOT EXISTS event (
    id int GENERATED ALWAYS AS IDENTITY PRIMARY KEY ,
    title text NOT NULL,
    banner text  DEFAULT NULL,
    content text NOT NULL,
    status int DEFAULT 0,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    deleted_at timestamptz  DEFAULT NULL
);

COMMIT;
