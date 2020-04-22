package main

import (
   "github.com/gin-gonic/gin"
   models "jameselgar.com/goldney/models"
   errors "jameselgar.com/goldney/errors"
  "github.com/gin-contrib/cors"
  "fmt"
  "time"
  "net/http"
)

func SetupRouter(env *Env) *gin.Engine {
    //gin.SetMode(gin.releaseMode)
    r := gin.Default()
    //r.Use(cors.New(cors.Config{
	//	AllowOrigins:     []string{"https://goldneyhall.com", "https://create.goldneyhall.com"},
	//	AllowMethods:     []string{"PUT", "PATCH", "GET"},
	//	AllowHeaders:     []string{"Origin"},
	//}))
    r.Use(cors.New(cors.Config{
      AllowOrigins:     []string{"https://create.goldneyhall.com", "https://goldneyhall.com", "http://goldneyhall.com", "http://create.goldneyhall.com"},
        //AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"POST", "PUT", "PATCH"},
		AllowHeaders:     []string{"Origin"},
        AllowCredentials: true,
		ExposeHeaders:    []string{"Content-Length"},
	}))

    //r.Use(cors.Default())

    // This was required for local testing
    //r.Use(cors.Default())
    // guessing this is pretty handy for version control :D
    api := r.Group("v1")
    tokenAuth := api.Group("/")
    tokenAuth.Use(AuthRequired())
    
    // An end point to test onnection works
    api.GET("/ping", ping)
    api.POST("/newTile", env.newTile)
    api.GET("/getTiles", env.getTiles)
    
    api.POST("/uploadImage", env.uploadImage)
    api.POST("/uploadAudio", env.uploadAudio)

    api.POST("/login", env.login)
   
    // Admin only function
    tokenAuth.POST("/updateTile", env.updateTile)
    tokenAuth.POST("/deleteSections", env.deleteSections)
    tokenAuth.POST("/toggleActivateTile", env.setTileActive)
    tokenAuth.GET("/getAllTiles", env.getAllTiles)

    return r
}

func ping(c *gin.Context){
    c.JSON(200, gin.H{
        "world": "Hello",
    })
}

func (e *Env) newTile (c *gin.Context){
    var t models.Tile
    binderr := c.ShouldBindJSON(&t)
    if binderr != nil {
      fmt.Println("The Binding was not a perfect match")
    }

    _, err := e.db.AddTile(&t)
    if err != nil {
      c.JSON(err.Code, err)
      return
    }
    c.JSON(200, "Success")
}

func (e *Env) getTiles (c *gin.Context) {
  tiles, err := e.db.GetActiveTiles()
  if err != nil {
    panic(err)
    c.JSON(err.Code, err)
    return
  }
  c.JSON(200, tiles)
}

func (e *Env) getAllTiles (c *gin.Context) {
  tiles, err := e.db.GetAllTiles()
  if err != nil {
    panic(err)
    c.JSON(err.Code, err)
    return
  }
  c.JSON(200, tiles)
}

func (e *Env) updateTile (c *gin.Context) {
  var t models.Tile
  binderr := c.ShouldBindJSON(&t)
  if binderr != nil {
    fmt.Println("The Binding was not a perfect match")
  }
  fmt.Println(t)

  _, err := e.db.UpdateTile(&t)
  if err != nil && err.Code == 409 {
      c.JSON(err.Code, err)
      return
  } else if err != nil {
      c.JSON(err.Code, err)
      return
  }
  c.JSON(200, "Successfully updated")
}

func (e * Env) uploadImage (c *gin.Context) {
  fmt.Println("uplaoding image endpoint reached")
  fileHeader, err := c.FormFile("undefined")
  if err != nil {
    panic(err)
    c.JSON(401, err)
    return
  }
  f, err := fileHeader.Open()
  if err != nil {
    c.JSON(402, err)
    return
  }
  fmt.Println("doing the uploaidng")
  fileName, err := e.da.ImageStore(f, fileHeader)
  fmt.Println("did the uploaidng")
  if err != nil {
    fmt.Println("panicing error")
    panic(err)
    c.JSON(403, err)
    return
  }
  fmt.Println("uplaoded the image")
  fmt.Println(fileName)
  c.JSON(200, fileName)

}

func (e * Env) uploadAudio (c *gin.Context) {
  fmt.Println("uplaoding image endpoint reached")
  fileHeader, err := c.FormFile("undefined")
  if err != nil {
    panic(err)
    c.JSON(401, err)
    return
  }
  f, err := fileHeader.Open()
  if err != nil {
    c.JSON(402, err)
    return
  }
  fmt.Println("doing the uploaidng")
  fileName, err := e.da.AudioStore(f, fileHeader)
  fmt.Println("did the uploaidng")
  if err != nil {
    fmt.Println("panicing error")
    panic(err)
    c.JSON(403, err)
    return
  }
  fmt.Println("uplaoded the image")
  fmt.Println(fileName)
  c.JSON(200, fileName)

}

func (e *Env) setTileActive (c *gin.Context) {
  fmt.Println("Toggling the tile active")
  var t models.Tile
  binderr := c.ShouldBindJSON(&t)
  if binderr != nil {
    fmt.Println("The Binding was not a perfect match")
  }
  fmt.Println(t)
  err := e.db.SetActive(&t)
  if err != nil {
    c.JSON(400, err)
  }
  c.JSON(200, "Successfully updated active")
}

func (e *Env) deleteSections (c *gin.Context) {
  fmt.Println("Deleting given section")
  //var sections = struct {
  //  ids []int  `json:"sections"`
  //}{
  //  []int{},
  //}
  //type DelSections struct {
  //  Ids []int `json:"secs"`
  //}
  var t models.Tile
  binderr := c.ShouldBindJSON(&t)
  if binderr != nil {
    fmt.Println("The Binding was not a perfect match")
  }
  fmt.Println(t)
  fmt.Println(t.DeleteSecs)
  for _, id := range t.DeleteSecs {
    fmt.Println()
    err := e.db.DeleteSection(id)
    if err != nil {
      if err.Err != nil {
        panic(err.Err)
      }
    }
  }
  //err := e.db.SetActive(&t)
  //if err != nil {
  //  c.JSON(400, err)
  //}
  c.JSON(200, "Successfully updated active")
}

func (e *Env) login (c *gin.Context) {
    var u models.User
    c.BindJSON(&u)

    err := e.db.Login(&u)
    if err != nil {
      panic(err)
      return
    }
    expirationTime := time.Now().Add(5 * time.Minute)

    claims := &models.Claims {
      Username: u.Username,
      StandardClaims: models.NewStandardClaims(expirationTime),
    }
    token := models.NewJWT(models.DefaultSignMethod(), claims)

    tokenString, erro := token.SignedString(models.GetKey())
    if erro != nil {
        fmt.Println("Error making token into signed string")
    }
    // TODO making these both false seemed to fix an issue but i dont want them to both be false im guessing 
    fmt.Println("Setting cookie")
    c.SetCookie(
        "token",
        tokenString,
        3600,
        "/",
        "",
        false,
        false)
    //c.JSON(200, gin.H{
    //    "token": tokenString,
    //})
}

func AuthRequired() gin.HandlerFunc {
    return func(c *gin.Context) {
        fmt.Println("Hello")
        cookie, err := c.Cookie("token")
        fmt.Println(cookie)

        //cookie, err := c.Request.Cookie("token")
        if err != nil {
            if err == http.ErrNoCookie {
                // If the cookie is not set, return an unauthorized status
                fmt.Println("Cookie not set")
                c.JSON(400, errors.ApiError{err, "No cookie supplied", 400})
                c.Abort()
	            return
            }
            c.JSON(500, errors.ApiError{err, "Error getting cookie", 400})
            c.Abort()
	    	return
	    }
        //tokenString, err := url.QueryUnescape(cookie.Value)
        claims := &models.Claims{}

        token, erro := models.ParseWClaims(cookie, claims)

        if !token.Valid {
            fmt.Println("Invalid Token")
            c.JSON(400, errors.ApiError{err, "Invalid Token", 400})
            c.Abort()
            return
        }

        if erro != nil {
            if erro == models.ErrSignatureInvalid{
                fmt.Println("Invalid Signature")
                c.JSON(400, errors.ApiError{err, "Invalid Signature", 400})
                c.Abort()
                return
            }
            fmt.Println("Unknown error")
            c.JSON(500, errors.ApiError{err, "Unknown error", 400})
            c.Abort()
            return
        }
        fmt.Println("Cookie all good")
        c.Next()
    }
}

