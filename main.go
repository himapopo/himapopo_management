package main

import (
	"net/http"
	"html/template"
	"database/sql"
	"log"
	"fmt"
	"strconv"
	"time"
	
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
	tm := time.Now()
	y := tm.Year()
	m := int(tm.Month())
	d := tm.Day()
	var dd string
	if d < 10{
		dd = strconv.Itoa(d)
		dd = "0" + dd 
	} else {
		dd = strconv.Itoa(d)
	}
	var mo string
	if m < 10{
		mo = strconv.Itoa(m)
		mo = "0" + mo
	} else {
		mo = strconv.Itoa(m)
	}
	yy := strconv.Itoa(y)
	datatime := yy + "-" + mo + "-" + dd
	Dbconnection, err := sql.Open("sqlite3", "./himapopo.sql")
	if err != nil {
		fmt.Println(err)
		log.Fatalln(err)
	}
	defer Dbconnection.Close()
	cmd := "SELECT * FROM management WHERE created_datetime LIKE ? ORDER BY id desc"
	rows, err := Dbconnection.Query(cmd, "%"+datatime+"%")
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
	t := template.Must(template.ParseFiles("views/home.html"))
	t.ExecuteTemplate(write, "home.html", mm)
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
	
	cmd := "SELECT * FROM management ORDER BY id desc LIMIT 30"
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

func sortHandler(write http.ResponseWriter, request *http.Request){
	year := request.FormValue("year")
	time := request.FormValue("time")
	datatime := year + "-" + time
	Dbconnection, err := sql.Open("sqlite3", "./himapopo.sql")
	if err != nil {
		fmt.Println(err)
		log.Fatalln(err)
	}
	defer Dbconnection.Close()
	
	cmd := "SELECT * FROM management WHERE created_datetime LIKE ? ORDER BY id desc"
	rows, err := Dbconnection.Query(cmd, "%"+datatime+"%")
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

	t := template.Must(template.ParseFiles("views/sort.html"))
	t.ExecuteTemplate(write, "sort.html",mm)
} 

// アップデート
func updateHandler(write http.ResponseWriter, request *http.Request){
	mgm := request.URL.Path[len("/update/"):]
	intmgm,_ := strconv.Atoi(mgm)
	wei := request.FormValue("weight")
	see := request.FormValue("seed")
	pel := request.FormValue("pellet")
	me := request.FormValue("memo")
	Dbconnection, err := sql.Open("sqlite3", "./himapopo.sql")
	if err != nil {
		log.Println(err)
	}
	defer Dbconnection.Close()
	cmd := "UPDATE management SET weight = ?, seed = ?, pellet = ?, memo = ? WHERE id = ?"
	_, err = Dbconnection.Exec(cmd, wei, see, pel, me, intmgm)
	if err != nil {
		log.Println(err)
	}
	http.Redirect(write, request, "/home/", http.StatusFound)
}

// 編集ページ
func editHandler(write http.ResponseWriter, request *http.Request){
	mgm := request.URL.Path[len("/edit/"):]
	intmgm,_ := strconv.Atoi(mgm)
	Dbconnection, _ := sql.Open("sqlite3", "./himapopo.sql")
	defer Dbconnection.Close()
	cmd := "SELECT * FROM management WHERE id = ?"
	row := Dbconnection.QueryRow(cmd,intmgm)
	var m Management
	err := row.Scan(&m.Id, &m.Name, &m.Weight, &m.Seed, &m.Pellet, &m.Memo, &m.Created_datetime)
	if err != sql.ErrNoRows{
		log.Println("no rows!!!")
	} else {
		log.Println(err)
	}
	t := template.Must(template.ParseFiles("views/edit.html"))
	t.ExecuteTemplate(write, "edit.html", m)
}


func main(){
	Dbconnection, err := sql.Open("sqlite3", "./himapopo.sql")
	if err != nil {
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
		log.Fatalln(err)
	}
	http.HandleFunc("/sort/", sortHandler)
	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("images/"))))
	http.HandleFunc("/update/", updateHandler)
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/index/", indexHandler)
	http.HandleFunc("/create/", createHandler)
	http.HandleFunc("/new/", newHandler)
	http.HandleFunc("/home/", homeHandler)
	http.ListenAndServe(":8000", nil)
}