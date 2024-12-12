BEGIN;

CREATE TABLE IF NOT EXISTS quotations_products (
    id_quotationsproduct VARCHAR(255) PRIMARY KEY NOT NULL,
    id_costumer VARCHAR(255) NOT NULL,
    id_product VARCHAR(255) NOT NULL,
    id_quotation VARCHAR(255) NOT NULL,
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
