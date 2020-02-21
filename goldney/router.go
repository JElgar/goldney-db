package main

import (
   "github.com/gin-gonic/gin"
   models "jameselgar.com/goldney/models"
  "github.com/gin-contrib/cors"
  "fmt"
)

func SetupRouter(env *Env) *gin.Engine {
    r := gin.Default()
    //r.Use(cors.New(cors.Config{
	//	AllowOrigins:     []string{"https://goldneyhall.com", "https://create.goldneyhall.com"},
	//	AllowMethods:     []string{"PUT", "PATCH", "GET"},
	//	AllowHeaders:     []string{"Origin"},
	//}))
    r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://create.goldneyhall.com"},
		AllowMethods:     []string{"POST", "PUT", "PATCH"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
	}))

    //r.Use(cors.Default())

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

    fmt.Println(t)

    _, err := e.db.AddTile(&t)
    if err != nil && err.Code == 409 {
        c.JSON(err.Code, err)
        return
    } else if err != nil {
        c.JSON(err.Code, err)
        return
    }
    c.JSON(200, "Success")
}

func (e *Env) getTiles (c *gin.Context) {
  tiles, err := e.db.GetTiles()
  if err != nil {
    c.JSON(err.Code, err)
    return
  }
  c.JSON(200, tiles)
}
