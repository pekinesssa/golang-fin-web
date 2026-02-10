package models

import (
	"context"
)

type ProviderType string

const (
	ProviderTypeStock = "STOCK"
)

// Интерфейс для провайдера финансовых данных (делаем универсальный модуль, а потом адаптеры к АПИ)
type Provider interface {
	Name() string
	Type() ProviderType

	GetPrice(ctx context.Context, ticker string) (float64, error)
	GetAssetInfo(ctx context.Context, ticker string) (*AssetInfo, error)

	HealthCheck(ctx context.Context) error
}

type AssetInfo struct {
	Ticker      string
	Name        string
	Exchange    string
	Description string
	Type        string
}
