CREATE TABLE IF NOT EXISTS users
(
    id UUID PRIMARY KEY,
    login VARCHAR(50) UNIQUE,
    Pass_hash Varchar
);

INSERT INTO users(id, login, Pass_hash)
VALUES ('212755e4-0922-49aa-a46b-7912488220a8', 'admin', 'admin')
