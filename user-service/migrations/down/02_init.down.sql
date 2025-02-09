DROP TRIGGER IF EXISTS before_update_credentials ON credentials;

DROP TABLE IF EXISTS credentials CASCADE;
DROP TABLE IF EXISTS profiles;

DROP TYPE role;