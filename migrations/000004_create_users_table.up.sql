CREATE TABLE users (
    id            BIGSERIAL    PRIMARY KEY,
    name          TEXT         NOT NULL,
    email         TEXT         NOT NULL,
    password_hash TEXT         NOT NULL,
    role          TEXT         NOT NULL DEFAULT 'customer'
                  CHECK (role IN ('admin', 'customer')),
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ  NOT NULL DEFAULT now()
);
