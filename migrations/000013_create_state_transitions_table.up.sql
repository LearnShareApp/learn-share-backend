CREATE TABLE IF NOT EXISTS public.state_transitions (
    transition_id SERIAL PRIMARY KEY,
    state_machine_id INTEGER NOT NULL REFERENCES state_machines(state_machine_id),
    current_state_id INTEGER NOT NULL REFERENCES states(state_id),
    next_state_id INTEGER NOT NULL REFERENCES states(state_id),
    UNIQUE (state_machine_id, current_state_id, next_state_id)
);