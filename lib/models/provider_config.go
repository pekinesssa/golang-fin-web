package models

import (
	"fmt"
	"os"
	"strings"
)

type ProviderConfig struct {
	Name     string `yaml:"name"`
	Enabled  bool   `yaml:"enabled"`
	Type     string `yaml:"type"`
	Priority int    `yaml:"priority"`

	// Аутентификация
	APIKey    string `yaml:"api_key"`
	APISecret string `yaml:"api_secret"`

	// Детальная конфигурация
	Config map[string]interface{} `yaml:"config"`
}

func (c *ProviderConfig) Validate() error {
	if c.Name == "" {
		return fmt.Errorf("provider name is required")
	}

	if c.Type == "" {
		return fmt.Errorf("provider type is required")
	}

	if c.Priority < 0 {
		return fmt.Errorf("priority must be >= 0")
	}

	return nil
}

func (c *ProviderConfig) GetAPIKey() string {
	return c.expandEnv(c.APIKey)
}

func (c *ProviderConfig) GetAPISecret() string {
	return c.expandEnv(c.APISecret)
}
func (c *ProviderConfig) expandEnv(value string) string {
	if strings.HasPrefix(value, "${") && strings.HasSuffix(value, "}") {
		varName := value[2 : len(value)-1]
		return os.Getenv(varName)
	}
	return value
}

func (c *ProviderConfig) GetConfigString(key string) string {
	if val, ok := c.Config[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

func (c *ProviderConfig) GetConfigInt(key string) int {
	if val, ok := c.Config[key]; ok {
		if num, ok := val.(int); ok {
			return num
		}
	}
	return 0
}

func (c *ProviderConfig) GetConfigMap(key string) map[string]interface{} {
	if val, ok := c.Config[key]; ok {
		if m, ok := val.(map[string]interface{}); ok {
			return m
		}
	}
	return make(map[string]interface{})
}

func (c *ProviderConfig) GetConfigBool(key string) bool {
	if val, ok := c.Config[key]; ok {
		if b, ok := val.(bool); ok {
			return b
		}
	}
	return false
}
