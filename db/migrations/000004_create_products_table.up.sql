BEGIN;

CREATE TABLE IF NOT EXISTS products (
    id_product VARCHAR(255) PRIMARY KEY NOT NULL,
    productname VARCHAR(255) NOT NULL,
    productcategory VARCHAR(255) NOT NULL,
    sellprice VARCHAR(50) NOT NULL,
    makeprice  VARCHAR(50) NOT NULL,
    pajak VARCHAR(50) NOT NULL,
    image VARCHAR(255) NOT NULL,
    qty VARCHAR(255) null,
    variant VARCHAR(3) NOT NULL CHECK (variant IN ('yes', 'no')),
    description TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);

COMMIT;