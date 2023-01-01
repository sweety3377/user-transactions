CREATE SCHEMA IF NOT EXISTS clients
CREATE TABLE IF NOT EXISTS clients
(
    id         TEXT UNIQUE,
    name       TEXT      NOT NULL,
    surname    TEXT      NOT NULL,
    lastname   TEXT      NOT NULL,
    balance    BIGINT    NOT NULL,
    created_at TIMESTAMP NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_clients_id ON clients.clients(id);
