CREATE TABLE IF NOT EXISTS public.reviews (
        review_id SERIAL PRIMARY KEY,
        teacher_id INTEGER NOT NULL REFERENCES teachers(teacher_id) ON DELETE CASCADE,
        student_id INTEGER NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
        category_id INTEGER NOT NULL REFERENCES categories(category_id) ON DELETE CASCADE,
        skill_id INTEGER NOT NULL REFERENCES skills(skill_id) ON DELETE CASCADE,
        rate SMALLINT NOT NULL,
        comment TEXT NOT NULL DEFAULT '',
        CONSTRAINT unique_review UNIQUE (teacher_id, student_id, category_id, skill_id)
);

-- Create function for updating ratings
CREATE OR REPLACE FUNCTION update_skill_and_teacher_on_review()
    RETURNS TRIGGER AS $$
BEGIN
    -- Update skills table
    UPDATE skills
    SET
        reviews_count = reviews_count + 1,
        total_rate_score = total_rate_score + NEW.rate,
        rate = (total_rate_score + NEW.rate)::decimal / (reviews_count + 1)
    WHERE skill_id = NEW.skill_id;

    -- Update teachers table
    UPDATE teachers
    SET
        reviews_count = reviews_count + 1,
        total_rate_score = total_rate_score + NEW.rate,
        rate = (total_rate_score + NEW.rate)::decimal / (reviews_count + 1)
    WHERE teacher_id = NEW.teacher_id;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create trigger for updating ratings
DO $$
    BEGIN
        IF NOT EXISTS (
            SELECT trigger_name
            FROM information_schema.triggers
            WHERE event_object_table = 'reviews'
              AND trigger_name = 'update_skill_and_teacher_on_review'
        ) THEN
            CREATE TRIGGER update_skill_and_teacher_on_review
                AFTER INSERT ON reviews
                FOR EACH ROW
            EXECUTE FUNCTION update_skill_and_teacher_on_review();
        END IF;
    END $$;