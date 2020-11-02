package main

import (
	"net/http"
	"html/template"
	"database/sql"
	"log"
	"fmt"
	_"github.com/mattn/go-sqlite3"
)

var Dbconnection *sql.DB

type Management struct{
	Id                  int
	Name             string
	Weight              int
	Seed                int
	Pellet              int
	Memo             string
	Created_datetime string

}

//ホーム画面
func homeHandler(write http.ResponseWriter, request *http.Request){
	t := template.Must(template.ParseFiles("views/home.html"))
	t.ExecuteTemplate(write, "home.html", nil)
}

//データ作成ページ
func newHandler(write http.ResponseWriter, request *http.Request){
	t := template.Must(template.ParseFiles("views/new.html"))
	t.ExecuteTemplate(write, "new.html", nil)
}


//データ記録
func createHandler(write http.ResponseWriter, request *http.Request){
	h_weight := request.FormValue("h_weight")
	h_seed := request.FormValue("h_seed")
	h_pellet := request.FormValue("h_pellet")
	h_memo := request.FormValue("h_memo")
	p_weight := request.FormValue("p_weight")
	p_seed := request.FormValue("p_seed")
	p_pellet := request.FormValue("p_pellet")
	p_memo := request.FormValue("p_memo")
	Dbconnection, err := sql.Open("sqlite3", "./himapopo.sql")
	if err != nil {
		fmt.Println(err)
		log.Fatalln(err)
	}
	defer Dbconnection.Close()
	cmd := "INSERT INTO management (name, weight, seed, pellet, memo) VALUES (?, ?, ?, ?, ?)"
	_, err = Dbconnection.Exec(cmd, "ひま", h_weight, h_seed, h_pellet, h_memo)
	if err != nil {
		fmt.Println(err)
		log.Fatalln(err)
	}
	cmd2 := "INSERT INTO management (name, weight, seed, pellet, memo) VALUES (?, ?, ?, ?, ?)"
	_, err = Dbconnection.Exec(cmd2, "ぽぽ", p_weight, p_seed, p_pellet, p_memo)
	if err != nil {
		fmt.Println(err)
		log.Fatalln(err)
	}
	http.Redirect(write, request, "/home/", http.StatusFound)
}


//一覧表示
func indexHandler(write http.ResponseWriter, request *http.Request){
	Dbconnection, err := sql.Open("sqlite3", "./himapopo.sql")
	if err != nil {
		fmt.Println(err)
		log.Fatalln(err)
	}
	defer Dbconnection.Close()
	
	cmd := "SELECT * FROM management"
	rows, err := Dbconnection.Query(cmd)
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()
	var mm []Management
	for rows.Next(){
		var m Management
		err = rows.Scan(&m.Id, &m.Name, &m.Weight, &m.Seed, &m.Pellet, &m.Memo, &m.Created_datetime)
		if err != nil {
			log.Println(err)
		} 
		mm = append(mm,m)
	}

	t := template.Must(template.ParseFiles("views/index.html"))
	t.ExecuteTemplate(write, "index.html",mm)
} 

// アップデート


// 編集ページ
func editHandler(write http.ResponseWriter, request *http.Request){
	t := template.Must(template.ParseFiles("views/edit.html"))
	t.ExecuteTemplate(write, "edit.html", nil)
}


func main(){
	Dbconnection, err := sql.Open("sqlite3", "./himapopo.sql")
	if err != nil {
		fmt.Println(err)
		log.Fatalln(err)
	}
	defer Dbconnection.Close()
	cmd :=`CREATE TABLE IF NOT EXISTS management(id INTEGER PRIMARY KEY,
																							 name STRING,
																							 weight INT,
																							 seed INT, 
																							 pellet INT, 
																							 memo STRING,
																							 created_datetime STRING DEFAULT (datetime('now','localtime'))
																							 )`
	_, err = Dbconnection.Exec(cmd)
	if err != nil {
		fmt.Println(err)
		log.Fatalln(err)
	}
	//http.HandleFunc("/update/", updateHandler)
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/index/", indexHandler)
	http.HandleFunc("/create/", createHandler)
	http.HandleFunc("/new/", newHandler)
	http.HandleFunc("/home/", homeHandler)
	http.ListenAndServe(":8000", nil)
}