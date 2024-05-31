CREATE TABLE IF NOT EXISTS items(
    id UUID PRIMARY KEY,
    type int,
    name VARCHAR(50),
    folder_id UUID REFERENCES Folders(id),
    user_id UUID REFERENCES Users(id),
    is_favorite BOOLEAN
)