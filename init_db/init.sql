CREATE TABLE
    IF NOT EXISTS items (
        id SERIAL PRIMARY KEY,
        url TEXT NOT NULL,
        title TEXT,
        image_url TEXT,
        current_price NUMERIC(10, 2) DEFAULT 0.00,
        status VARCHAR(50) DEFAULT 'pending',
        source VARCHAR(50) DEFAULT 'unknown',
        created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
    );

CREATE TABLE
    IF NOT EXISTS price_history (
        id SERIAL PRIMARY KEY,
        item_id INT NOT NULL REFERENCES items (id) ON DELETE CASCADE,
        price NUMERIC(10, 2) NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );