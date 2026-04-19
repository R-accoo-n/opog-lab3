CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS products (
                                        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL,
    category_name TEXT NOT NULL,
    category_tax DOUBLE PRECISION NOT NULL CHECK (category_tax >= 0),
    price DOUBLE PRECISION NOT NULL CHECK (price >= 0)
    );