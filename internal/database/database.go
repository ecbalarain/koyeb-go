package database

import (
	"crypto/tls"
	"database/sql"
	"fmt"
	"time"

	"github.com/go-sql-driver/mysql"
)

// Connect opens a connection pool to the MySQL/TiDB database.
func Connect(databaseURL string) (*sql.DB, error) {
	// Register TLS config for TiDB Cloud
	_ = mysql.RegisterTLSConfig("tidb", &tls.Config{
		MinVersion: tls.VersionTLS12,
	})

	// Determine DSN format and convert if needed
	var dsn string
	if len(databaseURL) > 8 && databaseURL[:8] == "mysql://" {
		// Parse URL format to Go MySQL driver format
		parsedDSN, err := parseDSN(databaseURL)
		if err != nil {
			return nil, fmt.Errorf("invalid DATABASE_URL: %w", err)
		}
		dsn = parsedDSN
	} else {
		// Already in Go MySQL driver format
		dsn = databaseURL
	}

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Connection pool settings
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Verify connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

// parseDSN converts a URL-style DSN to the Go MySQL driver format.
// Input:  mysql://user:pass@host:port/dbname?tls=true
// Output: user:pass@tcp(host:port)/dbname?tls=tidb&parseTime=true
func parseDSN(rawURL string) (string, error) {
	// Strip the mysql:// prefix
	if len(rawURL) > 8 && rawURL[:8] == "mysql://" {
		rawURL = rawURL[8:]
	}

	// Split user:pass@host:port/dbname?params
	atIdx := -1
	for i := len(rawURL) - 1; i >= 0; i-- {
		if rawURL[i] == '@' {
			atIdx = i
			break
		}
	}
	if atIdx == -1 {
		return "", fmt.Errorf("missing @ in DSN")
	}

	userPass := rawURL[:atIdx]
	rest := rawURL[atIdx+1:]

	// Split host:port/dbname?params
	slashIdx := -1
	for i := 0; i < len(rest); i++ {
		if rest[i] == '/' {
			slashIdx = i
			break
		}
	}
	if slashIdx == -1 {
		return "", fmt.Errorf("missing database name in DSN")
	}

	hostPort := rest[:slashIdx]
	dbAndParams := rest[slashIdx+1:]

	// Replace tls=true with tls=tidb (our registered config)
	dbAndParams = replaceTLS(dbAndParams)

	// Ensure parseTime=true is included
	if len(dbAndParams) > 0 && dbAndParams[len(dbAndParams)-1] != '?' {
		if !contains(dbAndParams, "parseTime") {
			if contains(dbAndParams, "?") {
				dbAndParams += "&parseTime=true"
			} else {
				dbAndParams += "?parseTime=true"
			}
		}
	}

	return fmt.Sprintf("%s@tcp(%s)/%s", userPass, hostPort, dbAndParams), nil
}

func replaceTLS(s string) string {
	result := []byte(s)
	target := []byte("tls=true")
	replacement := []byte("tls=tidb")
	for i := 0; i <= len(result)-len(target); i++ {
		if string(result[i:i+len(target)]) == string(target) {
			result = append(result[:i], append(replacement, result[i+len(target):]...)...)
			break
		}
	}
	return string(result)
}

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
