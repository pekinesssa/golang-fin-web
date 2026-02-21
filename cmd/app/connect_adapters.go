package app

import (
	"context"
	"golang-fin-web/lib/models"
	"net/http"
	"strings"
	"time"

	"golang.org/x/time/rate"
)

// Конфигурация для адаптера, который будет использовать провайдера
type ProviderConfig struct {
	Name     string                 `yaml:"name"`
	Type     string                 `yaml:"type"`
	Enabled  bool                   `yaml:"enabled"`
	Priority int                    `yaml:"priority"`
	Config   map[string]interface{} `yaml:"config"`
}

// Модель адаптера, который будет использовать провайдера для получения данных и предоставлять их в нужном формате для приложения
type AdapterModel struct {
	Name         string
	ProviderType models.ProviderType
	BaseURL      string
	Endpoints    map[string]string
	RateLimiter  *rate.Limiter
	Ticker       string
	HttpClient   http.Client
}

type ResponceMap struct {
	Price  string
	Symbol string
	Volume string
}

func NewAdapter(config ProviderConfig) (*AdapterModel, error) {
	cfg := config.Config

	endpoints := make(map[string]string)
	if eps, ok := cfg["endpoints"].(map[string]interface{}); ok {
		for i, j := range eps {
			endpoints[i] = j.(string)
		}
	}

	respMap := ResponceMap{}
	if rm, ok := cfg["responce_map"].(map[string]interface{}); ok {
		respMap.Price = rm["price"].(string)
		if i, ok := rm["symbol"].(string); ok {
			respMap.Symbol = i
		}
		if i, ok := rm["volume"].(string); ok {
			respMap.Symbol = i
		}
	}

	var limiter *rate.Limiter
	if rl, ok := cfg["rate_limit"].(map[string]interface{}); ok {
		if rps, ok := rl["requests_per_second"].(int); ok {
			limiter = rate.NewLimiter(rate.Limit(rps), rps)
		} else if rpm, ok := rl["requests_per_minute"].(int); ok {
			limiter = rate.NewLimiter(rate.Limit(rpm)/60, rpm)
		}
	}

	symbolMap := make(map[string]string)
	if sm, ok := cfg["symbol_mapping"].(map[string]interface{}); ok {
		for i, j := range sm {
			symbolMap[i] = j.(string)
		}
	}

	return &AdapterModel{
		Name:         config.Name,
		RateLimiter:  limiter,
		ProviderType: models.ProviderType(strings.ToUpper(config.Type)),
		BaseURL:      cfg["base_url"].(string),
		Endpoints:    endpoints,
		Ticker:       cfg["ticker"].(string),
		HttpClient: http.Client{
			Timeout: 10 * time.Second,
		},
	}, nil
}

func (a *AdapterModel) GetName() string {
	return a.Name
}

func (a *AdapterModel) GetProviderType() models.ProviderType {
	return a.ProviderType
}

func (a *AdapterModel) GetPrice() (ctx context.Context, price float64, err error) {
	// Здесь будет логика для получения данных от провайдера, например, через HTTP запросы к API, с учетом rate limiter и обработки ошибок
	return ctx, 0, nil
}
