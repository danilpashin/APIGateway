CREATE TABLE IF NOT EXISTS products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(150) NOT NULL,
    manufacturer VARCHAR(50) NOT NULL,
    price INT NOT NULL DEFAULT 100,
    amount INT NOT NULL DEFAULT 0,
    status BOOLEAN DEFAULT NULL,
    category VARCHAR(100) NOT NULL CHECK (category IN ('Household appliances', 'Smartphones and photographic equipment', 'TV, consoles, and audio', 'PCs, laptops, peripherals', 'PC accessories', 'Network equipment')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);