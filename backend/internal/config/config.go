package config

import "os"

type Config struct {
	DBConn    string
	JWTSecret string
}

func Load() Config {
	c := Config{
		DBConn:    os.Getenv("TODOAPP_DB"),
		JWTSecret: os.Getenv("TODOAPP_JWT"),
	}
	if c.DBConn == "" {
		c.DBConn = "user=postgres password=postgres dbname=todoapp host=127.0.0.1 port=5432 sslmode=disable"
	}
	if c.JWTSecret == "" {
		c.JWTSecret = "devsecret"
	}
	return c
}
