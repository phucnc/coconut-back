BEGIN;

CREATE TABLE IF NOT EXISTS account (
    id bigint GENERATED ALWAYS AS IDENTITY ,
    username text  DEFAULT NULL UNIQUE,
    avatar text DEFAULT NULL,
    cover text DEFAULT NULL,
    address text NOT NULL UNIQUE,
    info text DEFAULT NULL,
    twitter text DEFAULT NULL,
    facebook text DEFAULT NULL,
    tiktok text DEFAULT NULL,
    instagram text DEFAULT NULL,
    status  int NOT NULL DEFAULT 0,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    deleted_at timestamptz DEFAULT NULL
);

ALTER TABLE account ADD PRIMARY KEY (id);
CREATE INDEX IF NOT EXISTS account_username ON account (username);
CREATE INDEX IF NOT EXISTS account_address ON account (address);
ALTER TABLE account ADD CONSTRAINT account_username_check CHECK ( username = lower(username) );

COMMIT;
