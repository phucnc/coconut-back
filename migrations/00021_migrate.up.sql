BEGIN ;

CREATE TABLE IF NOT EXISTS notice (
                                              id bigint GENERATED ALWAYS AS IDENTITY ,
                                              status int DEFAULT 0,
                                              account_id bigint NOT NULL,
                                              content text NOT NULL ,
                                              created_at timestamptz NOT NULL DEFAULT now() ,
                                              updated_at timestamptz NOT NULL DEFAULT now() ,
                                              deleted_at timestamptz DEFAULT NULL
);




ALTER TABLE notice ADD PRIMARY KEY (id);
ALTER TABLE notice ADD FOREIGN KEY (account_id) REFERENCES account(id);

COMMIT;