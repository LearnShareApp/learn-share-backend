CREATE TABLE IF NOT EXISTS public.statuses (
        status_id SERIAL PRIMARY KEY,
        name TEXT UNIQUE NOT NULL
);

-- Seed initial statuses
INSERT INTO public.statuses (name)
VALUES
    ('ongoing'),
    ('cancel'),
    ('verification'),
    ('waiting'),
    ('finished')
ON CONFLICT (name) DO NOTHING;