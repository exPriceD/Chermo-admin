-- Таблица пользователей
CREATE TABLE IF NOT EXISTS users
(
    id         SERIAL PRIMARY KEY,                                                             -- Уникальный идентификатор пользователя
    username   VARCHAR(255) UNIQUE NOT NULL,                                                   -- Имя пользователя
    password   VARCHAR(255)        NOT NULL,                                                   -- Пароль пользователя
    role       VARCHAR(255)        NOT NULL,                                                   -- Роль пользователя (например, администратор, редактор)
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP AT TIME ZONE 'Europe/Moscow' -- Дата и время создания записи по МСК
);

-- Индексы для таблицы пользователей
CREATE INDEX idx_users_username ON users (username);
CREATE INDEX idx_users_role ON users (role);

-- Таблица мероприятий
CREATE TABLE IF NOT EXISTS events
(
    id          SERIAL PRIMARY KEY,                                                             -- Уникальный идентификатор мероприятия
    title       VARCHAR(255) NOT NULL,                                                          -- Название мероприятия
    description TEXT         NOT NULL,                                                          -- Описание мероприятия
    created_at  TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP AT TIME ZONE 'Europe/Moscow' -- Дата и время создания записи по МСК
);

-- Индексы для таблицы мероприятий
CREATE INDEX idx_events_id ON events (id);
CREATE INDEX idx_events_title ON events (title);

-- Таблица посетителей
CREATE TABLE IF NOT EXISTS visitors
(
    id         SERIAL PRIMARY KEY,                                                             -- Уникальный идентификатор посетителя
    first_name VARCHAR(255) NOT NULL,                                                          -- Имя посетителя
    last_name  VARCHAR(255) NOT NULL,                                                          -- Фамилия посетителя
    patronymic VARCHAR(255) NOT NULL,                                                          -- Отчество посетителя
    phone      VARCHAR(255) NOT NULL,                                                          -- Телефонный номер посетителя
    email      VARCHAR(255) NOT NULL,                                                          -- Электронная почта посетителя
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP AT TIME ZONE 'Europe/Moscow' -- Дата и время создания записи по МСК
);

-- Индексы для таблицы посетителей
CREATE INDEX idx_visitors_email ON visitors (email);
CREATE INDEX idx_visitors_phone ON visitors (phone);

-- Таблица для хранения расписания мероприятий
CREATE TABLE IF NOT EXISTS event_schedule
(
    id              SERIAL PRIMARY KEY,                                                             -- Уникальный идентификатор расписания мероприятия
    event_id        INTEGER REFERENCES events (id),                                                 -- Ссылка на мероприятие
    event_date      DATE    NOT NULL,                                                               -- Дата проведения мероприятия
    created_at      TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP AT TIME ZONE 'Europe/Moscow' -- Дата и время создания записи по МСК
);

-- Индексы для таблицы расписания мероприятий
CREATE INDEX idx_event_schedule_event_id ON event_schedule (event_id);
CREATE INDEX idx_event_schedule_event_date ON event_schedule (event_date);

-- Таблица для хранения временных слотов мероприятий
CREATE TABLE IF NOT EXISTS event_timeslots
(
    id              SERIAL PRIMARY KEY,                                                             -- Уникальный идентификатор временного слота
    schedule_id     INTEGER REFERENCES event_schedule (id),                                         -- Ссылка на расписание мероприятия
    start_time      TIME    NOT NULL,                                                               -- Время начала временного слота
    end_time        TIME    NOT NULL,                                                               -- Время окончания временного слота
    total_slots     INTEGER NOT NULL,                                                               -- Общее количество мест в этом временном слоте
    available_slots INTEGER NOT NULL,                                                               -- Доступное количество мест в этом временном слоте
    created_at      TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP AT TIME ZONE 'Europe/Moscow' -- Дата и время создания записи по МСК
);

-- Индексы для таблицы временных слотов мероприятий
CREATE INDEX idx_event_timeslots_schedule_id ON event_timeslots (schedule_id);
CREATE INDEX idx_event_timeslots_start_time ON event_timeslots (start_time);

-- Таблица для хранения записей участников на определенные даты и время
CREATE TABLE IF NOT EXISTS event_registrations
(
    id                SERIAL PRIMARY KEY,                                                             -- Уникальный идентификатор записи
    timeslot_id       INTEGER REFERENCES event_timeslots (id),                                        -- Ссылка на временной слот мероприятия
    visitor_id        INTEGER REFERENCES visitors (id),                                               -- Ссылка на посетителя
    registration_date TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP AT TIME ZONE 'Europe/Moscow' -- Дата и время регистрации по МСК
);

-- Индексы для таблицы записей участников
CREATE INDEX idx_event_registrations_timeslot_id ON event_registrations (timeslot_id);
CREATE INDEX idx_event_registrations_visitor_id ON event_registrations (visitor_id);