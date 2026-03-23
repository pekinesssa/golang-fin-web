package factory

import (
	adapters "golang-fin-web/internal/adapters/marketdata"
	"golang-fin-web/internal/adapters/marketdata/models"
)

type ProviderFactory struct {
	creators map[string]ProviderCreator
}

type ProviderCreator func(models.ProviderConfig) (models.Provider, error)

func NewProviderFactory() *ProviderFactory {
	factory := &ProviderFactory{
		creators: make(map[string]ProviderCreator),
	}
	//factory.Register("alphavantage", createAlphaVantageProvider)
	// Регистрируем провайдеры
	return factory
}

func (f *ProviderFactory) RegisterProvider(name string, creator ProviderCreator) {
	f.creators[name] = creator
}

func (f *ProviderFactory) CreateProvider(config models.ProviderConfig) (models.Provider, error) {
	if creator, ok := f.creators[config.Name]; ok {
		return creator(config)
	}
	return adapters.NewConfigurableAdapter(config)
}
