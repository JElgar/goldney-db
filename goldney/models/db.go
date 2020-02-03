package models

import (
  "database/sql"
  _ "github.com/lib/pq"
)

// This is the datastore interface which states what fucntions a datastore must define. This is very useful for testing as a "fake" datastore can be created.
type Datastore interface {

}


// This is the main DB which whill inherit the above interface
type DB struct {
  *sql.DB
}


func InitDB(DbUrl string) (*DB, error) {
    db, err := sql.Open("postgres", DbUrl)
    if err != nil {
        return nil, err
    }

    if err = db.Ping(); err != nil {
        return nil, err
    }

    return &DB{db}, nil
}
