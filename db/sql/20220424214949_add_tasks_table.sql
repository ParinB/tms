-- +goose Up
CREATE EXTENSION IF NOT EXISTS btree_gist;

CREATE TABLE tasks(
    id SERIAL PRIMARY KEY,
    description varchar not null,
    start_time timestamp not null CONSTRAINT start_time_greater_than_now CHECK(start_time > now()::timestamp),
    end_time timestamp not null CONSTRAINT end_time_after_start_time CHECK(end_time > start_time),
    reminder_period timestamp not null CONSTRAINT reminder_period_before_start_time CHECK(reminder_period < start_time),
    user_id integer not null,
    created_at timestamp not null default now()::timestamp,
    updated_at timestamp not null default now()::timestamp,
    EXCLUDE USING gist (user_id WITH =, tsrange(start_time, end_time) WITH &&),

    CONSTRAINT fk_task_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE tasks;
