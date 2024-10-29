package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	_ "github.com/joho/godotenv/autoload"
)

var tmpl *template.Template
var db *sql.DB

var Store = sessions.NewCookieStore([]byte("usermanagementsecret"))

func init() {
	tmpl, _ = template.ParseGlob("templates/*.html")

	//Set up Sessions
	Store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   3600 * 3,
		HttpOnly: true,
	}

}

var (
	database = os.Getenv("MYSQL_DATABASE")
	password = os.Getenv("MYSQL_PASSWORD")
	username = os.Getenv("MYSQL_USER")
	host     = "192.168.0.101"
	port     = "3306"
)

func initDB() {

	var err error
	// Initialize the db variable
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", username, password, host, port, database)
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		fmt.Println("Database connection failed 2", err)
		log.Fatal(err)
	}

	// Check the database connection
	if err = db.Ping(); err != nil {
		fmt.Println("Database connection failed")
		log.Fatal(err)
	}
}

func main() {

	gRouter := mux.NewRouter()

	//Setup MySQL
	initDB()
	defer db.Close()

	// Setup Static file handling for images

	fileServer := http.FileServer(http.Dir("./uploads"))
	gRouter.PathPrefix("/uploads/").Handler(http.StripPrefix("/uploads", fileServer))

	//All dynamic routes

	//gRouter.HandleFunc("/", handlers.Homepage(db, tmpl, Store))
	//
	//gRouter.HandleFunc("/register", handlers.RegisterPage(db, tmpl)).Methods("GET")
	//
	//gRouter.HandleFunc("/register", handlers.RegisterHandler(db, tmpl)).Methods("POST")
	//
	//gRouter.HandleFunc("/login", handlers.LoginPage(db, tmpl)).Methods("GET")
	//
	//gRouter.HandleFunc("/login", handlers.LoginHandler(db, tmpl, Store)).Methods("POST")
	//
	//gRouter.HandleFunc("/edit", handlers.Editpage(db, tmpl, Store)).Methods("GET")
	//
	//gRouter.HandleFunc("/edit", handlers.UpdateProfileHandler(db, tmpl, Store)).Methods("POST")
	//
	//gRouter.HandleFunc("/upload-avatar", handlers.AvatarPage(db, tmpl, Store)).Methods("GET")
	//
	//gRouter.HandleFunc("/upload-avatar", handlers.UploadAvatarHandler(db, tmpl, Store)).Methods("POST")
	//
	//gRouter.HandleFunc("/logout", handlers.LogoutHandler(Store)).Methods("GET")

	err := http.ListenAndServe(":4002", gRouter)

	if err != nil {
		fmt.Println("Error starting server: ", err)
		log.Fatal(err)
	}

}
