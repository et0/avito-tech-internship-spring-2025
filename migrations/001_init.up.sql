CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    role TEXT NOT NULL CHECK (role IN ('employee', 'moderator')),
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS pvz (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    registration_date TIMESTAMP NOT NULL DEFAULT NOW(),
    city TEXT NOT NULL CHECK (city IN ('Москва', 'Санкт-Петербург', 'Казань'))
);

CREATE TABLE IF NOT EXISTS receptions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    pvz_id UUID NOT NULL REFERENCES pvz(id),
    status TEXT NOT NULL CHECK (status IN ('in_progress', 'close'))
);

CREATE UNIQUE INDEX unique_active_reception ON receptions (pvz_id) WHERE (status = 'in_progress');

CREATE TABLE IF NOT EXISTS products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    type TEXT NOT NULL CHECK (type IN ('электроника', 'одежда', 'обувь')),
    reception_id UUID NOT NULL REFERENCES receptions(id)
);