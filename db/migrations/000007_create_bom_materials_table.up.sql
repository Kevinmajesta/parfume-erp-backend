BEGIN;

CREATE TABLE IF NOT EXISTS bom_materials (
    id_bommaterial VARCHAR(255) PRIMARY KEY NOT NULL,
    id_material VARCHAR(255) NOT NULL,
    id_bom VARCHAR(255) NOT NULL,
    materialname VARCHAR(50) NOT NULL, 
    quantity VARCHAR(50) NOT NULL,
    unit VARCHAR(50) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);

COMMIT;
