CREATE TABLE IF NOT EXISTS public.lessons (
        lesson_id SERIAL PRIMARY KEY,
        student_id INTEGER NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
        teacher_id INTEGER NOT NULL REFERENCES teachers(teacher_id) ON DELETE CASCADE,
        category_id INTEGER NOT NULL REFERENCES categories(category_id) ON DELETE CASCADE,
        schedule_time_id INTEGER UNIQUE NOT NULL REFERENCES schedule_times(schedule_time_id) ON DELETE CASCADE,
        price INTEGER NOT NULL DEFAULT 0,
        status_id INTEGER DEFAULT NULL REFERENCES statuses(status_id) ON DELETE CASCADE,
        state_machine_item_id INTEGER NOT NULL REFERENCES state_machines_items(item_id)
);

-- Create function for default lesson status
CREATE OR REPLACE FUNCTION set_default_lesson_status()
    RETURNS TRIGGER AS $$
BEGIN
    IF NEW.status_id IS NULL THEN
        NEW.status_id := (
            SELECT status_id
            FROM statuses
            WHERE name = 'verification'
            LIMIT 1
        );
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create trigger for default lesson status
DO $$
    BEGIN
        IF NOT EXISTS (
            SELECT trigger_name
            FROM information_schema.triggers
            WHERE event_object_table = 'lessons'
              AND trigger_name = 'set_lesson_status'
        ) THEN
            CREATE TRIGGER set_lesson_status
                BEFORE INSERT ON lessons
                FOR EACH ROW
            EXECUTE FUNCTION set_default_lesson_status();
        END IF;
    END $$;