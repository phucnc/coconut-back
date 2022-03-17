BEGIN ;

CREATE TABLE IF NOT EXISTS collectible (
                                              id bigint GENERATED ALWAYS AS IDENTITY ,
                                              guid uuid DEFAULT uuid_generate_v4() NOT NULL ,
                                              title text NOT NULL ,
                                              description text NOT NULL DEFAULT '',
                                              upload_file text NOT NULL DEFAULT '' ,
                                              properties jsonb NOT NULL DEFAULT '{}'::jsonb ,
                                              royalties numeric(5,4) NOT NULL ,
                                              instant_sale_price numeric(128,8) NOT NULL DEFAULT 0,
                                              unlock_once_purchased bool NOT NULL DEFAULT FALSE ,
                                              token text DEFAULT NULL,
                                              token_id numeric(120, 0) DEFAULT NULL,
                                              created_at timestamptz NOT NULL DEFAULT now() ,
                                              updated_at timestamptz NOT NULL DEFAULT now() ,
                                              deleted_at timestamptz DEFAULT NULL
);



CREATE TABLE IF NOT EXISTS collectible_category (
                                                       collectible_id bigint NOT NULL ,
                                                       category_id bigint NOT NULL ,
                                                       created_at timestamptz NOT NULL DEFAULT now() ,
                                                       updated_at timestamptz NOT NULL DEFAULT now() ,
                                                       deleted_at timestamptz DEFAULT NULL
);







ALTER TABLE collectible ADD PRIMARY KEY (id);
ALTER TABLE collectible ADD CONSTRAINT collectible_guid_key UNIQUE (guid);
ALTER TABLE collectible ADD CONSTRAINT collectible_royalties_check CHECK ( royalties >= 0  AND royalties <= 1 );
ALTER TABLE collectible ADD CONSTRAINT collectible_instant_sale_price_check CHECK ( instant_sale_price >= 0 );
CREATE INDEX IF NOT EXISTS collectible_guid_idx ON collectible(guid);
CREATE INDEX IF NOT EXISTS collectible_created_at_idx ON collectible(created_at);
CREATE INDEX IF NOT EXISTS collectible_created_at_id_idx ON collectible(created_at, id);


ALTER TABLE collectible_category ADD PRIMARY KEY (collectible_id, category_id);
CREATE INDEX IF NOT EXISTS collectible_category_collectible_id_category_id_idx ON collectible_category(collectible_id, category_id);
ALTER TABLE collectible_category ADD FOREIGN KEY (collectible_id) REFERENCES collectible(id);
ALTER TABLE collectible_category ADD FOREIGN KEY (category_id) REFERENCES category(id);

COMMIT ;