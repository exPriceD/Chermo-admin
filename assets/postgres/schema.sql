-- Таблица пользователей
CREATE TABLE IF NOT EXISTS users
(
    id         SERIAL PRIMARY KEY,                                                             -- Уникальный идентификатор пользователя
    username   VARCHAR(255) NOT NULL,                                                          -- Имя пользователя
    password   VARCHAR(255) NOT NULL,                                                          -- Пароль пользователя
    role       VARCHAR(255) NOT NULL,                                                          -- Роль пользователя (например, администратор, редактор)
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP AT TIME ZONE 'Europe/Moscow' -- Дата и время создания записи по МСК
);

-- Таблица мероприятий
CREATE TABLE IF NOT EXISTS events
(
    id          SERIAL PRIMARY KEY,                                                             -- Уникальный идентификатор мероприятия
    title       VARCHAR(255) NOT NULL,                                                          -- Название мероприятия
    description TEXT         NOT NULL,                                                          -- Описание мероприятия
    created_at  TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP AT TIME ZONE 'Europe/Moscow' -- Дата и время создания записи по МСК
);

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

-- Таблица для хранения расписания мероприятий
CREATE TABLE IF NOT EXISTS event_schedule
(
    id              SERIAL PRIMARY KEY,                                                             -- Уникальный идентификатор расписания мероприятия
    event_id        INTEGER REFERENCES events (id),                                                 -- Ссылка на мероприятие
    event_date      DATE    NOT NULL,                                                               -- Дата проведения мероприятия
    start_time      TIME    NOT NULL,                                                               -- Время начала мероприятия
    end_time        TIME    NOT NULL,                                                               -- Время окончания мероприятия
    total_slots     INTEGER NOT NULL,                                                               -- Общее количество мест
    available_slots INTEGER NOT NULL,                                                               -- Доступное количество мест
    created_at      TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP AT TIME ZONE 'Europe/Moscow' -- Дата и время создания записи по МСК
);

-- Таблица для хранения записей участников на определенные даты и время
CREATE TABLE IF NOT EXISTS event_registrations
(
    id                SERIAL PRIMARY KEY,                                                             -- Уникальный идентификатор записи
    schedule_id       INTEGER REFERENCES event_schedule (id),                                         -- Ссылка на расписание мероприятия
    visitor_id        INTEGER REFERENCES visitors (id),                                               -- Ссылка на посетителя
    registration_date TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP AT TIME ZONE 'Europe/Moscow' -- Дата и время регистрации по МСК
);

