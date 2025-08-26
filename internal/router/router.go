package router

import (
	"log"
	"net/http"

	"github.com/artyomkorchagin/first-task/internal/middleware"
	orderservice "github.com/artyomkorchagin/first-task/internal/service"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handler struct {
	orderService *orderservice.Service
}

func NewHandler(orderService *orderservice.Service) *Handler {
	return &Handler{
		orderService: orderService,
	}
}

func (h *Handler) InitRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Static("/static", "../../static/")
	router.Use(middleware.LoggerMiddleware())

	main := router.Group("/")
	{
		main.GET("/", h.wrap(h.renderIndex))
		main.GET("/orders/:id", h.wrap(h.getOrder))

		main.GET("/status", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		})

		main.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}
	log.Println("Routes initialized")
	return router
}
