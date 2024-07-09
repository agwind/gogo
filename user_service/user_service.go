package user_service

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/yaml"
	_ "github.com/mattn/go-sqlite3"
)

type database struct {
	db        *sql.DB
	user_stmt *sql.Stmt
}

type user struct {
	Id         int
	First_name string
	Last_name  string
	Email      string
}

var db database

func getUser(w http.ResponseWriter, r *http.Request) {
	user_param := r.PathValue("user")

	fmt.Printf("got /user request for user [%s]\n", user_param)
	user_id, err := strconv.Atoi(user_param)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "Bad user id\n")
		return
	}

	var u user
	err = db.user_stmt.QueryRow(user_id).Scan(&u.Id, &u.First_name, &u.Last_name, &u.Email)
	if err == sql.ErrNoRows {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("User email: %s\n", u.Email)
	}
	user_json, err := json.Marshal(u)
	if err != nil {
		log.Fatal(err)
	}
	io.WriteString(w, string(user_json))
}

func setup_db() {
	database, err := sql.Open("sqlite3", config.String("db_file"))
	if err != nil {
		log.Fatal(err)
	}
	db.db = database
	db.user_stmt, err = db.db.Prepare("select id, first_name, last_name, email from users where id = ?")
	if err != nil {
		log.Fatal(err)
	}
}

func shutdown_db() {
	db.user_stmt.Close()
	db.db.Close()
}

func main() {
	config.AddDriver(yaml.Driver)

	err := config.LoadFiles("config.yaml")
	if err != nil {
		panic(err)
	}
	// fmt.Printf("config data: \n %#v\n", config.Data())

	setup_db()
	defer shutdown_db()

	http.HandleFunc("/user/{user}", getUser)

	port := config.String("port")
	fmt.Printf("Listening on port: " + port + "\n")
	err = http.ListenAndServe(":"+port, nil)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Println("Server closed")
	} else if err != nil {
		fmt.Printf("Error starting server: %s\n", err)
		os.Exit(1)
	}
}
