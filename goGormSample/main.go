package main

import (
	"encoding/json"
	"fmt"
	"gorm.io/driver/mysql"
	_ "gorm.io/gen"
	"gorm.io/gorm"
	"io"
	"log"
	"net/http"
)

type Article struct {
	ID      int    `json:"id" gorm:"id"`
	Title   string `json:"title" gorm:"title"`
	Desc    string `json:"desc" gorm:"desc"`
	Content string `json:"content" gorm:"content"`
}

type User struct {
	ID   int    `json:"id" gorm:"id"`
	Name string `json:"name" gorm:"name"`
}

const dbName = "sql_sample"

func Connect() *gorm.DB {
	const userName = "root"
	const password = "heas3real9ract2ZIRK"
	const dsn = userName + ":" + password + "@/" + dbName
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err.Error())
	}

	return db
}

func Disconnect(db *gorm.DB) {
	_db, err := db.DB()
	if err != nil {
		println(err.Error())
	}
	err = _db.Close()
	if err != nil {
		println(err.Error())
	}
	fmt.Println("close db")
}

func handleRequests() {
	server := http.Server{
		Addr: ":8000",
	}
	http.HandleFunc("/", homePage)
	http.HandleFunc("/articles", articles)
	http.HandleFunc("/users", users)
	err := server.ListenAndServe()
	if err != nil {
		panic(err.Error())
	}
}

func homePage(w http.ResponseWriter, r *http.Request) {
	if r.RequestURI == "/favicon.ico" {
		return
	}
	_, err := fmt.Fprintf(w, "Welcome to the HomePage")
	if err != nil {
		log.Fatalf("%v", err)
		return
	}
	fmt.Println("Endpoint Hit: homePage")
}

func articles(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		getArticles(w, r)
	case "POST":
		postArticles(w, r)
	default:
		w.WriteHeader(405)
	}
}

func getArticles(w http.ResponseWriter, _ *http.Request) {
	db := Connect()
	defer Disconnect(db)

	var articles []Article
	result := db.Find(&articles)

	if result.Error != nil {
		panic(result.Error.Error())
	}

	fmt.Println("Endpoint Hit: articles")
	err := json.NewEncoder(w).Encode(articles)
	if err != nil {
		log.Fatalf("%v", err)
	}
	w.WriteHeader(http.StatusOK)
}

func postArticles(w http.ResponseWriter, r *http.Request) {
	db := Connect()
	defer Disconnect(db)

	uri := r.RequestURI
	fmt.Println(uri)
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatalf("%v", err)
	}
	var article Article
	err = json.Unmarshal(body, &article)
	if err != nil {
		log.Fatalf("%v", err)
	}

	// Selectで必要Columnを指定か、Omitで不要なColumn指定
	db.Omit("Id").Create(&article)

	err = json.NewEncoder(w).Encode(map[string]int{"id": article.ID})
	if err != nil {
		panic(err.Error())
	}
	w.WriteHeader(http.StatusCreated)
}

func users(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		getUsers(w, r)
	default:
		w.WriteHeader(405)
	}
}

func getUsers(w http.ResponseWriter, _ *http.Request) {
	db := Connect()
	defer Disconnect(db)

	var users []User
	result := db.Find(&users)

	if result.Error != nil {
		panic(result.Error.Error())
	}

	fmt.Println("Endpoint Hit: articles")
	err := json.NewEncoder(w).Encode(users)
	if err != nil {
		log.Fatalf("%v", err)
	}

	w.WriteHeader(http.StatusOK)
}

func main() {
	handleRequests()
}
