CREATE TABLE IF NOT EXISTS public.state_transitions (
    transition_id SERIAL PRIMARY KEY,
    state_machine_id INTEGER NOT NULL REFERENCES state_machines(state_machine_id),
    current_state_id INTEGER NOT NULL REFERENCES states(state_id),
    next_state_id INTEGER NOT NULL REFERENCES states(state_id),
    UNIQUE (state_machine_id, current_state_id, next_state_id)
);

INSERT INTO public.state_transitions (transition_id, state_machine_id, current_state_id, next_state_id)
VALUES
    --lesson

    (1, 1, 1, 3),
    (2, 1, 1, 4),
    (3, 1, 4, 5),
    (4, 1, 4, 6),
    (5, 1, 6, 5),
    (6, 1, 6, 7),
    (7, 1, 7, 8),
    (8, 1, 7, 9),
    (9, 1, 8, 9),

    --skill
    (10, 2, 1, 2),
    (11, 2, 1, 3),
    (12, 2, 2, 1),
    (13, 2, 3, 1)

ON CONFLICT DO NOTHING;