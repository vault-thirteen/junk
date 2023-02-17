--{ Создание таблицы simple_events }--
CREATE TABLE simple_events
(
    id bigserial NOT NULL CONSTRAINT simple_events_pk PRIMARY KEY,
    user_id varchar(300) NOT NULL,
    file_id varchar(300) NOT NULL,
    event_type_id smallint NOT NULL CONSTRAINT simple_events_event_types_id_fk REFERENCES event_types,
    event_time timestamp with time zone DEFAULT now() NOT NULL
);

CREATE INDEX simple_events_user_id_index
    ON simple_events (user_id);

CREATE INDEX simple_events_file_id_index
    ON simple_events (file_id);

CREATE INDEX simple_events_event_time_index
    ON simple_events (event_time);
