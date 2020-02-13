package main

import (
//   "database/sql"
  models "jameselgar.com/goldney/models"
  "log"
  "os"
)

// Stuct to store environment variables for application
type Env struct {
    db models.Datastore
}

func main() {
    // Setup inital config for sever
    // config.InitConfig()

    // Set up connection to postgres
    db, err := models.InitDB(os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD") )
    if err != nil {
        log.Panic(err)
    }
    env := &Env{db: db}

    //Gin stuff
    r := SetupRouter(env)
    r.Run(":8080")
}
