CREATE TABLE IF NOT EXISTS public.users(
       user_id SERIAL PRIMARY KEY NOT NULL,
       kratos_identity_id UUID UNIQUE, --NOT NULL
       email TEXT UNIQUE NOT NULL,
       name TEXT NOT NULL,
       surname TEXT NOT NULL,
       password TEXT NOT NULL,
       registration_date TIMESTAMPTZ DEFAULT NOW(),
       birthdate DATE NOT NULL,
       avatar TEXT NOT NULL DEFAULT ''
);