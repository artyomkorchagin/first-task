package router

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/artyomkorchagin/first-task/internal/types"
	"github.com/gin-gonic/gin"
)

func (h *Handler) readOrder(c *gin.Context) error {
	orderID := c.Param("order_uid")
	order, err := h.orderService.ReadOrder(c, orderID)
	if err != nil {
		return err
	}

	c.HTML(http.StatusOK, "order.html", gin.H{
		"Order": order,
	})

	return nil
}

func (h *Handler) CreateOrderKafka(ctx context.Context, msg []byte) error {
	order := &types.Order{}
	if err := json.Unmarshal(msg, order); err != nil {
		return err
	}
	if err := h.orderService.CreateOrder(ctx, order); err != nil {
		return err
	}
	return nil
}

func (h *Handler) renderIndex(c *gin.Context) error {
	c.HTML(http.StatusOK, "index.html", nil)
	return nil
}
