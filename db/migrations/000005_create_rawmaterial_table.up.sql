BEGIN;

CREATE TABLE IF NOT EXISTS rawmaterial (
    id_material VARCHAR(255) PRIMARY KEY NOT NULL,
    materialname VARCHAR(255) NOT NULL,
    materialcategory VARCHAR(255) NOT NULL,
    sellprice VARCHAR(50) NOT NULL,
    makeprice  VARCHAR(50) NOT NULL,
    unit VARCHAR(50) NOT NULL,
    image VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);

COMMIT;