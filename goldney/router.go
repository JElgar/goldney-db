package main

import (
   "github.com/gin-gonic/gin"
)

func SetupRouter(env *Env) *gin.Engine {
    r := gin.Default()

    // This was required for local testing
    //r.Use(cors.Default())
    // guessing this is pretty handy for version control :D
    api := r.Group("api/v1")
    
    // An end point to test connection works
    api.GET("/ping", ping)

    return r
}

func ping(c *gin.Context){
    c.JSON(200, gin.H{
        "world": "Hello",
    })
}

