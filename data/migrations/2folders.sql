CREATE TABLE IF NOT EXISTS folders
(
    id      UUID PRIMARY KEY,
    user_id UUID REFERENCES Users (id),
    name    VARCHAR(50)
);
INSERT INTO folders(id, user_id, name)
VALUES ('212754e4-0922-49aa-a46b-7912488220a8', '212755e4-0922-49aa-a46b-7912488220a8', 'name')
