CREATE TABLE IF NOT EXISTS public.schedule_times (
        schedule_time_id SERIAL PRIMARY KEY,
        teacher_id INTEGER NOT NULL REFERENCES teachers(teacher_id) ON DELETE CASCADE,
        datetime TIMESTAMPTZ NOT NULL,
        is_available BOOLEAN NOT NULL DEFAULT TRUE,
        CONSTRAINT unique_teacher_schedule_time UNIQUE (teacher_id, datetime)
);