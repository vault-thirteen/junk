create table public.boolean_index_table_b
(
	id serial
		constraint boolean_index_table_b_pk
			primary key,
	"Name" text not null,
	"IsSpecial" boolean default false not null
);

alter table public.boolean_index_table_b owner to test;

create index boolean_index_table_b_isspecial_index
	on public.boolean_index_table_b ("IsSpecial" desc);

