BEGIN;

CREATE TABLE IF NOT EXISTS rfqs (
    id_rfq VARCHAR(255) PRIMARY KEY NOT NULL,
    order_date TIMESTAMPTZ NOT NULL,
    id_vendor VARCHAR(255) NOT NULL,
    status VARCHAR(50) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);

COMMIT;
