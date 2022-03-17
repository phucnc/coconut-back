BEGIN;

CREATE TABLE IF NOT EXISTS category (
    id int GENERATED ALWAYS AS IDENTITY PRIMARY KEY ,
    name text NOT NULL,
    ext_id int GENERATED ALWAYS AS IDENTITY,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    deleted_at timestamptz DEFAULT NULL
);

ALTER TABLE category ADD CONSTRAINT  category_name_key UNIQUE (name);
ALTER TABLE category ADD CONSTRAINT  category_name_check CHECK ( name = lower(name) );
CREATE INDEX IF NOT EXISTS category_name_idx ON category(name);

INSERT INTO category (
    name
)
VALUES
    ('foods'),
    ('animal'),
    ('dance/sing'),
    ('funny'),
    ('satisfying'),
    ('beauty'),
    ('fashion'),
    ('tricks/skills'),
    ('sports'),
    ('outdoor activities'),
    ('daily-life'),
    ('animation/arts'),
    ('transportation'),
    ('science'),
    ('education'),
    ('travel'),
    ('health'),
    ('gaming')

ON CONFLICT ON CONSTRAINT category_name_key DO NOTHING;

CREATE INDEX IF NOT EXISTS category_name_idx ON category(name);

COMMIT;
