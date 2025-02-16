CREATE TABLE purchases (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    item_name VARCHAR(255) NOT NULL,
    price INT NOT NULL CHECK (price > 0),
    created_at TIMESTAMP DEFAULT now()
);
