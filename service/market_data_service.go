// service/market_data_service.go
package service

import (
	"context"
	"fmt"
	factory "golang-fin-web/lib"
	"golang-fin-web/lib/models"
)

type MarketDataService struct {
	registry *factory.ProviderRegistry
}

func NewMarketDataService(configPath string) (*MarketDataService, error) {
	registry, err := factory.NewRegistryFromConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create provider registry: %w", err)
	}

	return &MarketDataService{
		registry: registry,
	}, nil
}

func (s *MarketDataService) GetPrice(ctx context.Context, ticker string, assetType models.AssetType) (*models.Price, error) {
	return s.registry.GetPrice(ctx, ticker, assetType)
}

func (s *MarketDataService) GetCryptoPrice(ctx context.Context, ticker string) (*models.Price, error) {
	return s.GetPrice(ctx, ticker, models.AssetTypeCrypto)
}

func (s *MarketDataService) GetStockPrice(ctx context.Context, ticker string) (*models.Price, error) {
	return s.GetPrice(ctx, ticker, models.AssetTypeStock)
}

func (s *MarketDataService) ListProviders() []string {
	return s.registry.ListProviders()
}

func (s *MarketDataService) HealthCheck(ctx context.Context) map[string]error {
	return s.registry.HealthCheck(ctx)
}
