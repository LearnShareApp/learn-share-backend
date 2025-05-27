CREATE TABLE IF NOT EXISTS public.state_machines_items (
    item_id SERIAL PRIMARY KEY,
    state_machine_id INTEGER NOT NULL REFERENCES state_machines(state_machine_id),
    state_id INTEGER NOT NULL REFERENCES states(state_id)
);