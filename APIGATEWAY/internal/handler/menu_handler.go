package handler

import (
	"apigateway/internal/middleware"
	"net/http"

	menuPB "apigateway/proto/menu"
	"github.com/gin-gonic/gin"
)

func InitMenuRoutes(r *gin.Engine, client menuPB.MenuServiceClient) {
	protected := r.Group("/menu")
	protected.Use(middleware.JWTAuthMiddleware())

	protected.GET("", func(c *gin.Context) {
		res, err := client.ListMenuItems(c, &menuPB.ListMenuItemsRequest{
			Limit: 10, Skip: 0,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, res.Items)
	})

	protected.POST("", func(c *gin.Context) {
		var req menuPB.CreateMenuItemRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		res, err := client.CreateMenuItem(c, &req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"id": res.Id})
	})

	protected.GET("/:id", func(c *gin.Context) {
		id := c.Param("id")
		res, err := client.GetMenuItemByID(c, &menuPB.GetMenuItemByIDRequest{Id: id})
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, res.Item)
	})
	protected.POST("/search", func(c *gin.Context) {
		var req struct {
			Limit    int64  `json:"limit"`
			Skip     int64  `json:"skip"`
			Search   string `json:"search"`
			Category string `json:"category"`
			SortBy   string `json:"sort_by"`
			SortAsc  bool   `json:"sort_asc"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request: " + err.Error()})
			return
		}

		res, err := client.ListMenuItems(c, &menuPB.ListMenuItemsRequest{
			Limit:    req.Limit,
			Skip:     req.Skip,
			Search:   req.Search,
			Category: req.Category,
			SortBy:   req.SortBy,
			SortAsc:  req.SortAsc,
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"total_count": res.TotalCount,
			"items":       res.Items,
		})
	})

	protected.PATCH("/:id", func(c *gin.Context) {
		id := c.Param("id")
		var req menuPB.UpdateMenuItemRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		req.Id = id
		res, err := client.UpdateMenuItem(c, &req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": res.Message})
	})

	protected.DELETE("/:id", func(c *gin.Context) {
		id := c.Param("id")
		res, err := client.DeleteMenuItem(c, &menuPB.DeleteMenuItemRequest{Id: id})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": res.Message})
	})

	protected.POST("/multiple", func(c *gin.Context) {
		var req menuPB.GetMultipleMenuItemsRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		res, err := client.GetMultipleMenuItems(c, &req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, res.Items)
	})
}
