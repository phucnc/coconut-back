BEGIN ;

CREATE TABLE IF NOT EXISTS KYC (
                                              id bigint GENERATED ALWAYS AS IDENTITY ,
                                              status int DEFAULT 0,
                                              account_id bigint NOT NULL,
                                              fullname text NOT NULL ,
                                              birthday text NOT NULL,
                                              email text NOT NULL ,
                                              city text NOT NULL ,
                                              country text NOT NULL ,
                                              front_id text NOT NULL,
                                              back_id text not NULL,
                                              selfie_note  text not NULL,

                                              created_at timestamptz NOT NULL DEFAULT now() ,
                                              updated_at timestamptz NOT NULL DEFAULT now() ,
                                              deleted_at timestamptz DEFAULT NULL
);

ALTER TABLE KYC ADD PRIMARY KEY (id);
ALTER TABLE KYC ADD FOREIGN KEY (account_id) REFERENCES account(id);

ALTER TABLE account ADD COLUMN IF NOT EXISTS kyc int DEFAULT 0;

COMMIT;