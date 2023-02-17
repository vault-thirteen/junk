--{ Создание таблицы download_events }--
CREATE TABLE download_events
(
    id bigserial NOT NULL CONSTRAINT download_events_pk PRIMARY KEY,
    user_id varchar(300) NOT NULL,
    file_id varchar(300) NOT NULL,
    event_time timestamp with time zone DEFAULT now() NOT NULL
);

CREATE INDEX download_events_user_id_index
    ON download_events (user_id);

CREATE INDEX download_events_file_id_index
    ON download_events (file_id);

CREATE INDEX download_events_event_time_index
    ON download_events (event_time);
