package models

import (
	"context"
	"time"
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

	Ping(ctx context.Context) error
}

type AssetInfo struct {
	Ticker      string
	Name        string
	Exchange    string
	Description string
	Type        string
}

type PriceChange struct {
	Symbol    string
	Price     float64
	Volume24  float64
	Change24  float64
	High24    float64
	Low24     float64
	Timestamp time.Time
	Source    string
}
