-- Таблица музеев
CREATE TABLE IF NOT EXISTS museums
(
    id   SERIAL PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL
);

-- Таблица пользователей
CREATE TABLE IF NOT EXISTS users
(
    id         SERIAL PRIMARY KEY,
    username   VARCHAR(255) UNIQUE NOT NULL,
    password   VARCHAR(255)        NOT NULL,
    museum_id  INTEGER REFERENCES museums (id),
    role       VARCHAR(255)        NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Таблица мероприятий
CREATE TABLE IF NOT EXISTS events
(
    id          SERIAL PRIMARY KEY,
    title       VARCHAR(255) NOT NULL,
    description TEXT         NOT NULL,
    museum_id   INTEGER REFERENCES museums (id),
    image_url   TEXT         NOT NULL,
    created_at  TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Таблица посетителей
CREATE TABLE IF NOT EXISTS visitors
(
    id         SERIAL PRIMARY KEY,
    first_name VARCHAR(255) NOT NULL,
    last_name  VARCHAR(255) NOT NULL,
    patronymic VARCHAR(255) NOT NULL,
    phone      VARCHAR(255) NOT NULL,
    email      VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT unique_email UNIQUE (email)
);

-- Таблица для хранения расписания мероприятий
CREATE TABLE IF NOT EXISTS event_schedule
(
    id         SERIAL PRIMARY KEY,
    event_id   INTEGER REFERENCES events (id),
    event_date DATE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Таблица для хранения временных слотов мероприятий
CREATE TABLE IF NOT EXISTS event_timeslots
(
    id              SERIAL PRIMARY KEY,
    schedule_id     INTEGER REFERENCES event_schedule (id),
    start_time      TIME    NOT NULL,
    end_time        TIME    NOT NULL,
    total_slots     INTEGER NOT NULL,
    available_slots INTEGER NOT NULL,
    created_at      TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Таблица для хранения записей участников на определенные даты и время
CREATE TABLE IF NOT EXISTS event_registrations
(
    id                SERIAL PRIMARY KEY,
    timeslot_id       INTEGER REFERENCES event_timeslots (id),
    visitor_id        INTEGER REFERENCES visitors (id),
    status            VARCHAR(64)              DEFAULT 'pending' NOT NULL,
    registration_date TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT unique_timeslot_visitor UNIQUE (timeslot_id, visitor_id)
);

-- Индексы для таблицы музеев
CREATE INDEX IF NOT EXISTS idx_museum_id ON museums (id);

-- Индексы для таблицы пользователей
CREATE INDEX IF NOT EXISTS idx_users_username ON users (username);
CREATE INDEX IF NOT EXISTS idx_users_role ON users (role);

-- Индексы для таблицы мероприятий
CREATE INDEX IF NOT EXISTS idx_events_id ON events (id);
CREATE INDEX IF NOT EXISTS idx_events_title ON events (title);

-- Индексы для таблицы посетителей
CREATE INDEX IF NOT EXISTS idx_visitors_email ON visitors (email);
CREATE INDEX IF NOT EXISTS idx_visitors_phone ON visitors (phone);

-- Индексы для таблицы расписания мероприятий
CREATE INDEX IF NOT EXISTS idx_event_schedule_event_id ON event_schedule (event_id);
CREATE INDEX IF NOT EXISTS idx_event_schedule_event_date ON event_schedule (event_date);

-- Индексы для таблицы временных слотов мероприятий
CREATE INDEX IF NOT EXISTS idx_event_timeslots_schedule_id ON event_timeslots (schedule_id);
CREATE INDEX IF NOT EXISTS idx_event_timeslots_start_time ON event_timeslots (start_time);

-- Индексы для таблицы записей участников
CREATE INDEX IF NOT EXISTS idx_event_registrations_timeslot_id ON event_registrations (timeslot_id);
CREATE INDEX IF NOT EXISTS idx_event_registrations_visitor_id ON event_registrations (visitor_id);