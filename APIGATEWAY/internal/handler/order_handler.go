package handler

import (
	"apigateway/internal/middleware"
	"net/http"

	orderPB "apigateway/proto/order"
	"github.com/gin-gonic/gin"
)

func InitOrderRoutes(r *gin.Engine, client orderPB.OrderServiceClient) {
	protected := r.Group("/orders")
	protected.Use(middleware.JWTAuthMiddleware())

	protected.POST("", func(c *gin.Context) {
		var req orderPB.CreateOrderRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		res, err := client.CreateOrder(c, &req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"order_id": res.Id})
	})

	protected.GET("/:id", func(c *gin.Context) {
		id := c.Param("id")
		res, err := client.GetOrder(c, &orderPB.GetOrderRequest{Id: id})
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, res.Order)
	})

	protected.GET("", func(c *gin.Context) {
		res, err := client.ListOrders(c, &orderPB.ListOrdersRequest{
			Limit: 10,
			Skip:  0,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, res.Orders)
	})
	protected.GET("/user/:userId", func(c *gin.Context) {
		userId := c.Param("userId")

		res, err := client.ListOrdersByUser(c, &orderPB.ListOrdersByUserRequest{
			UserId: userId,
			Limit:  100,
			Skip:   0,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, res.Orders)
	})
	protected.PUT("/:id", func(c *gin.Context) {
		id := c.Param("id")
		var req orderPB.UpdateOrderRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		req.Id = id
		res, err := client.UpdateOrder(c, &req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": res.Message})
	})

	protected.PATCH("/:id/status", func(c *gin.Context) {
		id := c.Param("id")
		var req orderPB.PatchOrderStatusRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		req.Id = id
		res, err := client.PatchOrderStatus(c, &req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": res.Message})
	})

	protected.DELETE("/:id", func(c *gin.Context) {
		id := c.Param("id")
		res, err := client.DeleteOrder(c, &orderPB.DeleteOrderRequest{Id: id})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": res.Message})
	})
}
