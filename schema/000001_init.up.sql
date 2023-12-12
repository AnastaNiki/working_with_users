CREATE TABLE users
(
	id serial not null unique,
	firstName varchar(255) not null,
	lastName varchar(255) not null,
	patronymic varchar(255),
	login varchar(255) not null,
	archive boolean not null
);
