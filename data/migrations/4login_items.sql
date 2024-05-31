CREATE TABLE IF NOT EXISTS login_items (
    id UUID PRIMARY KEY,
    item_id UUID REFERENCES items (id) ON DELETE CASCADE,
    login VARCHAR(50),
    encrypt_password TEXT
)