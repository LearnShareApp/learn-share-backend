CREATE TABLE IF NOT EXISTS public.categories(
        category_id SERIAL PRIMARY KEY NOT NULL,
        name TEXT UNIQUE NOT NULL,
        min_age INTEGER NOT NULL
);

-- Seed initial categories
INSERT INTO public.categories (name, min_age)
VALUES
    ('Cooking', 12),
    ('Programming', 6),
    ('Drawing', 3),
    ('Dancing', 6),
    ('English', 6),
    ('Russian', 6),
    ('Public Speaking', 14),
    ('Physics', 6),
    ('Biology', 6),
    ('History', 6),
    ('Maths', 6),
    ('Music', 6)
ON CONFLICT (name) DO NOTHING;