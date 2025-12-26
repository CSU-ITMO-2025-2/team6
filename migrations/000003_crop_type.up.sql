create table crop_type
(
	id         uuid primary key,
	name       text,
	created_at timestamp,
	updated_at timestamp
);
COMMENT ON TABLE crop_type IS 'Тип сельхоз культуры';
COMMENT ON COLUMN crop_type.id IS 'UUID сельхоз культуры';
COMMENT ON COLUMN crop_type.name IS 'Название сельхоз культуры';
