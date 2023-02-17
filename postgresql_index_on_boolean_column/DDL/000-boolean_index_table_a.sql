create table public.boolean_index_table_a
(
	id serial
		constraint boolean_index_table_a_pk
			primary key,
	"Name" text not null,
	"IsSpecial" boolean default false not null
);

alter table public.boolean_index_table_a owner to test;

