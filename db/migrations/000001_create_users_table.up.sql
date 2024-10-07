BEGIN;

CREATE TABLE IF NOT EXISTS users (
    id_user UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    fullname VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    phone VARCHAR(255) NOT NULL,
    role VARCHAR(255) NOT NULL,
    status BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ,
    reset_code VARCHAR(255),
    reset_code_expires_at TIMESTAMPTZ,
    verification BOOLEAN DEFAULT false,
    verification_code VARCHAR(255)
);

COMMIT;
