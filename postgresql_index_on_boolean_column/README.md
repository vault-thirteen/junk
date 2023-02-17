# postgresql_index_on_boolean_column
A performance test of an index on the boolean column in the PostgreSQL database.

## Why and what for ?
Taking part in technical interviews for one of the Golang Dev. positions, I got 
a question about why an index on a boolean column in PostgreSQL does not work 
faster than a simple query without indices. I told that it may occur if the 
table is very small and indices do not show the acceleration on small data. The 
"correct" answer was that in PostgreSQL indices on boolean columns do not show 
performance gain.   

"Hm..."

In this test you can see with your own eyes how an index on a boolean value 
column shows a significant increase in performance in PostgreSQL database.     

![](https://raw.githubusercontent.com/vault-thirteen/postgresql_index_on_boolean_column/main/test/Screenshot_2021-10-07_16-49-17.png)
