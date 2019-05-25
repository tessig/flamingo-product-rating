create table ratings
(
	id int auto_increment,
	name varchar(100) not null,
	title varchar(255) not null,
	text TEXT not null,
	created_at datetime not null,
	constraint ratings_pk
		primary key (id)
);

create index ratings_created_at_index
	on ratings (created_at);

create index ratings_name_index
	on ratings (name);

