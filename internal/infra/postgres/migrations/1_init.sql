-- +goose Up

-- CREATE TYPE withdraw_status AS ENUM (
--     'NEW',
--     'REJECTED',
--     'PROCESSING',
--     'PROCESSED'
-- );

CREATE TYPE order_status AS ENUM (
    'NEW',
    'INVALID',
    'PROCESSING',
    'PROCESSED'
);

CREATE TABLE "user" (
    id UUID NOT NULL PRIMARY KEY,
    login character varying(255) NOT NULL UNIQUE,
    password char(255) NOT NULL
);

CREATE TABLE "account" (
    user_id UUID NOT NULL REFERENCES "user" ("id") UNIQUE,
    balance numeric(20, 2) NOT NULL DEFAULT 0
);

CREATE TABLE "order" (
    id UUID NOT NULL PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES "user" ("id"),
    order_number varchar NOT NULL UNIQUE,
    amount numeric(20, 2) NOT NULL,
    status order_status NOT NULL DEFAULT 'NEW',
    created_at timestamp without time zone DEFAULT null,
    updated_at timestamp without time zone DEFAULT null
);

CREATE TABLE "withdraw" (
    id UUID NOT NULL PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES "user" ("id"),
    order_number varchar NOT NULL,
    amount numeric(20, 2) NOT NULL,
--     status withdraw_status NOT NULL DEFAULT 'NEW',
    created_at timestamp without time zone DEFAULT null,
    updated_at timestamp without time zone DEFAULT null
);

-- +goose Down
DROP TABLE "withdraw";
DROP TABLE "order";
DROP TABLE "account";
DROP TABLE "user";

DROP TYPE "order_status";
-- DROP TYPE "withdraw_status";
