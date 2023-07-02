package migrations

import (
	"context"
	"database/sql"
	"embed"

	"github.com/SanaripEsep/esep-backend/pkg/logging"
	_ "github.com/lib/pq"

	"github.com/pressly/goose/v3"
)

//go:embed *.sql
var migs embed.FS

func Up(ctx context.Context, url string, log logging.GooseLogger) error {

	goose.SetBaseFS(migs)

	conn, err := sql.Open("postgres", url)
	if err != nil {
		return err
	}
	defer conn.Close()

	goose.SetLogger(log)

	return goose.Up(conn, ".")
}
