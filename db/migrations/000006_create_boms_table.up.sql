BEGIN;

CREATE TABLE IF NOT EXISTS boms (
    id_bom VARCHAR(255) PRIMARY KEY NOT NULL,
    id_product VARCHAR(255) NOT NULL,
    productname VARCHAR(50) NOT NULL, 
    productPreference VARCHAR(50) NOT NULL,
    quantity VARCHAR(50) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);

COMMIT;
