BEGIN ;

CREATE TABLE IF NOT EXISTS comment (
                                              id bigint GENERATED ALWAYS AS IDENTITY ,
                                              content text NOT NULL ,
                                              collectible_id bigint NOT NULL,
                                              account_id bigint NOT NULL,
                                              created_at timestamptz NOT NULL DEFAULT now() ,
                                              updated_at timestamptz NOT NULL DEFAULT now() ,
                                              deleted_at timestamptz DEFAULT NULL
);

CREATE TABLE IF NOT EXISTS comment_like (        
                                                       comment_id bigint NOT NULL,
                                                       account_id bigint NOT NULL,
                                                       created_at timestamptz NOT NULL DEFAULT now() ,
                                                       updated_at timestamptz NOT NULL DEFAULT now() 
);


ALTER TABLE comment ADD PRIMARY KEY (id);
ALTER TABLE comment ADD FOREIGN KEY (collectible_id) REFERENCES collectible(id);
ALTER TABLE comment ADD FOREIGN KEY (account_id) REFERENCES account(id);

ALTER TABLE comment_like ADD PRIMARY KEY (comment_id, account_id);
ALTER TABLE comment_like ADD FOREIGN KEY (comment_id) REFERENCES comment(id);
ALTER TABLE comment_like ADD FOREIGN KEY (account_id) REFERENCES account(id);



CREATE TABLE IF NOT EXISTS collectible_report (
                                              id bigint GENERATED ALWAYS AS IDENTITY ,
                                              content text  DEFAULT NULL ,
                                              collectible_id bigint NOT NULL,
                                              account_id bigint NOT NULL,
                                              created_at timestamptz NOT NULL DEFAULT now() ,
                                              updated_at timestamptz NOT NULL DEFAULT now() ,
                                              deleted_at timestamptz DEFAULT NULL
);

ALTER TABLE collectible_report ADD PRIMARY KEY (id);
ALTER TABLE collectible_report ADD FOREIGN KEY (collectible_id) REFERENCES collectible(id);
ALTER TABLE collectible_report ADD FOREIGN KEY (account_id) REFERENCES account(id);


CREATE TABLE IF NOT EXISTS account_report (
                                              id bigint GENERATED ALWAYS AS IDENTITY ,
                                              content text NOT NULL ,
                                              account_report_id bigint NOT NULL,
                                              account_id bigint NOT NULL,
                                              created_at timestamptz NOT NULL DEFAULT now() ,
                                              updated_at timestamptz NOT NULL DEFAULT now() ,
                                              deleted_at timestamptz DEFAULT NULL
);

ALTER TABLE account_report ADD PRIMARY KEY (id);
ALTER TABLE account_report ADD FOREIGN KEY (account_report_id) REFERENCES account(id);
ALTER TABLE account_report ADD FOREIGN KEY (account_id) REFERENCES account(id);


CREATE TABLE IF NOT EXISTS collectible_like (
                                              id bigint GENERATED ALWAYS AS IDENTITY ,
                                              collectible_id bigint NOT NULL,
                                              account_id bigint NOT NULL,
                                              created_at timestamptz NOT NULL DEFAULT now() ,
                                              updated_at timestamptz NOT NULL DEFAULT now() 
);

ALTER TABLE collectible_like ADD PRIMARY KEY (collectible_id, account_id);
ALTER TABLE collectible_like ADD FOREIGN KEY (collectible_id) REFERENCES collectible(id);
ALTER TABLE collectible_like ADD FOREIGN KEY (account_id) REFERENCES account(id);


COMMIT ;