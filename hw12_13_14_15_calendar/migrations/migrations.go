package migrations

import (
	"database/sql"
	"embed"

	_ "github.com/jackc/pgx/v5/stdlib" //nolint
	_ "github.com/lib/pq"              //nolint
	"github.com/pressly/goose/v3"      //nolint
)

//go:embed *.sql
var embedMigrations embed.FS

func Migrate(db *sql.DB) {
	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		panic(err)
	}

	if err := goose.Up(db, "."); err != nil {
		panic(err)
	}
}
