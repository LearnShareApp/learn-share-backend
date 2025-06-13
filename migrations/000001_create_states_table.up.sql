CREATE TABLE IF NOT EXISTS public.states (
    state_id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);

INSERT INTO public.states (state_id, name)
VALUES
    (1, 'pending'),  -- skill + lesson (old verification)
    (2, 'approved'), -- skill
    (3, 'rejected'), -- skill + lesson (old cancel as opt   )


    (4, 'planned'), -- lesson (old waiting)
    (5, 'cancel'), -- lesson mb skill
    (6, 'ongoing'), -- lesson
    (7, 'finished'), -- lesson
    (8, 'conflicted'), -- lesson
    (9, 'completed') -- lesson

ON CONFLICT DO NOTHING;