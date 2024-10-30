drop table if exists employees;
CREATE TABLE employees (
    id SERIAL PRIMARY KEY,
    fio varchar (100),
    salary integer,
    age integer
   );

INSERT INTO employees (fio, salary, age) VALUES ('Алексей Задонский', 50000, 61);
INSERT INTO employees (fio, salary, age) VALUES ('Ярослав Петров', 60000, 19);
INSERT INTO employees (fio, salary, age) VALUES ('Коновалова Таиссия', 25000, 25);
INSERT INTO employees (fio, salary, age) VALUES ('Перепелица Максим', 120000, 40);

select * from employees;