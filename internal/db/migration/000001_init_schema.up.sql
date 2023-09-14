CREATE TABLE IF NOT EXISTS fio_data (
    id SERIAL PRIMARY KEY,
    "name" VARCHAR(255),
    surname VARCHAR(255),
    patronymic VARCHAR(255),
    age INT,
    gender VARCHAR(10),
    nationality VARCHAR(255)
);