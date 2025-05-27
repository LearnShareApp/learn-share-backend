CREATE TABLE IF NOT EXISTS public.categories(
        category_id SERIAL PRIMARY KEY NOT NULL,
        name TEXT UNIQUE NOT NULL,
        min_age INTEGER NOT NULL
);

-- Seed initial categories
INSERT INTO public.categories (category_id, name, min_age)
VALUES
    (1, 'Cooking', 12),
    (2, 'Programming', 6),
    (3, 'Drawing', 3),
    (4, 'Dancing', 6),
    (5, 'English', 6),
    (6, 'Russian', 6),
    (7, 'Public Speaking', 14),
    (8, 'Physics', 6),
    (9, 'Biology', 6),
    (10,'History', 6),
    (11,'Maths', 6),
    (12,'Music', 6)
ON CONFLICT DO NOTHING;