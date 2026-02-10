package models

import (
	"net/http"

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
	name         string
	providerType ProviderType
	baseURL      string
	endpoints    map[string]string
	requestLimit *rate.Limiter
	ticker       string
	httpClient   http.Client
}
