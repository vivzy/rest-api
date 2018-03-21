package main

import (
  "github.com/gin-gonic/gin"
  "github.com/appleboy/gin-jwt"
  "rest-api/app"
  "time"
)

func SetupRouter() *gin.Engine {
  router := gin.Default()

  // Middleware for Auth
  authMiddleware := &jwt.GinJWTMiddleware{
    Realm:      "test zone",
    Key:        []byte("verybigsecret"),
    Timeout:    time.Hour,
    MaxRefresh: time.Hour,
    Authenticator: app.AuthUsers,
    Authorizator: func(userId string, c *gin.Context) bool {
      return true
    },
    Unauthorized: func(c *gin.Context, code int, message string) {
      c.JSON(code, gin.H{
        "code":    code,
        "message": message,
      })
    },
    TokenLookup: "header:Authorization",
    TokenHeadName: "Bearer",
    TimeFunc: time.Now,
  }

  // Routes
  router.POST("/login", authMiddleware.LoginHandler)
  router.POST("/adduser", app.PostUser)

  // Routes with Auth
  v1 := router.Group("api/v1")
  v1.Use(authMiddleware.MiddlewareFunc())
  {

    v1.GET("/users", app.GetUsers)
    v1.GET("/users/:id", app.GetUser)
    v1.PUT("/users/:id" , app.UpdateUser)
    v1.DELETE("/users/:id", app.DeleteUser)
    v1.GET("/refresh_token", authMiddleware.RefreshHandler)
  }

  router.Use(gin.Logger())
  return router
}

func main() {
  router := SetupRouter()
  router.Run(":8080")
}
