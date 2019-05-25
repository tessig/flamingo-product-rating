alter table ratings
	add product_id int not null;

create index ratings_product_id_index
	on ratings (product_id);

