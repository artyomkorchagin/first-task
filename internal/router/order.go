package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) readOrder(c *gin.Context) error {
	return nil
}

func (h *Handler) createOrder(c *gin.Context) error {
	return nil
}

func (h *Handler) renderIndex(c *gin.Context) error {
	c.HTML(http.StatusOK, "index.html", gin.H{})
	return nil
}
