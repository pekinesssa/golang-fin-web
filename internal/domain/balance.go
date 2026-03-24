package domain

import (
	"errors"
)

type AssetType string

const (
	AssetTypeCrypto AssetType = "crypto"
	AssetTypeStock  AssetType = "stock"
	AssetTypeFiat   AssetType = "fiat"
)

type CryptoBalance struct {
	ticker    string
	assetType AssetType
	valueRUB  float64
	valueUSD  float64

	availableAmount float64
	lockedAmount    float64
	totalAmount     float64
	contractAddress string

	network string
}

func NewCryptoBalance(ticker string, assetType AssetType, totalAmount, availableAmount, lockedAmount, valueRUB, valueUSD float64) (*CryptoBalance, error) {
	if ticker == "" {
		return nil, errors.New("ticker cannot be empty")
	}
	if totalAmount < 0 || availableAmount < 0 || lockedAmount < 0 {
		return nil, errors.New("amounts cannot be negative")
	}
	if valueRUB < 0 || valueUSD < 0 {
		return nil, errors.New("values cannot be negative")
	}

	return &CryptoBalance{
		assetType:       assetType,
		valueRUB:        valueRUB,
		valueUSD:        valueUSD,
		availableAmount: availableAmount,
		lockedAmount:    lockedAmount,
		totalAmount:     totalAmount,
	}, nil
}

func (b *CryptoBalance) Ticker() string {
	return b.ticker
}

func (b *CryptoBalance) AssetType() AssetType {
	return b.assetType
}

func (b *CryptoBalance) ValueRUB() float64 {
	return b.valueRUB
}

func (b *CryptoBalance) ValueUSD() float64 {
	return b.valueUSD
}

func (b *CryptoBalance) AvailableAmount() float64 {
	return b.availableAmount
}

func (b *CryptoBalance) LockedAmount() float64 {
	return b.lockedAmount
}

func (b *CryptoBalance) TotalAmount() float64 {
	return b.totalAmount
}

func (b *CryptoBalance) ContractAddress() string {
	return b.contractAddress
}

func (b *CryptoBalance) WithNetwork(network string) *CryptoBalance {
	b.network = network
	return b
}

func (b *CryptoBalance) WithContractAddress(addr string) *CryptoBalance {
	b.contractAddress = addr
	return b
}
