BEGIN;

CREATE TABLE IF NOT EXISTS rfqs_products (
    id_rfqproduct VARCHAR(255) PRIMARY KEY NOT NULL,
    id_vendor VARCHAR(255) NOT NULL,
    id_product VARCHAR(255) NOT NULL,
    id_rfq VARCHAR(255) NOT NULL,
    productname VARCHAR(255) NOT NULL, 
    quantity VARCHAR(50) NOT NULL,
    unitprice VARCHAR(50) NOT NULL,
    tax VARCHAR(50) NOT NULL,
    subtotal VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);

COMMIT;
