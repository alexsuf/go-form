package main

import (
	"html/template"
	"net/http"
	"strconv"
)

type FormData struct {
	Fio    string
	Salary int
	Age    int
	Date   string
}

var formTemplate = `
<!DOCTYPE html>
<html>
<head>
	<title>Форма ввода</title>
</head>
<body>
	<h1>Введите данные</h1>
	<form method="POST">
		<label>Фамилия и имя:</label>
		<input type="text" name="fio" required><br><br>
		<label>Зарплата:</label>
		<input type="number" name="salary" required><br><br>
		<label>Возраст:</label>
		<input type="number" name="age" required><br><br>
		<input type="submit" value="Отправить">
	</form>
	{{ if .Fio }}
		<h2>Вы ввели:</h2>
		<p>Фамилия и имя: {{ .Fio }}</p>
		<p>Зарплата: {{ .Salary }}</p>
		<p>Возраст: {{ .Age }}</p>
	{{ end }}
</body>
</html>
`

func main() {
	http.HandleFunc("/", indexHandler)
	http.ListenAndServe(":8080", nil)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	var data FormData

	if r.Method == http.MethodPost {
		r.ParseForm()
		// Преобразуем строку в int
		if salary, err := strconv.Atoi(r.FormValue("salary")); err == nil {
			data.Salary = salary
		} else {
			http.Error(w, "Неверное значение зарплаты", http.StatusBadRequest)
			return
		}
		if age, err := strconv.Atoi(r.FormValue("age")); err == nil {
			data.Age = age
		} else {
			http.Error(w, "Неверный возраст", http.StatusBadRequest)
			return
		}
		data.Fio = r.FormValue("fio")
	}

	tmpl, err := template.New("form").Parse(formTemplate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, data)
}
