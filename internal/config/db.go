package core

import (
	"context"
	"goapp/ent"
	"log"

	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"

	_ "github.com/mattn/go-sqlite3"
)

func NewDBClient(dsn string, autoMigrate bool) *ent.Client {

	drv, err := sql.Open(dialect.SQLite, dsn)
	if err != nil {
		log.Fatalf("failed opening sqlite connection: %v", err)
	}

	drv.DB().SetMaxOpenConns(1)

	client := ent.NewClient(ent.Driver(drv))

	if autoMigrate {
		if err := client.Schema.Create(context.Background()); err != nil {
			log.Fatalf("failed creating schema: %v", err)
		}
	}

	return client.Debug()
}
