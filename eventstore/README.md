
Support this data schema

```
pflow=# select * from events;
-[ RECORD 1 ]----------------------------------
id       | 3e1e0cac-5402-4073-83c7-f226be62746e
schema   | octoe
action   | ON
multiple | 1
payload  |
state    | {1,1,1,1,1,1,1,1,1,1,0,0,0,1}
ts       | 2019-06-02 20:01:29.087686
uuid     | 35486c8c-1ee4-4b0e-8a53-b298d851f898
parent   | 00000000-0000-0000-0000-000000000000
-[ RECORD 2 ]----------------------------------
id       | 3e1e0cac-5402-4073-83c7-f226be62746e
schema   | octoe
action   | EXEC
multiple | 1
payload  |
state    | {1,1,1,1,1,1,1,1,1,0,1,0,0,1}
ts       | 2019-06-02 20:01:29.090497
uuid     | 71c5b788-2ec7-414c-ac0a-b79378af0dd6
parent   | 35486c8c-1ee4-4b0e-8a53-b298d851f898
```


```
pflow=# select * from states;
-[ RECORD 1 ]---------------------------------
id      | 3e1e0cac-5402-4073-83c7-f226be62746e
schema  | octoe
state   | {1,1,1,1,1,1,1,1,1,0,1,0,0,1}
head    | 71c5b788-2ec7-414c-ac0a-b79378af0dd6
created | 2019-06-02 20:01:29.087686
updated | 2019-06-02 20:01:29.087686
```
