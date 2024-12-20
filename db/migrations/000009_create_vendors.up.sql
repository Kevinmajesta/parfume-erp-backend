BEGIN;

CREATE TABLE
    IF NOT EXISTS vendors (
        id_vendor VARCHAR(255) PRIMARY KEY NOT NULL,
        vendorname VARCHAR(255) NOT NULL,
        addressone VARCHAR(255) NOT NULL,
        addresstwo VARCHAR(255) NULL,
        city VARCHAR(255) NULL,
        state VARCHAR(255) NULL,
        zip VARCHAR(255) NULL,
        country VARCHAR(255) NULL,
        phone VARCHAR(255) NULL,
        email VARCHAR(255) NULL,
        website VARCHAR(255) NULL,
        status VARCHAR(255) NOT NULL,
        created_at TIMESTAMPTZ NOT NULL DEFAULT now (),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT now (),
        deleted_at TIMESTAMPTZ
    );

COMMIT;