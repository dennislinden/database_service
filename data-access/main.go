package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB
var tpl *template.Template

//struct for database
type Tag struct {
	Name    sql.NullString
	Age     sql.NullInt64
	Number  sql.NullInt64
	Address sql.NullString
}

//struct for website data
type Sub struct {
	Name    string
	Age     string
	Number  string
	Address string
}

func connectToDatabase() {
	// Get a database handle.
	var err error
	db, err = sql.Open("mysql", "root:my-password@tcp(127.0.0.1:3306)/users")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Database Connected!")
}

func insertToDatabase(name string, age int64, number int64, address string) {
	// perform a db.Query insert
	insert, err := db.Query("INSERT INTO test VALUES(?, ?, ?, ?)", name, age, number, address)

	// be careful deferring Queries if you are using transactions
	defer insert.Close()

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Data added!")
}

func executeQuery() {
	// Execute the query
	results, err := db.Query("SELECT Name, Age, Number, Address  FROM test")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	for results.Next() {
		var tag Tag
		// for each row, scan the result into our tag composite object
		err = results.Scan(&tag.Name, &tag.Age, &tag.Number, &tag.Address)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
		// and then print out the tag's Name attribute
		log.Println(tag.Name, tag.Age, tag.Number, tag.Address)
	}
}

func processGetHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("processGetHandler running")
	var s Sub
	s.Name = r.FormValue("userName")
	s.Age = r.FormValue("ageValue")
	s.Number = r.FormValue("numberValue")
	s.Address = r.FormValue("addressValue")
	fmt.Println("Username:", s.Name, "Age:", s.Age, "Number:", s.Number, "Address:", s.Address)
	tpl.ExecuteTemplate(w, "getform.html", s)
	insertToDatabase(string(s.Name), 12, 42334, s.Address)
}

func main() {
	connectToDatabase()
	tpl, _ = tpl.ParseGlob("templates/*.html")
	http.HandleFunc("/processget", processGetHandler)
	http.ListenAndServe(":8082", nil) //startup html service
}
