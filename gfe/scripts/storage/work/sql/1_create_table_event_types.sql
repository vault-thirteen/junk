--{ Создание таблицы event_types }--
CREATE TABLE event_types
(
    id smallint NOT NULL CONSTRAINT event_types_pk PRIMARY KEY,
    description_ru text NOT NULL UNIQUE,
    description_en text NOT NULL UNIQUE,
    is_simple boolean NOT NULL, is_aggregated boolean NOT NULL
);

INSERT INTO public.event_types (id, description_ru, description_en, is_simple, is_aggregated) VALUES (1, 'Создание', 'Creation', true, false);
INSERT INTO public.event_types (id, description_ru, description_en, is_simple, is_aggregated) VALUES (2, 'Загрузка', 'Upload', true, false);
INSERT INTO public.event_types (id, description_ru, description_en, is_simple, is_aggregated) VALUES (3, 'Скачивание', 'Download', false, true);
INSERT INTO public.event_types (id, description_ru, description_en, is_simple, is_aggregated) VALUES (4, 'Изменение', 'Modification', true, false);
