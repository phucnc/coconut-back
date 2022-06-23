/*banner*/
BEGIN ;

CREATE TABLE IF NOT EXISTS BANNER (
                                              id bigint GENERATED ALWAYS AS IDENTITY ,
                                              status int DEFAULT 0,
                                              name text NOT NULL ,
                                              picture text NOT NULL,
                                              link text DEFAULT NULL ,
                                              in_order int DEFAULT 0,
                                              admin int DEFAULT 0,  
                                              begin_date DATE  NULL ,  
                                              end_date DATE  NULL ,       
                                              created_at timestamptz NOT NULL DEFAULT now() ,
                                              updated_at timestamptz NOT NULL DEFAULT now() ,
                                              deleted_at timestamptz DEFAULT NULL
);

ALTER TABLE BANNER ADD PRIMARY KEY (id);

COMMIT;