CREATE TABLE "user"
(
	id                 uuid primary key,
	username           TEXT      NOT NULL UNIQUE,
	pass_hash          bytea     NOT NULL,
	lastname           TEXT,
	name               TEXT,
	patronymic         TEXT,
	company_id         uuid      NULL,
	is_company_admin   boolean default FALSE,
	photo_id           uuid,
	birth_date         date,
	email              TEXT UNIQUE,
	is_email_confirmed boolean,
	created_at         timestamp NOT NULL,
	modified_at        timestamp NULL,
	deleted_at         timestamp NULL
);
COMMENT ON TABLE "user" IS 'Пользователи';
COMMENT ON COLUMN "user".id IS 'UUID пользователя';
COMMENT ON COLUMN "user".username IS 'Логин пользователя';
COMMENT ON COLUMN "user".pass_hash IS 'Хеш пароля пользователя';
COMMENT ON COLUMN "user".lastname IS 'Фамилия';
COMMENT ON COLUMN "user".name IS 'Имя';
COMMENT ON COLUMN "user".patronymic IS 'Отчество';
COMMENT ON COLUMN "user".company_id IS 'Компания пользователя';
COMMENT ON COLUMN "user".is_company_admin IS 'Является ли он админом своей компании';
COMMENT ON COLUMN "user".photo_id IS 'UUID аватарки пользователя';
COMMENT ON COLUMN "user".birth_date IS 'Дата рождения';
COMMENT ON COLUMN "user".email IS 'Почта';
COMMENT ON COLUMN "user".is_email_confirmed IS 'Подтверждена ли почта';

