package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/gorilla/mux"
	_ "github.com/joho/godotenv/autoload"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	_ "github.com/golang-migrate/migrate/v4"
)

var tmpl *template.Template
var db *sql.DB

func init() {
	tmpl, _ = template.ParseGlob("templates/*.html")
}

type Task struct {
	Id   int
	Task string
	Done bool
}

var (
	database = os.Getenv("MYSQL_DATABASE")
	username = os.Getenv("MYSQL_USER")
	password = os.Getenv("MYSQL_PASSWORD")
)

func initDB() {
	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/%s?parseTime=true", username, password, database)
	db, err = sql.Open("mysql", dsn) // Correct assignment to package variable
	if err != nil {
		fmt.Println("Error opening database")
		log.Fatal(err)
	}

	// Check the database connection
	if err = db.Ping(); err != nil {
		fmt.Println("Error pinging database")
		log.Fatal(err)
	}
}

func main() {
	gRouter := mux.NewRouter()

	// Setup MySQL
	initDB()
	defer func() {
		if err := db.Close(); err != nil {
			fmt.Println("Error closing database")
			log.Fatal(err)
		}
	}()

	gRouter.HandleFunc("/", Homepage)
	gRouter.HandleFunc("/tasks", fetchTasks).Methods("GET")
	gRouter.HandleFunc("/newtaskform", getTaskForm)
	gRouter.HandleFunc("/tasks", addTask).Methods("POST")
	gRouter.HandleFunc("/gettaskupdateform/{id}", getTaskUpdateForm).Methods("GET")
	gRouter.HandleFunc("/tasks/{id}", updateTask).Methods("PUT", "POST")
	gRouter.HandleFunc("/tasks/{id}", deleteTask).Methods("DELETE")

	gRouter.NotFoundHandler = http.HandlerFunc(notFoundHandler)

	err := http.ListenAndServe(":4001", gRouter)
	if err != nil {
		fmt.Println("Error starting server")
		log.Fatal(err)
	}
	fmt.Println("Server started on port 4001")
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	//http.Error(w, "Page not found", http.StatusNotFound)
	tmpl.ExecuteTemplate(w, "404.html", nil)
}

func Homepage(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "home.html", nil)
}

func fetchTasks(w http.ResponseWriter, r *http.Request) {
	todos, err := getTasks(db)
	if err != nil {
		http.Error(w, "Error fetching tasks", http.StatusInternalServerError)
		return
	}
	tmpl.ExecuteTemplate(w, "todoList", todos)
}

func getTaskForm(w http.ResponseWriter, r *http.Request) {

	tmpl.ExecuteTemplate(w, "addTaskForm", nil)
}

func addTask(w http.ResponseWriter, r *http.Request) {

	task := r.FormValue("task")

	fmt.Println(task)

	query := "INSERT INTO tasks (task, done) VALUES (?, ?)"

	stmt, err := db.Prepare(query)

	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	_, executeErr := stmt.Exec(task, 0)

	if executeErr != nil {
		log.Fatal(executeErr)
	}

	// Return a new list of Todos
	todos, _ := getTasks(db)

	//You can also just send back the single task and append it
	//I like returning the whole list just to get everything fresh, but this might not be the best strategy
	tmpl.ExecuteTemplate(w, "todoList", todos)

}

func getTaskUpdateForm(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	//Convert string id from URL to integer
	taskId, _ := strconv.Atoi(vars["id"])

	task, err := getTaskByID(db, taskId)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	tmpl.ExecuteTemplate(w, "updateTaskForm", task)

}

func updateTask(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	taskItem := r.FormValue("task")
	//taskStatus, _ := strconv.ParseBool(r.FormValue("done"))
	var taskStatus bool

	fmt.Println(r.FormValue("done"))

	//Check the string value of the checkbox
	switch strings.ToLower(r.FormValue("done")) {
	case "yes", "on":
		taskStatus = true
	case "no", "off":
		taskStatus = false
	default:
		taskStatus = false
	}

	taskId, _ := strconv.Atoi(vars["id"])

	task := Task{
		taskId, taskItem, taskStatus,
	}

	updateErr := updateTaskById(db, task)

	if updateErr != nil {
		log.Fatal(updateErr)
	}

	//Refresh all Tasks
	todos, _ := getTasks(db)

	tmpl.ExecuteTemplate(w, "todoList", todos)
}

func deleteTask(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	taskId, _ := strconv.Atoi(vars["id"])

	err := deleteTaskWithID(db, taskId)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	//Return list
	todos, _ := getTasks(db)

	tmpl.ExecuteTemplate(w, "todoList", todos)
}

func getTasks(dbPointer *sql.DB) ([]Task, error) {
	query := "SELECT id, task, done FROM tasks"
	rows, err := dbPointer.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var todo Task
		rowErr := rows.Scan(&todo.Id, &todo.Task, &todo.Done)
		if rowErr != nil {
			return nil, rowErr
		}
		tasks = append(tasks, todo)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return tasks, nil
}

func getTaskByID(dbPointer *sql.DB, id int) (*Task, error) {

	query := "SELECT id, task, done FROM tasks WHERE id = ?"

	var task Task

	row := dbPointer.QueryRow(query, id)
	err := row.Scan(&task.Id, &task.Task, &task.Done)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("No task was found with task %d", id)
		}
		return nil, err
	}

	return &task, nil

}

func updateTaskById(dbPointer *sql.DB, task Task) error {

	query := "UPDATE tasks SET task = ?, done = ? WHERE id = ?"

	result, err := dbPointer.Exec(query, task.Task, task.Done, task.Id)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		fmt.Println("No rows updated")
	} else {
		fmt.Printf("%d row(s) updated\n", rowsAffected)
	}

	return nil

}

func deleteTaskWithID(dbPointer *sql.DB, id int) error {

	query := "DELETE FROM tasks WHERE id = ?"

	stmt, err := dbPointer.Prepare(query)

	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no task found with id %d", id)
	}

	fmt.Printf("Deleted %d task(s)\n", rowsAffected)
	return nil

}
