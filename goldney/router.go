package main

import (
   "github.com/gin-gonic/gin"
   models "jameselgar.com/goldney/models"
)

func SetupRouter(env *Env) *gin.Engine {
    r := gin.Default()

    // This was required for local testing
    //r.Use(cors.Default())
    // guessing this is pretty handy for version control :D
    api := r.Group("v1")
    
    // An end point to test onnection works
    api.GET("/ping", ping)
    api.POST("/newTile", env.newTile)
    api.GET("/getTiles", env.getTiles)

    return r
}

func ping(c *gin.Context){
    c.JSON(200, gin.H{
        "world": "Hello",
    })
}

func (e *Env) newTile (c *gin.Context){
    // TODO on success return user and enventually JSON web token
    var t models.Tile
    c.BindJSON(&t)

    _, err := e.db.CreateTile(&t)
    if err != nil && err.Code == 409 {
        c.JSON(err.Code, err)
        return
    } else if err != nil {
        c.JSON(err.Code, err)
        return
    }
}

func (e *Env) getTiles (c *gin.Context) {
    e.db.GetTiles()
}
