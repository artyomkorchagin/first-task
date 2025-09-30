package orderservice

import (
	"context"
	"fmt"

	"github.com/artyomkorchagin/first-task/internal/types"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (s *Service) CreateOrder(ctx context.Context, o *types.Order) error {
	logger := s.logger.With(zap.String("wallet_uuid", o.OrderUUID))
	if o.OrderUUID == "" {
		return types.ErrBadRequest(fmt.Errorf("orderUUID is empty"))
	}

	_, err := uuid.Parse(o.OrderUUID)
	if err != nil {
		logger.Warn("CreateOrder: invalid UUID",
			zap.Error(err))
		return types.ErrBadRequest(fmt.Errorf("orderUUID is not valid: %w", err))
	}

	if err := s.repo.CreateOrder(ctx, o); err != nil {
		logger.Error("CreateOrder: failed to create order",
			zap.Error(err))
		return err
	}
	logger.Info("CreateOrder: successfully created order", zap.String("orderUUID", o.OrderUUID))

	return nil
}
