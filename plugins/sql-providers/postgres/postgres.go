package postgres

import (
	"database/sql"
	"ele/config"
	providers "ele/plugins/sql-providers"
	"ele/utils/doc"
	"fmt"
	"log"
	"net"
	"os"
)

// Provider ...
type Provider struct {
}

// CreateDB ...
func (prov *Provider) CreateDB(options providers.CreateOptions) {
	config := options.Database.Clone()

	config.DBName = "postgres" // Change to default db
	connectionstring := prov.ToConnectionString(&config)

	db, err := sql.Open(options.Database.Driver, connectionstring)
	defer db.Close()
	if err != nil {
		log.Fatal(err)
	}

	var count int
	qerr := db.QueryRow("SELECT COUNT(*) RowCount FROM pg_database WHERE datname = $1", options.Name).Scan(&count)
	if qerr != nil {
		log.Fatal(qerr)
	}

	if count == 0 {
		_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", options.Name))

		if err != nil {
			log.Fatal(err)
		}
	}
}

// CanConnect ...
func (prov *Provider) CanConnect(database *config.DatabaseConfig) (bool, error) {
	connectionstring := prov.ToConnectionString(database)
	db, err := sql.Open(database.Driver, connectionstring)
	defer db.Close()
	if err != nil {
		return false, err
	}

	err = db.Ping()

	if _, ok := err.(net.Error); ok {
		// You know it's a net.Error instance
		// Expected to be nil if known error
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return true, nil
}

// NewProvider ...
func NewProvider() *Provider {
	return &Provider{}
}

// ToConnectionString ...
func (prov *Provider) ToConnectionString(database *config.DatabaseConfig) string {
	if database == nil {
		doc.Paragraph(doc.Error("No Database settings provided"))
		os.Exit(1)
	}
	// port=%d host=%s
	cstring := fmt.Sprintf("dbname=%s user=%s password=%s sslmode=%s", database.DBName, database.User, database.Password, database.SSLMode)
	if database.Host != "" {
		cstring += fmt.Sprintf(" host=%s", database.Host)
	}

	if database.Port > 0 {
		cstring += fmt.Sprintf(" port=%d", database.Port)
	}

	return cstring
}
