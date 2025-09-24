package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) readOrder(c *gin.Context) error {
	orderID := c.Param("id")

	order, err := h.orderService.ReadOrder(c, orderID)
	if err != nil {
		return err
	}
	c.JSON(http.StatusOK, order)

	return nil
}

func (h *Handler) createOrder(c *gin.Context) error {

	return nil
}

func (h *Handler) renderIndex(c *gin.Context) error {
	c.HTML(http.StatusOK, "index.html", gin.H{})
	return nil
}
