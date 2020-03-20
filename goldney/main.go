package main

import (
//   "database/sql"
  models "jameselgar.com/goldney/models"
  secret "jameselgar.com/goldney/secret"
  "github.com/joho/godotenv"
  "log"
  "os"
  "fmt"
)

// Stuct to store environment variables for application
type Env struct {
    db models.Datastore // Currently in production this is postgres
    da models.AssetsDatastore // Currently in production this is s3
}

func main() {
    // Setup inital config for sever
    // config.InitConfig()
    err := godotenv.Load()
    if err != nil {
      log.Fatal("Error loading .env file")
    }

    // Set up connection to postgres
    db, err := models.InitDB(os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD") )
    if err != nil {
        log.Panic(err)
    }

    da, err := models.InitAssetsDatastore(secret.AWS_KEY_ID, secret.AWS_KEY, secret.AWS_TOKEN)
    if err != nil {
        fmt.Println("I relly shouldn't be paiced here")
        log.Panic(err)
    }

    env := &Env{db: db, da: da}

    //Gin stuff
    fmt.Println("Im gonna go setup the routers")
    r := SetupRouter(env)
    r.Run(":8080")
}
