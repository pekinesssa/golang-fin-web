package main

import (
	"context"
	"fmt"
	"golang-fin-web/lib/models"
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
			fmt.Printf(" %s: NOT OK %v\n", name, err)
		} else {
			fmt.Printf(" %s: OK \n", name)
		}
	}
	fmt.Println()

	fmt.Println(" Fetching crypto prices...")

	cryptoTickers := []string{"BTC", "ETH", "SOL", "ADA", "DOT"}

	for _, ticker := range cryptoTickers {
		price, err := marketDataService.GetCryptoPrice(ctx, ticker)
		if err != nil {
			fmt.Printf("  NOT OK %s: %v\n", ticker, err)
			continue
		}

		printCryptoPrice(price)
		time.Sleep(2 * time.Second)
	}

	fmt.Println("Fetching Stock Prices...")

	stockTickers := []string{"AAPL", "MSFT", "GOOGL", "TSLA", "NVDA"}

	for _, ticker := range stockTickers {
		price, err := marketDataService.GetStockPrice(ctx, ticker)
		if err != nil {
			fmt.Printf("  NOT OK %s: %v\n", ticker, err)
			continue
		}

		printStockPrice(price)
		time.Sleep(2 * time.Second)
	}

}

func formatPrice(price float64) string {
	if price >= 1000 {
		return fmt.Sprintf("%.2f", price)
	} else if price >= 1 {
		return fmt.Sprintf("%.4f", price)
	} else {
		return fmt.Sprintf("%.6f", price)
	}
}

func formatLargeNumber(num float64) string {
	if num >= 1000000000000 {
		return fmt.Sprintf("%.2fT", num/1000000000000)
	} else if num >= 1000000000 {
		return fmt.Sprintf("%.2fB", num/1000000000)
	} else if num >= 1000000 {
		return fmt.Sprintf("%.2fM", num/1000000)
	} else if num >= 1000 {
		return fmt.Sprintf("%.2fK", num/1000)
	}
	return fmt.Sprintf("%.2f", num)
}

func printCryptoPrice(p *models.Price) {
	changeDirection := "UP"
	if p.ChangePct24 < 0 {
		changeDirection = "DOWN"
	}

	fmt.Printf(" %s %s\n", changeDirection, p.Symbol)
	fmt.Printf("	Price:		$%s\n", formatPrice(p.Price))

	if p.ChangePct24 != 0 {
		fmt.Printf("     24h Change:  %+.2f%%", p.ChangePct24)
		if p.Change24 != 0 {
			fmt.Printf(" (%+.2f)", p.Change24)
		}
	}

	if p.Volume24 != 0 {
		fmt.Printf("	24h Volume:		$%s\n", formatLargeNumber(p.Volume24))
	}

	if p.High24 != 0 && p.Low24 != 0 {
		fmt.Printf("	24h Range:		$%s - $%s\n", formatPrice(p.Low24), formatPrice(p.High24))
	}

	if p.MarketCap != 0 {
		fmt.Printf("	Market Cap:		$%s\n", formatLargeNumber(p.MarketCap))
	}

	if p.Supply != 0 {
		fmt.Printf("	Supply:		%s\n", formatLargeNumber(p.Supply))
	}

	fmt.Printf("	Source:		%s\n", p.Source)
}

func printStockPrice(p *models.Price) {
	changeDirection := "UP"
	if p.ChangePct24 < 0 {
		changeDirection = "DOWN"
	}

	fmt.Printf("  %s %s\n", changeDirection, p.Symbol)
	fmt.Printf("	Price:	$%s\n", formatPrice(p.Price))

	if p.Open != 0 {
		fmt.Printf("	Open:	$%s\n", formatPrice(p.Open))
	}

	if p.PreviousClose != 0 {
		fmt.Printf("	Prev Close:	$%s\n", formatPrice(p.PreviousClose))
	}

	if p.ChangePct24 != 0 {
		fmt.Printf(" 	Change:	%+.2f%%", p.ChangePct24)
		if p.Change24 != 0 {
			fmt.Printf(" (%+.2f)", p.Change24)
		}
	}

	if p.High24 != 0 && p.Low24 != 0 {
		fmt.Printf("	Day Range:	$%s - $%s\n", formatPrice(p.Low24), formatPrice(p.High24))
	}

	if p.Volume24 != 0 {
		fmt.Printf("	Volume:	%s shares\n", formatLargeNumber(p.Volume24))
	}

	if p.MarketCap != 0 {
		fmt.Printf("	Market Cap:	$%s\n", formatLargeNumber(p.MarketCap))
	}

	fmt.Printf("	 Source:	%s\n", p.Source)

}
