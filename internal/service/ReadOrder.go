package orderservice

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/artyomkorchagin/first-task/internal/types"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (s *Service) ReadOrder(ctx context.Context, orderUUID string) (*types.Order, error) {
	var order *types.Order
	logger := s.logger.With(zap.String("wallet_uuid", orderUUID))
	if orderUUID == "" {
		return nil, types.ErrBadRequest(fmt.Errorf("orderUUID is empty"))
	}

	_, err := uuid.Parse(orderUUID)
	if err != nil {
		logger.Warn("ReadOrder: invalid UUID",
			zap.Error(err))
		return nil, types.ErrBadRequest(fmt.Errorf("orderUUID is not valid: %w", err))
	}

	key := fmt.Sprintf("order:%s", orderUUID)
	if val, err := s.redis.Get(ctx, key).Result(); err == nil {
		if err := json.Unmarshal([]byte(val), &order); err != nil {
			logger.Warn("ReadOrder: unmarshal error")
		}
		return order, nil
	}

	order, err = s.repo.ReadOrder(ctx, orderUUID)
	if err != nil {
		return nil, types.ErrNotFound(err)
	}

	encodedOrder, err := json.Marshal(order)
	if err != nil {
		logger.Warn("ReadOrder: failed to marshal order",
			zap.Error(err))
	} else {
		if err := s.redis.Set(ctx, key, encodedOrder, 15*time.Second).Err(); err != nil {
			logger.Warn("ReadOrder: failed to set cache",
				zap.String("cache_key", key),
				zap.Error(err))
		} else {
			logger.Debug("GetBalance: balance cached",
				zap.ByteString("order", encodedOrder),
				zap.Duration("ttl", 15*time.Second))
		}
	}

	return order, nil
}
