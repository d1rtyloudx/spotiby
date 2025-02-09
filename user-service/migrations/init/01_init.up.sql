CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS CITEXT;

CREATE TYPE role AS ENUM ('admin', 'user');

CREATE TABLE IF NOT EXISTS credentials
(
    id            UUID PRIMARY KEY                  DEFAULT uuid_generate_v4(),
    username      VARCHAR(32) UNIQUE       NOT NULL,
    email         VARCHAR(64) UNIQUE       NOT NULL,
    hash_pass     BYTEA                    NOT NULL,
    role          role                     NOT NULL DEFAULT 'user',
    created_at    TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMP WITH TIME ZONE          DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS profiles
(
    id             UUID PRIMARY KEY                  DEFAULT uuid_generate_v4(),
    display_name   VARCHAR(32)              NOT NULL DEFAULT '',
    first_name     VARCHAR(32)              NOT NULL DEFAULT '',
    last_name      VARCHAR(32)              NOT NULL DEFAULT '',
    description    VARCHAR(256)             NOT NULL DEFAULT '',
    credential_id  UUID REFERENCES credentials(id) ON DELETE CASCADE UNIQUE,
    avatar_url     VARCHAR(128)             NOT NULL DEFAULT ''
);

CREATE TABLE IF NOT EXISTS followers
(
    follower_id   UUID REFERENCES profiles (id) ON DELETE CASCADE,
    following_id  UUID REFERENCES profiles (id) ON DELETE CASCADE,
    created_at    TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    PRIMARY KEY (follower_id, following_id)
);

CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at := NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER before_update_credentials
BEFORE UPDATE ON credentials
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();