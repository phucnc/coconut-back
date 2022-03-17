package database

import (
	"fmt"
	"os"
)

type PostgresConfig struct {
	Host    string
	Port    string
	Usr     string
	Pwd     string
	Db      string
	SSLMode string
}

func (c *PostgresConfig) ToConnStr() string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s", c.Usr, c.Pwd, c.Host, c.Port, c.Db, c.SSLMode)
}

func NewPostgreConfigFromEnv() *PostgresConfig {
	config := &PostgresConfig{
		Host:    os.Getenv("PG_HOST"),
		Port:    os.Getenv("PG_PORT"),
		Usr:     os.Getenv("PG_USR"),
		Pwd:     os.Getenv("PG_PWD"),
		Db:      os.Getenv("PG_DB"),
		SSLMode: os.Getenv("PG_SSL_MODE"),
	}
	return config
}
