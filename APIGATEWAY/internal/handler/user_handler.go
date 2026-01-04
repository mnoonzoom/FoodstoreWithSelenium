package handler

import (
	"apigateway/internal/middleware"
	"context"
	"net/http"

	userPB "apigateway/proto/user"
	"github.com/gin-gonic/gin"
)

func InitUserRoutes(r *gin.Engine, client userPB.UserServiceClient) {
	r.POST("/register", func(c *gin.Context) {
		var req userPB.RegisterRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		res, err := client.Register(context.Background(), &req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"user_id": res.Id})
	})

	r.POST("/login", func(c *gin.Context) {
		var req userPB.LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		res, err := client.Login(context.Background(), &req)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, res)
	})
	protected := r.Group("/users")
	protected.Use(middleware.JWTAuthMiddleware())

	protected.GET("/:id", func(c *gin.Context) {
		id := c.Param("id")
		res, err := client.GetUser(c, &userPB.GetUserRequest{Id: id})
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, res.User)
	})
}
