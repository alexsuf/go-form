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
		h1 {
			text-align: center;
		}
		table {
			width: 80%;
			margin: 0 auto;
			border-collapse: collapse;
		}
		th, td {
			border: 1px solid white;
			padding: 10px;
			text-align: left;
		}
		th {
			background-color: #444;
		}
	</style>
</head>
<body>
	<h1>Список сотрудников</h1>
	<table>
		<tr>
			<th>ID</th>
			<th>ФИО</th>
			<th>Зарплата</th>
			<th>Возраст</th>
		</tr>
		{{ range . }}
			<tr>
				<td>{{ .ID }}</td>
				<td>{{ .FIO }}</td>
				<td>{{ .Salary }}</td>
				<td>{{ .Age }}</td>
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

	http.HandleFunc("/", employeeHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func employeeHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, fio, salary, age FROM employees")
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
