BEGIN;

CREATE TABLE IF NOT EXISTS schedules (
    id_schedules UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    qty_kolam VARCHAR(255) NOT NULL,
    date_schedules TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);

COMMIT;
