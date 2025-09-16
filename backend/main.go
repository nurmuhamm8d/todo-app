package main

import (
	"context"
	"database/sql"
	"embed"
	"log"
	"os"

	"github.com/lib/pq"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed frontend/dist/*
var assets embed.FS

func mustDB() *sql.DB {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=127.0.0.1 port=5432 user=postgres password=postgres dbname=todoapp sslmode=disable"
	}
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	if err := ensureSchema(db); err != nil {
		log.Fatal(err)
	}
	return db
}

func main() {
	_ = pq.Array
	db := mustDB()
	app := NewApp(db)

	appOptions := &options.App{
		Title:  "Todo App",
		Width:  1200,
		Height: 800,
		OnStartup: func(ctx context.Context) {
			app.startup(ctx)
		},
		AssetServer: &assetserver.Options{
			Assets:  assets,
			Handler: nil,
		},
		Bind: []interface{}{
			app,
		},
		Debug: options.Debug{
			OpenInspectorOnStartup: false,
		},
		WindowStartState: options.Normal,
	}

	err := wails.Run(appOptions)
	if err != nil {
		log.Fatal(err)
	}
}
