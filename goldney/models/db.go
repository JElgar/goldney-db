package models

import (
  "database/sql"
 _ "github.com/lib/pq"
)

// This is the datastore interface which states what fucntions a datastore must define. This is very useful for testing as a "fake" datastore can be created.
type Datastore interface {
  TileStore
}


// This is the main DB which whill inherit the above interface
type DB struct {
  *sql.DB
}


func InitDB(DbHost string, DbUser string, DbPassword string) (*DB, error) {
//    psqlInfo := fmt.Sprintf("host=%s port=%d user=%s"+ "password=%s dbname=%s sslmode=disable", 
//    DbHost, "5432", DbUser, DbPassword, "test")
//psqlInfo := "postgres://"+ DbUser + ":" + DbPassword + "@db/goldney?sslmode=disable"
    psqlInfo := "postgres://test:test@db/goldney?sslmode=disable"
//  db, err := sql.Open("postgres", DbUser + ":" + DbPassword + "@" + DbHost)
    db, err := sql.Open("postgres", psqlInfo)
    if err != nil {
        return nil, err
    }

    if err = db.Ping(); err != nil {
        return nil, err
    }

    return &DB{db}, nil
}
