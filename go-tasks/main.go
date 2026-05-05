package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/glebarez/go-sqlite"
)

// Структура студента (для БД та JSON)
type Student struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Grade int    `json:"grade"`
}

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("sqlite", "./students.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Створюємо таблицю
	db.Exec(`CREATE TABLE IF NOT EXISTS students (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT,
		grade INTEGER
	)`)

	// 1. ЗАПУСКАЄМО ВЕБ-СЕРВЕР НА ФОНІ (go routine)
	go func() {
		http.HandleFunc("/students", httpHandler)
		fmt.Println("\n[INFO] Веб-сервер запущено на http://localhost:8080/students")
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Printf("Помилка сервера: %v", err)
		}
	}()

	// 2. ОСНОВНИЙ ЦИКЛ ПРОГРАМИ (МЕНЮ)
	for {
		fmt.Println("\n--- МЕНЮ КЕРУВАННЯ (Група 304-ТН) ---")
		fmt.Println("1. Додати студента")
		fmt.Println("2. Список усіх")
		fmt.Println("3. Вийти")
		fmt.Print("Вибір: ")

		var choice int
		fmt.Scan(&choice)

		switch choice {
		case 1:
			addStudent()
		case 2:
			listStudents()
		case 3:
			fmt.Println("Вихід...")
			os.Exit(0)
		}
	}
}

// Функція для роботи через браузер/Postman
func httpHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Повертаємо список студентів у форматі JSON
	rows, _ := db.Query("SELECT id, name, grade FROM students")
	defer rows.Close()

	var list []Student
	for rows.Next() {
		var s Student
		rows.Scan(&s.ID, &s.Name, &s.Grade)
		list = append(list, s)
	}
	json.NewEncoder(w).Encode(list)
}

func addStudent() {
	var name string
	var grade int
	fmt.Print("Ім'я: ")
	fmt.Scan(&name)
	fmt.Print("Оцінка: ")
	fmt.Scan(&grade)

	db.Exec("INSERT INTO students (name, grade) VALUES (?, ?)", name, grade)
	fmt.Println("Додано!")
}

func listStudents() {
	rows, _ := db.Query("SELECT id, name, grade FROM students")
	defer rows.Close()
	fmt.Println("\nID | Студент | Оцінка")
	for rows.Next() {
		var id, grade int
		var name string
		rows.Scan(&id, &name, &grade)
		fmt.Printf("%d | %s | %d\n", id, name, grade)
	}
}
