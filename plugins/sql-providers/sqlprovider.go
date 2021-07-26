package providers

import (
	"ele/config"
)

// SQLProvider ...
type SQLProvider interface {
	CreateDB(options CreateOptions)
	CanConnect(database *config.DatabaseConfig) (bool, error)
	ToConnectionString(database *config.DatabaseConfig) string
}

// CreateOptions ...
type CreateOptions struct {
	Name     string
	Database *config.DatabaseConfig
}

var sqlproviders map[string]*SQLProvider = make(map[string]*SQLProvider, 0)

// Get ...
func Get(driverName string) (SQLProvider, bool) {
	provider, ok := sqlproviders[driverName]

	return *provider, ok
}

// Register ...
func Register(driverName string, provider SQLProvider) {
	sqlproviders[driverName] = &provider
}
