package blockchain

import (
	"context"
	"golang-fin-web/internal/domain"
)

type BlockchainAdapter interface {
	Name() string
	SupportedNetworks() []string
	GetNativeBalances(ctx context.Context, address, network string) (*domain.CryptoBalance, error)
	GetTokenBalances(ctx context.Context, address, network string) ([]*domain.CryptoBalance, error)
}
