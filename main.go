package main

import (
        "database/sql"
        "encoding/json"
        "fmt"
        // "io/ioutil"
        "log"
        "net/http"
        "os"
        "time"
        // "strings"
        _ "github.com/go-sql-driver/mysql"
        "github.com/gorilla/handlers"
        "github.com/gorilla/mux"
)

func connect() (*sql.DB, error) {
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	database := os.Getenv("DB_NAME") // optional, for flexibility

	if password == "" || host == "" {
		return nil, fmt.Errorf("missing DB credentials")
	}

	dsn := fmt.Sprintf("root:%s@tcp(%s:3306)/%s", password, host, database)
	return sql.Open("mysql", dsn)
}

func blogHandler(w http.ResponseWriter, r *http.Request) {
        db, err := connect()
        if err != nil {
                w.WriteHeader(500)
                return
        }
        defer db.Close()

        rows, err := db.Query("SELECT title FROM blog")
        if err != nil {
                w.WriteHeader(500)
                return
        }
        var titles []string
        for rows.Next() {
                var title string
                err = rows.Scan(&title)
                titles = append(titles, title)
        }
        json.NewEncoder(w).Encode(titles)
}

func main() {
        log.Print("Prepare db...")
        if err := prepare(); err != nil {
                log.Fatal(err)
        }

        log.Print("Listening 8000")
        r := mux.NewRouter()
        r.HandleFunc("/", blogHandler)
        log.Fatal(http.ListenAndServe(":8000", handlers.LoggingHandler(os.Stdout, r)))
}

func prepare() error {
        db, err := connect()
        if err != nil {
                return err
        }
        defer db.Close()

        for i := 0; i < 60; i++ {
                if err := db.Ping(); err == nil {
                        break
                }
                time.Sleep(time.Second)
        }

        if _, err := db.Exec("DROP TABLE IF EXISTS blog"); err != nil {
                return err
        }

        if _, err := db.Exec("CREATE TABLE IF NOT EXISTS blog (id int NOT NULL AUTO_INCREMENT, title varchar(255), PRIMARY KEY (id))"); err != nil {
                return err
        }

        for i := 0; i < 5; i++ {
                if _, err := db.Exec("INSERT INTO blog (title) VALUES (?);", fmt.Sprintf("Blog post #%d", i)); err != nil {
                        return err
                }
        }
        return nil
}
