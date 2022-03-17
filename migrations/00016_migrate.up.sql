BEGIN ;


CREATE TABLE IF NOT EXISTS report_type (
    id int GENERATED ALWAYS AS IDENTITY PRIMARY KEY ,
    description text NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);


INSERT INTO report_type (
    description
)
VALUES
    ('Sexual Content'),
    ('Violent or repulsive content'),
    ('Hateful or abusive content'),
    ('Hamful or dangerous content'),
    ('Spam of misleading'),
    ('Copyright');


ALTER TABLE collectible_report ADD COLUMN IF NOT EXISTS report_type_id int  NOT NULL;
ALTER TABLE collectible_report ADD FOREIGN KEY (report_type_id) REFERENCES report_type(id);

ALTER TABLE collectible ADD COLUMN IF NOT EXISTS view int  NOT NULL DEFAULT 0;
ALTER TABLE collectible ADD COLUMN IF NOT EXISTS creator text DEFAULT NULL; 


COMMIT;