CREATE TYPE user_role AS ENUM ('admin', 'staff', 'customer');

CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    clerk_id TEXT NOT NULL UNIQUE, -- Clerk's user_xxx ID, the bridge between both systems
    email TEXT,
    first_name TEXT,
    last_name TEXT,
    role user_role NOT NULL DEFAULT 'customer', -- local copy, Clerk's session claims stay the source of truth
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
