package orderservice

import (
	"testing"
)

func TestCreateOrder(t *testing.T) {
	tests := []struct {
		name    string
		uuid    string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid",
			uuid:    "c1d5628b-4482-4307-8f7c-85e43bd0f6f7",
			wantErr: false,
		},
		{
			name:    "empty",
			uuid:    "",
			wantErr: true,
			errMsg:  "orderUUID is empty",
		},
		{
			name:    "invalid",
			uuid:    "wss-22r-fdf",
			wantErr: true,
			errMsg:  "orderUUID is invalid",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

		})
	}
}

// func (s *Service) CreateOrder(ctx context.Context, o *types.Order) error {
// 	logger := s.logger.With(zap.String("order_uuid", o.OrderUUID))
// 	if o.OrderUUID == "" {
// 		return types.ErrBadRequest(fmt.Errorf("orderUUID is empty"))
// 	}

// 	_, err := uuid.Parse(o.OrderUUID)
// 	if err != nil {
// 		logger.Warn("CreateOrder: invalid UUID",
// 			zap.Error(err))
// 		return types.ErrBadRequest(fmt.Errorf("orderUUID is not valid: %w", err))
// 	}

// 	if err := s.repo.CreateOrder(ctx, o); err != nil {
// 		logger.Error("CreateOrder: failed to create order",
// 			zap.Error(err))
// 		return err
// 	}
// 	logger.Info("CreateOrder: successfully created order", zap.String("orderUUID", o.OrderUUID))

// 	return nil
// }
