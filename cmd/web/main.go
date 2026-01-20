package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/anditakaesar/uwa-go-fullstack/internal/env"
	"github.com/anditakaesar/uwa-go-fullstack/internal/infra"
	"github.com/anditakaesar/uwa-go-fullstack/internal/server"
	"github.com/anditakaesar/uwa-go-fullstack/internal/xlog"
)

func main() {
	db, err := infra.NewDatabase()
	if err != nil {
		xlog.Logger.Error(fmt.Sprintf("unable to connect to database: %v", err))
		os.Exit(1)
	}
	defer db.Close()

	svr := server.SetupServer(&server.ServerDependency{
		DB: db,
	})
	http.ListenAndServe(env.Values.Port, svr)
}
