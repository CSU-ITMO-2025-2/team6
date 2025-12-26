create table study_classes
(
	id           uuid primary key,
	name         text,
	crop_type_id uuid,
	created_at   timestamp,
	updated_at   timestamp
);
COMMENT ON TABLE study_classes IS 'Классы сельхоз культур представленные в МЛ модели';
COMMENT ON COLUMN study_classes.id IS 'Название класса в МЛ модели';
COMMENT ON COLUMN study_classes.crop_type_id IS 'UUID сельхоз культуры, которую этот класс представляет';
