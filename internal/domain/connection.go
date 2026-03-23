package domain

import (
	"errors"
	"time"
)

var (
	ErrConnectionFailed   = errors.New("connection failed")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrEmptyAddress       = errors.New("empty address")
)

type ConnectionType string

const (
	ConnectionTypeWallet   ConnectionType = "wallet"
	ConnectionTypeExchange ConnectionType = "exchange"
)

type ConnectionStatus string

const (
	StatusActive  ConnectionStatus = "active"
	StatusPending ConnectionStatus = "pending"
	StatusError   ConnectionStatus = "Error"
)

type Connection struct {
	id          string
	userID      string
	portfolioID string
	name        string
	conType     ConnectionType
	status      ConnectionStatus
	provider    string

	apiKey    string
	secretKey string
	address   string

	updatedAt time.Time
	lastError string
}

func NewConnection(id, userID, portfolioID, name string, conType ConnectionType, provider string) (*Connection, error) {
	if id == "" {
		return nil, errors.New("id cannot be empty")
	}
	if userID == "" {
		return nil, errors.New("userID cannot be empty")
	}
	if portfolioID == "" {
		return nil, errors.New("portfolioID cannot be empty")
	}
	if name == "" {
		return nil, errors.New("name cannot be empty")
	}

	now := time.Now()
	return &Connection{
		id:          id,
		userID:      userID,
		portfolioID: portfolioID,
		name:        name,
		conType:     conType,
		status:      StatusPending,
		provider:    provider,
		updatedAt:   now,
	}, nil
}

func (c *Connection) ID() string {
	return c.id
}

func (c *Connection) UserID() string {
	return c.userID
}

func (c *Connection) PortfolioID() string {
	return c.portfolioID
}

func (c *Connection) Name() string {
	return c.name
}

func (c *Connection) Type() ConnectionType {
	return c.conType
}

func (c *Connection) Status() ConnectionStatus {
	return c.status
}

func (c *Connection) Provider() string {
	return c.provider
}

func (c *Connection) Address() string {
	return c.address
}
func (c *Connection) APIKey() string {
	return c.apiKey
}

func (c *Connection) SecretKey() string {
	return c.secretKey
}

func (c *Connection) SetCredentials(apiKey, secretKey string) error {
	if apiKey == "" || secretKey == "" {
		return ErrInvalidCredentials
	}
	if c.conType == ConnectionTypeExchange {
		return errors.New("credentials are only applicable for exchange connections")
	}
	c.apiKey = apiKey
	c.secretKey = secretKey
	c.updatedAt = time.Now()
	return nil
}

func (c *Connection) SetAddress(address string) error {
	if address == "" {
		return ErrEmptyAddress
	}
	if c.conType == ConnectionTypeWallet {
		return errors.New("address is only applicable for wallet connections")
	}
	c.address = address
	c.updatedAt = time.Now()
	return nil
}

func (c *Connection) UpdateStatusActive(status ConnectionStatus) {
	c.status = StatusActive
	c.lastError = ""
	c.updatedAt = time.Now()
}

func (c *Connection) MarkError(err error) {
	c.status = StatusError
	c.lastError = err.Error()
	c.updatedAt = time.Now()
}

func (c *Connection) IsExchange() bool {
	return c.conType == ConnectionTypeExchange
}

func (c *Connection) IsWallet() bool {
	return c.conType == ConnectionTypeWallet
}
