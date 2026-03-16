package adapters

import (
	"context"
	"encoding/json"
	"fmt"
	"golang-fin-web/lib/models"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/tidwall/gjson"
	"golang.org/x/time/rate"
)

type ConfigurableAdapter struct {
	*BaseAdapter
	baseURL         string
	endpoints       map[string]string
	authType        string
	authConfig      map[string]string
	tickerFormat    string
	tickerMapping   map[string]string
	responseMapping ResponseMapping
	requestParams   map[string]string
}

type ResponseMapping struct {
	PricePath    string
	Volume24Path string
	Change24Path string
	High24Path   string
	Low24Path    string
	JSONPath     bool
}

func NewConfigurableAdapter(config models.ProviderConfig) (*ConfigurableAdapter, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	baseURL := config.GetConfigString("base_url")
	if baseURL == "" {
		return nil, fmt.Errorf("base_url is required in provider config")
	}

	providerType := models.ProviderType(config.Type)
	base := NewBaseAdapter(config.Name, providerType, config.Priority)
	adapter := &ConfigurableAdapter{
		BaseAdapter:     base,
		baseURL:         baseURL,
		endpoints:       parceEndpoints(config.GetConfigMap("endpoints")),
		authType:        parceAuthType(config.GetConfigString("auth_type")),
		authConfig:      parceAuthConfig(config.GetConfigMap("auth_config")),
		tickerFormat:    config.GetConfigString("ticker_format"),
		tickerMapping:   parceTickerMapping(config.GetConfigMap("ticker_mapping")),
		responseMapping: parceResponceMapping(config.GetConfigMap("response_mapping")),
		requestParams:   parceRequestParams(config.GetConfigMap("request_params")),
	}

	adapter.SetBaseURL(baseURL)
	rateLimitConfig := config.GetConfigMap("rate_limit")
	if limiter := createRateLimiter(rateLimitConfig); limiter != nil {
		adapter.rateLimiter = limiter
	}
	return adapter, nil

}

func parceEndpoints(config map[string]interface{}) map[string]string {
	endpoints := make(map[string]string)
	for key, value := range config {
		if strVal, ok := value.(string); ok {
			endpoints[key] = strVal
		}
	}
	return endpoints
}

func parceTickerMapping(config map[string]interface{}) map[string]string {
	mapping := make(map[string]string)
	for key, value := range config {
		if strVal, ok := value.(string); ok {
			mapping[key] = strVal
		}
	}
	return mapping
}

func parceResponceMapping(config map[string]interface{}) ResponseMapping {
	rm := ResponseMapping{}
	if val, ok := config["price_path"].(string); ok {
		rm.PricePath = val
	}
	if val, ok := config["volume24_path"].(string); ok {
		rm.Volume24Path = val
	}
	if val, ok := config["change24_path"].(string); ok {
		rm.Change24Path = val
	}
	if val, ok := config["high24_path"].(string); ok {
		rm.High24Path = val
	}
	if val, ok := config["low24_path"].(string); ok {
		rm.Low24Path = val
	}
	if val, ok := config["json_path"].(bool); ok {
		rm.JSONPath = val
	}

	if !rm.JSONPath {
		if val, ok := config["json_path"].(bool); ok {
			rm.JSONPath = val
		}
		if val, ok := config["nested"].(bool); ok {
			rm.JSONPath = val
		}
	}
	return rm
}

func parceAuthType(authType string) string {
	switch authType {
	case "api_key":
		return "api_key"
	default:
		return "none"
	}
}

func parceAuthConfig(auth map[string]interface{}) map[string]string {
	config := make(map[string]string)
	for key, value := range auth {
		if strVal, ok := value.(string); ok {
			config[key] = strVal
		}
	}
	return config
}

func parceRequestParams(config map[string]interface{}) map[string]string {
	params := make(map[string]string)
	for key, value := range config {
		if strVal, ok := value.(string); ok {
			params[key] = strVal
		}
	}
	return params
}

func createRateLimiter(config map[string]interface{}) *rate.Limiter {
	if rps, ok := config["requests_per_second"]; ok {
		if rpsInt, ok := rps.(int); ok && rpsInt > 0 {
			return rate.NewLimiter(rate.Limit(rpsInt), rpsInt)
		}
		if rpsFloat, ok := rps.(float64); ok && rpsFloat > 0 {
			rpsInt := int(rpsFloat)
			return rate.NewLimiter(rate.Limit(rpsInt), rpsInt)
		}
	}

	if rpm, ok := config["requests_per_minute"]; ok {
		if rpmInt, ok := rpm.(int); ok && rpmInt > 0 {
			rps := float64(rpmInt) / 60.0
			return rate.NewLimiter(rate.Limit(rps), rpmInt)
		}
		if rpmFloat, ok := rpm.(float64); ok && rpmFloat > 0 {
			rps := rpmFloat / 60.0
			rpmInt := int(rpmFloat)
			return rate.NewLimiter(rate.Limit(rps), rpmInt)
		}
	}

	return nil
}

func (a *ConfigurableAdapter) preparePath(path string, ticker string) string {
	path = strings.ReplaceAll(path, "{ticker}", ticker)
	path = strings.ReplaceAll(path, "{ticker_lower}", strings.ToLower(ticker))
	path = strings.ReplaceAll(path, "{ticker_upper}", strings.ToUpper(ticker))
	return path
}

func (a *ConfigurableAdapter) formatTicker(ticker string) string {
	if mapped, ok := a.tickerMapping[ticker]; ok {
		return mapped
	}
	if a.tickerFormat == "" {
		return ticker
	}
	formatted := strings.ReplaceAll(a.tickerFormat, "{ticker}", ticker)
	formatted = strings.ReplaceAll(formatted, "{ticker_upper}", strings.ToUpper(ticker))
	formatted = strings.ReplaceAll(formatted, "{ticker_lower}", strings.ToLower(ticker))
	return formatted
}

func convertToFloat64(val interface{}) float64 {
	switch v := val.(type) {
	case float64:
		return v
	case float32:
		return float64(v)
	case int:
		return float64(v)
	case int64:
		return float64(v)
	case string:
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			return f
		}
	}
	return 0
}

func (a *ConfigurableAdapter) GetPrice(ctx context.Context, ticker string) (*models.Price, error) {
	apiTicker := a.formatTicker(ticker)
	endpoint := a.endpoints["price"]
	if endpoint == "" {
		return nil, fmt.Errorf("price endpoint not configured for provider %s", a.Name())
	}

	endpoint = strings.ReplaceAll(endpoint, "{ticker}", apiTicker)
	endpoint = strings.ReplaceAll(endpoint, "{ticker_lower}", strings.ToLower(apiTicker))

	url := a.baseURL + endpoint

	resp, err := a.MakeRequest(ctx, http.MethodGet, url)
	if err != nil {
		return nil, fmt.Errorf("error making request to %s: %w", url, err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non-200 response from %s: %d - %s", url, resp.StatusCode, string(body))
	}

	price, err := a.parsePrice(body, apiTicker, ticker)
	if err != nil {
		return nil, fmt.Errorf("error parsing price from %s: %w", url, err)
	}

	price.Source = a.Name()
	price.Timestamp = time.Now()

	return price, nil
}

func (a *ConfigurableAdapter) parsePrice(data []byte, apiTicker string, originalSymbol string) (*models.Price, error) {
	price := &models.Price{
		Symbol: originalSymbol,
	}

	if a.responseMapping.JSONPath {
		priceField := a.preparePath(a.responseMapping.PricePath, apiTicker)
		priceValue := gjson.GetBytes(data, priceField)
		if !priceValue.Exists() {
			return nil, fmt.Errorf("price field not found in response using path: %s", priceField)
		}
		price.Price = priceValue.Float()

		if a.responseMapping.Volume24Path != "" {
			volume24Field := a.preparePath(a.responseMapping.Volume24Path, apiTicker)
			volume24Value := gjson.GetBytes(data, volume24Field)
			if volume24Value.Exists() {
				price.Volume24 = volume24Value.Float()
			}
		}

		if a.responseMapping.Change24Path != "" {
			change24Field := a.preparePath(a.responseMapping.Change24Path, apiTicker)
			change24Value := gjson.GetBytes(data, change24Field)
			if change24Value.Exists() {
				price.Change24 = change24Value.Float()
			}
		}

		if a.responseMapping.High24Path != "" {
			high24Field := a.preparePath(a.responseMapping.High24Path, apiTicker)
			high24Value := gjson.GetBytes(data, high24Field)
			if high24Value.Exists() {
				price.High24 = high24Value.Float()
			}
		}

		if a.responseMapping.Low24Path != "" {
			low24Field := a.preparePath(a.responseMapping.Low24Path, apiTicker)
			low24Value := gjson.GetBytes(data, low24Field)
			if low24Value.Exists() {
				price.Low24 = low24Value.Float()
			}
		}
	} else {
		var resultMap map[string]interface{}
		if err := json.Unmarshal(data, &resultMap); err != nil {
			return nil, fmt.Errorf("error unmarshaling response: %w", err)
		}

		if val, ok := resultMap[a.responseMapping.PricePath]; ok {
			price.Price = convertToFloat64(val)
		} else {
			return nil, fmt.Errorf("price field not found in response using key: %s", a.responseMapping.PricePath)
		}

		if a.responseMapping.Volume24Path != "" {
			if val, ok := resultMap[a.responseMapping.Volume24Path]; ok {
				price.Volume24 = convertToFloat64(val)
			}
		}

		if a.responseMapping.Change24Path != "" {
			if val, ok := resultMap[a.responseMapping.Change24Path]; ok {
				price.Change24 = convertToFloat64(val)
			}
		}
	}

	return price, nil
}

func (a *ConfigurableAdapter) GetMultiplePrices(ctx context.Context, tickers []string) (map[string]*models.Price, error) {
	prices := make(map[string]*models.Price)
	for _, ticker := range tickers {
		price, err := a.GetPrice(ctx, ticker)
		if err != nil {
			fmt.Errorf("error getting price for %s: %w", ticker, err)
			continue
		}
		prices[ticker] = price
	}

	return prices, nil
}
