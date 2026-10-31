CREATE TYPE role as ENUM( 'Admin', 'Staff', 'Customer' )

CREATE TABLE IF NOT EXISTS users(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid();
    email VARCHAR(255) UNIQUE NOT NULL;
    role
);