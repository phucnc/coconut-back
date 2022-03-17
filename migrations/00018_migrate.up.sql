BEGIN;

ALTER TABLE collectible_report ADD COLUMN IF NOT EXISTS status int NOT NULL DEFAULT 0 ;

CREATE TABLE IF NOT EXISTS trend (
    id int GENERATED ALWAYS AS IDENTITY PRIMARY KEY ,
    collectible_id bigint NOT NULL,
    in_order int NOT NULL default 1,
    adsvertisement boolean NOT NULL default false,
    created_at timestamptz NOT NULL DEFAULT now()
);

ALTER TABLE trend ADD FOREIGN KEY (collectible_id) REFERENCES collectible(id);

COMMIT;
