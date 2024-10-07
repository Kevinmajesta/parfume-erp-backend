BEGIN;

CREATE TABLE IF NOT EXISTS suggestions (
    id_suggestion UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id_user) ON DELETE CASCADE,
    type TEXT,
    message TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);

COMMIT;
