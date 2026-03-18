package models

import (
	"context"
	"time"
)

type ProviderType string

const (
	ProviderTypeStock  = "STOCK"
	ProviderTypeCrypto = "CRYPTO"
	ProviderTypeBoth   = "BOTH"
	ЗProviderTypeForex = "FOREX"
)

type AssetType string

const (
	AssetTypeStock  = "STOCK"
	AssetTypeCrypto = "CRYPTO"
	AssetTypeETF    = "ETF"
	AssetTypeForex  = "FOREX"
)

// Интерфейс для провайдера финансовых данных (делаем универсальный модуль, а потом адаптеры к АПИ)
type Provider interface {
	Name() string
	Type() ProviderType
	Priority() int

	GetPrice(ctx context.Context, ticker string) (*Price, error)
	GetMultiplePrices(ctx context.Context, tickers []string) (map[string]*Price, error)

	Ping(ctx context.Context) error
}

type AssetInfo struct {
	Ticker      string
	Name        string
	Exchange    string
	Description string
	Type        string
}

type Price struct {
	Symbol string
	Price  float64

	Volume24    float64
	Change24    float64
	ChangePct24 float64
	High24      float64
	Low24       float64

	Open          float64
	PreviousClose float64
	MarketCap     float64
	Supply        float64 // Для криптовалют

	Timestamp time.Time
	Source    string
}
