CREATE TABLE IF NOT EXISTS public.state_machines (
    state_machine_id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    start_state_id INTEGER NOT NULL REFERENCES states(state_id)
);

INSERT INTO public.state_machines (state_machine_id, name, start_state_id)
VALUES
    (1, 'lesson', 1),
    (2, 'skill', 1)
ON CONFLICT DO NOTHING;