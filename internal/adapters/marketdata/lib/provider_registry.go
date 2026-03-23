package factory

import (
	"context"
	"fmt"
	"golang-fin-web/internal/adapters/marketdata/models"
	"log"
	"os"
	"sort"

	"gopkg.in/yaml.v3"
)

type ProviderRegistry struct {
	cryptoProviders []ProviderWithPriority
	stockProviders  []ProviderWithPriority
	allProviders    map[string]models.Provider

	factory *ProviderFactory
}

type ProviderWithPriority struct {
	Provider models.Provider
	Priority int
}

func NewRegistryFromConfig(configPath string) (*ProviderRegistry, error) {

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var config struct {
		Providers []models.ProviderConfig `yaml:"providers"`
	}

	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	registry := &ProviderRegistry{
		cryptoProviders: []ProviderWithPriority{},
		stockProviders:  []ProviderWithPriority{},
		allProviders:    make(map[string]models.Provider),
		factory:         NewProviderFactory(),
	}

	for _, pc := range config.Providers {
		if !pc.Enabled {
			log.Printf("Skipping disabled provider: %s", pc.Name)
			continue
		}

		provider, err := registry.factory.CreateProvider(pc)
		if err != nil {
			log.Printf("Failed to create provider %s: %v", pc.Name, err)
			continue
		}

		log.Printf("Created provider: %s (type=%s, priority=%d)",
			pc.Name, provider.Type(), pc.Priority)

		pwp := ProviderWithPriority{
			Provider: provider,
			Priority: pc.Priority,
		}

		switch provider.Type() {
		case models.ProviderTypeCrypto:
			registry.cryptoProviders = append(registry.cryptoProviders, pwp)
		case models.ProviderTypeStock:
			registry.stockProviders = append(registry.stockProviders, pwp)
		case models.ProviderTypeBoth:
			registry.cryptoProviders = append(registry.cryptoProviders, pwp)
			registry.stockProviders = append(registry.stockProviders, pwp)
		}

		registry.allProviders[pc.Name] = provider
	}

	sort.Slice(registry.cryptoProviders, func(i, j int) bool {
		return registry.cryptoProviders[i].Priority < registry.cryptoProviders[j].Priority
	})

	sort.Slice(registry.stockProviders, func(i, j int) bool {
		return registry.stockProviders[i].Priority < registry.stockProviders[j].Priority
	})

	log.Printf("Registry initialized: %d crypto providers, %d stock providers",
		len(registry.cryptoProviders), len(registry.stockProviders))

	return registry, nil
}

func (r *ProviderRegistry) GetPrice(ctx context.Context, ticker string, assetType models.AssetType) (*models.Price, error) {
	var providers []ProviderWithPriority

	switch assetType {
	case models.AssetTypeCrypto:
		providers = r.cryptoProviders
	case models.AssetTypeStock, models.AssetTypeETF:
		providers = r.stockProviders
	default:
		return nil, fmt.Errorf("unsupported asset type: %s", assetType)
	}

	if len(providers) == 0 {
		return nil, fmt.Errorf("no providers available for type %s", assetType)
	}

	var lastErr error
	for _, p := range providers {
		log.Printf("Trying provider: %s (priority=%d) for %s",
			p.Provider.Name(), p.Priority, ticker)

		price, err := p.Provider.GetPrice(ctx, ticker)
		if err == nil {
			log.Printf("Got price for %s from %s: $%.2f",
				ticker, p.Provider.Name(), price.Price)
			return price, nil
		}

		log.Printf("Provider %s failed: %v", p.Provider.Name(), err)
		lastErr = err
	}

	return nil, fmt.Errorf("all providers failed for %s: %w", ticker, lastErr)
}

func (r *ProviderRegistry) Get(name string) (models.Provider, error) {
	provider, ok := r.allProviders[name]
	if !ok {
		return nil, fmt.Errorf("provider %s not found", name)
	}
	return provider, nil
}

func (r *ProviderRegistry) ListProviders() []string {
	names := make([]string, 0, len(r.allProviders))
	for name := range r.allProviders {
		names = append(names, name)
	}
	return names
}

func (r *ProviderRegistry) HealthCheck(ctx context.Context) map[string]error {
	results := make(map[string]error)

	for name, provider := range r.allProviders {
		err := provider.Ping(ctx)
		results[name] = err

		if err != nil {
			log.Printf("Health check failed for %s: %v", name, err)
		} else {
			log.Printf("Health check OK for %s", name)
		}
	}

	return results
}
