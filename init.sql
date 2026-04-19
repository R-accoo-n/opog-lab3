CREATE TABLE IF NOT EXISTS travellers (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    first_name TEXT NOT NULL,
    last_name  TEXT NOT NULL,
    age        INT  NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS uq_travellers_name_age
    ON travellers (first_name, last_name, age);
