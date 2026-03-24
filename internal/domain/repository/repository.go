package repository

import (
	"context"

	"golang-fin-web/internal/domain"
)

type ConnectionRepository interface {
	CreateConnection(ctx context.Context, conn *domain.Connection) error
	FindByID(ctx context.Context, connectionID string) (*domain.Connection, error)
	GetConnectionsByUserID(ctx context.Context, userID string) ([]*domain.Connection, error)
	UpdateConnectionStatus(ctx context.Context, connectionID string, status domain.ConnectionStatus) error
	DeleteConnection(ctx context.Context, connectionID string) error
}
