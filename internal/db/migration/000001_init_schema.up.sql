CREATE TABLE IF NOT EXISTS people (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    name TEXT,
    surname TEXT,
    patronymic TEXT,
    age INT,
    gender TEXT,
    nationality TEXT
);