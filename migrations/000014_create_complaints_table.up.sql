CREATE TABLE IF NOT EXISTS public.complaints (
        complaint_id SERIAL PRIMARY KEY,
        complainer_id INTEGER NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
        reported_id INTEGER NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
        reason TEXT NOT NULL,
        description TEXT NOT NULL,
        created_at TIMESTAMPTZ DEFAULT NOW()
        -- CONSTRAINT unique_complaint UNIQUE (complainer_id, reported_id)
);