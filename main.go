package main

import (
    "database/sql"
    "fmt"
    _ "github.com/lib/pq"
)

func main() {
	user := ""
	password:= ""
	dbname:= "postgres"
    connStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",user,password,dbname);
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        panic(err)
    }
	defer db.Close()

    err = db.Ping()
    if err != nil {
        panic(err)
    }

    fmt.Printf("Successfully connected to user: %s, database: %s",user,dbname)
}

