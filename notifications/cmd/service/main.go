package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"zenport/internal/config"
	"zenport/internal/system"
	"zenport/internal/web"
	"zenport/notifications"
	"zenport/ntps/migrations"

	_ "github.com/jackc/pgx/v4/stdlib"
)

func main() {
	if err := run(); err != nil {
		fmt.Printf("stores exitted abnormally: %s\n", err)
		os.Exit(1)
	}
}

func run() (err error) {
	var cfg config.AppConfig
	cfg, err = config.InitConfig()
	if err != nil {
		return err
	}
	s, err := system.NewSystem(cfg)
	if err != nil {
		return err
	}
	defer func(db *sql.DB) {
		if err = db.Close(); err != nil {
			return
		}
	}(s.DB())

	if err = s.MigrateDB(migrations.FS); err != nil {
		return err
	}
	s.Mux().Mount("/", http.FileServer(http.FS(web.WebUI)))
	// call the module composition root
	if err = notifications.Root(s.Waiter().Context(), s); err != nil {
		return err
	}
	fmt.Println("started stores service")
	defer fmt.Println("stopped stores service")

	s.Waiter().Add(
		s.WaitForWeb,
		s.WaitForRPC,
		s.WaitForStream,
	)

	return s.Waiter().Wait()
}
