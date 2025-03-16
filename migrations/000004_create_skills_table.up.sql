CREATE TABLE IF NOT EXISTS public.skills (
        skill_id SERIAL PRIMARY KEY,
        teacher_id INTEGER NOT NULL REFERENCES teachers(teacher_id) ON DELETE CASCADE,
        category_id INTEGER NOT NULL REFERENCES categories(category_id) ON DELETE CASCADE,
        video_card_link TEXT,
        about TEXT,
        rate FLOAT NOT NULL DEFAULT 0,
        total_rate_score INTEGER NOT NULL DEFAULT 0,
        reviews_count INTEGER NOT NULL DEFAULT 0,
        is_active BOOLEAN NOT NULL DEFAULT TRUE,
        CONSTRAINT unique_teacher_category UNIQUE (teacher_id, category_id)
);