package main

import (
//   "database/sql"
  models "jameselgar.com/goldney/models"
  "log"
  "os"
)

// Stuct to store environment variables for application
type Env struct {
    db models.Datastore // Currently in production this is postgres
    da models.AssetsDatastore // Currently in production this is s3
}

func main() {
    // Setup inital config for sever
    // config.InitConfig()

    // Set up connection to postgres
    db, err := models.InitDB(os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD") )
    if err != nil {
        log.Panic(err)
    }

    da, err := models.InitAssetsDatastore()
    env := &Env{db: db, da: da}

    //Gin stuff
    r := SetupRouter(env)
    r.Run(":8080")
}
