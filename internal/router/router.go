package router

import (
	"net/http"

	orderservice "github.com/artyomkorchagin/first-task/internal/service"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
)

type Handler struct {
	orderService *orderservice.Service
	logger       *zap.Logger
}

func NewHandler(orderService *orderservice.Service, logger *zap.Logger) *Handler {
	return &Handler{
		orderService: orderService,
		logger:       logger,
	}
}

func (h *Handler) InitRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.LoadHTMLGlob("static/html/*")
	router.Static("/static", "./static/")
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	main := router.Group("/")
	{
		main.GET("/", h.wrap(h.renderIndex))
		main.GET("/order", h.wrap(h.readOrder))

		main.GET("/status", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		})

		main.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}
	h.logger.Info("Routes initialized")
	return router
}
