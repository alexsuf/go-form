package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

type Employee struct {
	ID     int
	FIO    string
	Salary int
	Age    int
}

var db *sql.DB

var pageTemplate = `
<!DOCTYPE html>
<html>
<head>
	<title>Сотрудники</title>
	<style>
		body {
			background-color: black;
			color: white;
			font-family: Arial, sans-serif;
		}
		h2 {
			text-align: center;
		}
		table {
			width: 95%;
			margin: 0 auto;
			border-collapse: collapse;
		}
		th, td {
			border: 1px solid white;
			padding: 5px;
			text-align: center;
		}
		th {
			background-color: #444;
		}
		input[type="text"], input[type="number"], input[type="date"] {
			color: black;
			font-size: 16px; /* Измените на нужный вам размер */
		}
		.button {
			background-color: #444;
			color: white;
			border: none;
			padding: 10px 20px;
			text-align: center;
			text-decoration: none;
			margin: 4px 2px;
			cursor: pointer;
			border-radius: 5px;
		}
		.center {
			text-align: center;
		}
		}
	</style>
</head>
<body>
	<h2>Список сотрудников компании "Рога и копыта"</h2>

	<div class="center">
		<form method="POST" action="/add">
			<input type="text" name="Fio" placeholder="ФИО" required style="font-size: 16px;">
			<input type="number" name="Salary" placeholder="Зарплата" required style="font-size: 16px;">
			<input type="number" name="Age" placeholder="Возраст" required style="font-size: 16px;">
			<input type="submit" class="Button" value="Добавить запись" style="background-color: blue; color: white; font-size: 16px;">
		</form>
	</div>

	<table>
		<tr>
			<th>ID</th>
			<th>ФИО</th>
			<th>Зарплата</th>
			<th>Возраст</th>
			<th>Действия</th>
		</tr>
		{{ range . }}
			<tr>
				<td>{{ .ID }}</td>
				<td>
					<form method="POST" action="/edit/{{ .ID }}">
						<input type="text" name="fio" value="{{ .FIO }}" required style="font-size: 16px;">
				</td>
				<td>
					<input type="number" name="salary" value="{{ .Salary }}" required>
				</td>
				<td>
					<input type="number" name="age" value="{{ .Age }}" required>
				</td>
				<td>
					<input type="submit" class="button" value="Сохранить" style="background-color: green; color: white; font-size: 16px;">
					</form>
					<form method="POST" action="/delete/{{ .ID }}" style="display:inline;">
						<input type="submit" class="button" value="Удалить"  style="background-color: red; color: white; font-size: 16px;" onclick="return confirm('Вы уверены, что хотите удалить эту запись?');">
					</form>
				</td>
			</tr>
		{{ end }}
	</table>
</body>
</html>
`

func main() {
	var err error
	connStr := "user=postgres password=secret host=postgres port=5432 dbname=postgres sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	http.HandleFunc("/", employeeHandler)
	http.HandleFunc("/edit/", editEmployeeHandler)
	http.HandleFunc("/delete/", deleteEmployeeHandler)
	http.HandleFunc("/add", addEmployeeHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func employeeHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, fio, salary, age FROM employees order by id")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var employees []Employee
	for rows.Next() {
		var emp Employee
		if err := rows.Scan(&emp.ID, &emp.FIO, &emp.Salary, &emp.Age); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		employees = append(employees, emp)
	}

	tmpl, err := template.New("page").Parse(pageTemplate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, employees)
}

func editEmployeeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		id := r.URL.Path[len("/edit/"):]

		fio := r.FormValue("fio")
		salary := r.FormValue("salary")
		age := r.FormValue("age")

		_, err := db.Exec("UPDATE employees SET fio = $1, salary = $2, age = $3 WHERE id = $4",
			fio, salary, age, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func deleteEmployeeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		id := r.URL.Path[len("/delete/"):]

		_, err := db.Exec("DELETE FROM employees WHERE id = $1", id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func addEmployeeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		fio := r.FormValue("fio")
		salary := r.FormValue("salary")
		age := r.FormValue("age")

		_, err := db.Exec("INSERT INTO employees (fio, salary, age) VALUES ($1, $2, $3)",
			fio, salary, age)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
