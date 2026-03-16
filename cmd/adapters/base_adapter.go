package adapters

import (
	"context"
	"golang-fin-web/lib/models"
	"net/http"
	"time"

	"golang.org/x/time/rate"
)

// Базовый адаптер, который будет использоваться для всех провайдеров, и который будет реализовывать общую логику для получения данных и предоставления их в нужном формате для приложения
type BaseAdapter struct {
	name         string
	providerType models.ProviderType
	baseURL      string
	priority     int
	rateLimiter  *rate.Limiter
	ticker       string
	httpClient   http.Client
}

func NewBaseAdapter(name string, providerType models.ProviderType, priority int) *BaseAdapter {
	return &BaseAdapter{
		name:         name,
		providerType: providerType,
		priority:     priority,
		httpClient: http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (b *BaseAdapter) Name() string {
	return b.name
}

func (b *BaseAdapter) Type() models.ProviderType {
	return b.providerType
}

func (b *BaseAdapter) Priority() int {
	return b.priority
}

func (b *BaseAdapter) MakeRequest(ctx context.Context, method, url string) (*http.Response, error) {
	var resp *http.Response
	var err error

	if b.rateLimiter != nil {
		err = b.rateLimiter.Wait(ctx)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err = b.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (b *BaseAdapter) SetBaseURL(url string) {
	b.baseURL = url
}

func (b *BaseAdapter) SetRateLimiter(limiter *rate.Limiter) {
	b.rateLimiter = limiter
}

func (b *BaseAdapter) Ping(ctx context.Context) error {
	if b.baseURL == "" {
		return nil // Если нет базового URL, просто считаем, что провайдер доступен (например, для провайдеров, которые не требуют HTTP)
	}

	resp, err := b.MakeRequest(ctx, "GET", b.baseURL)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}
