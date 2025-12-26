CREATE TYPE study_status AS ENUM ('NEW', 'QUEUED', 'PROCESSING', 'COMPLETED', 'FAILED');

create table study
(
	id                    uuid primary key,
	name                  text,
	status                study_status,
	owner_id              uuid references "user" (id),
	image_id              uuid,
	predicted_class_id    uuid references study_classes (id),
	predicted_class_score float,
	error_description     text null,
	created_at            timestamp,
	updated_at            timestamp
);
COMMENT ON TABLE study IS 'Информация об исследовании';
COMMENT ON COLUMN study.id IS 'UUID исследования';
COMMENT ON COLUMN study.name IS 'Название исследования';
COMMENT ON COLUMN study.owner_id IS 'Владелец исследования';
COMMENT ON COLUMN study.image_id IS 'UUID изображения предоставленного для исследования';
COMMENT ON COLUMN study.predicted_class_id IS 'Класс распознанный МЛ моделью';
COMMENT ON COLUMN study.predicted_class_score IS 'Точность';
COMMENT ON COLUMN study.status IS 'Тип статуса исследования(NEW, QUEUED, PROCESSING, COMPLETED, FAILED)';
COMMENT ON COLUMN study.error_description IS 'Описание ошибки, если статус - FAILED';
