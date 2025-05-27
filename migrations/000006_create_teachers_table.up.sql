CREATE TABLE IF NOT EXISTS public.teachers(
      teacher_id SERIAL PRIMARY KEY NOT NULL,
      user_id INTEGER UNIQUE NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
      rate FLOAT NOT NULL DEFAULT 0.0,
      total_rate_score INTEGER NOT NULL DEFAULT 0,
      reviews_count INTEGER NOT NULL DEFAULT 0
);