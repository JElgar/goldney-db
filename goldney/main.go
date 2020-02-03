package main

import (
//   "database/sql"
  models "github.com/jelgar/goLdney/goldney-ar/goldney/models"
  "log"
)

// Stuct to store environment variables for application
type Env struct {
    db models.Datastore
}

func main() {
    // Setup inital config for sever
    // config.InitConfig()

    // Set up connection to postgres
    db, err := models.InitDB("fakeString")
    // Example connection string (this was for aws)
    // postgresql://admin:test123@ec2-35-178-198-24.eu-west-2.compute.amazonaws.com/secta?sslmode=disable
    if err != nil {
        log.Panic(err)
    }
    env := &Env{db: db}

    //Gin stuff
    r := SetupRouter(env)
    r.Run(":8080")
}
