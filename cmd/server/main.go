package main

import (
	"context"
	"fmt"
	"golang-fin-web/service"
	"log"
	"time"
)

func main() {
	fmt.Println("Starting Market Data Service...")

	marketDataService, err := service.NewMarketDataService("config/providers.yaml")
	if err != nil {
		log.Fatalf("Failed to create service: %v", err)
	}

	fmt.Println("Service created successfully")

	providers := marketDataService.ListProviders()
	fmt.Printf("Available providers: %v\n\n", providers)

	fmt.Println("Running health check...")
	ctx := context.Background()
	healthResults := marketDataService.HealthCheck(ctx)

	for name, err := range healthResults {
		if err != nil {
			fmt.Printf(" %s: NOT OK%v\n", name, err)
		} else {
			fmt.Printf(" %s: OK\n", name)
		}
	}
	fmt.Println()

	fmt.Println(" Fetching crypto prices...")

	cryptoTickers := []string{"BTC", "ETH", "SOL", "ADA"}

	for _, ticker := range cryptoTickers {
		price, err := marketDataService.GetCryptoPrice(ctx, ticker)
		if err != nil {
			fmt.Printf("  NOT OK %s: %v\n", ticker, err)
			continue
		}

		fmt.Printf("  OK %s: $%.2f", ticker, price.Price)

		if price.Change24 != 0 {
			fmt.Printf(" (24h: %.2f%%)", price.Change24)
		}

		if price.Volume24 != 0 {
			fmt.Printf(" [Vol: $%.0f]", price.Volume24)
		}

		fmt.Printf(" - from %s\n", price.Source)

		time.Sleep(2 * time.Second)
	}
}
