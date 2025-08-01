CREATE TABLE IF NOT EXISTS public.statuses (
        status_id SERIAL PRIMARY KEY,
        name TEXT UNIQUE NOT NULL
);

-- Seed initial statuses
INSERT INTO public.statuses (status_id, name)
VALUES
    (1, 'ongoing'),
    (2, 'cancel'),
    (3, 'verification'),
    (4, 'waiting'),
    (5, 'finished')
ON CONFLICT DO NOTHING;