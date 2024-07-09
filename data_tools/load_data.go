package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/brianvoe/gofakeit"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	var count int
	flag.IntVar(&count, "count", 100, "The number of users to add")

	var data_file string
	flag.StringVar(&data_file, "file", "../data.db", "The data file to load test users")

	flag.Parse()

	fmt.Printf("Count [%d]\n", count)
	db, err := sql.Open("sqlite3", data_file)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := tx.Prepare("insert into users(first_name, last_name, email) values(?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	for i := 0; i < count; i++ {
		first_name := gofakeit.FirstName()
		last_name := gofakeit.LastName()
		insert_count := 0
		var email string
		for {
			if insert_count == 0 {
				email = fmt.Sprintf("%s.%s@example.com", strings.ToLower(first_name[0:1]), strings.ToLower(last_name))
			} else {
				email = fmt.Sprintf("%s.%s.%d@example.com", strings.ToLower(first_name[0:1]), strings.ToLower(last_name), insert_count)
			}
			//fmt.Printf("%s %s %s", first_name, last_name, email)
			_, err = stmt.Exec(first_name, last_name, email)
			if err != nil {
				log.Println(err)
				log.Printf("Email [%s] already used", email)
				insert_count++
			} else if insert_count > 1000 {
				log.Fatalf("Too many inserts for %s %s", first_name, last_name)
			} else {
				break
			}
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}

}
