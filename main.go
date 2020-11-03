package main

import (
	"database/sql"
	"log"
	"net/http"
	"himapopo_management/controllers"

	_ "github.com/mattn/go-sqlite3"
)

var Dbconnection *sql.DB

func main() {
	Dbconnection, err := sql.Open("sqlite3", "./himapopo.sql")
	if err != nil {
		log.Fatalln(err)
	}
	defer Dbconnection.Close()
	cmd := `CREATE TABLE IF NOT EXISTS management(id INTEGER PRIMARY KEY,
																							            name STRING,
																							             weight INT,
																							               seed INT, 
																							             pellet INT, 
																							            memo STRING,
					 created_datetime STRING DEFAULT (datetime('now','localtime'))
																							 )`
	_, err = Dbconnection.Exec(cmd)
	if err != nil {
		log.Fatalln(err)
	}
	http.HandleFunc("/sort/", controllers.SortHandler)
	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("images/"))))
	http.HandleFunc("/update/", controllers.UpdateHandler)
	http.HandleFunc("/edit/", controllers.EditHandler)
	http.HandleFunc("/index/", controllers.IndexHandler)
	http.HandleFunc("/create/", controllers.CreateHandler)
	http.HandleFunc("/new/", controllers.NewHandler)
	http.HandleFunc("/home/", controllers.HomeHandler)
	http.HandleFunc("/", controllers.HomeHandler)
	http.ListenAndServe(":8000", nil)
}
