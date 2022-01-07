package main

import (
	"database/sql"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/Tursunkhuja/crud/cmd/app"
	"github.com/Tursunkhuja/crud/pkg/customers"
)

func main() {
	host := "0.0.0.0"
	port := "9999"
	dsn := "postgres://app:pass@localhost:5432/db"

	if err := execute(host, port, dsn); err != nil {
		log.Print(err)
		os.Exit(1)
	}
}

// Func start server
func execute(host string, port string, dns string) (err error) {
	db, err := sql.Open("pgx", dns)
	if err != nil {
		return nil
	}
	defer func() {
		if cerr := db.Close(); cerr != nil {
			if err == nil {
				err = cerr
				return
			}
			log.Println(err)
		}
	}()
	mux := http.NewServeMux()
	customerSvc := customers.NewService(db)
	server := app.NewServer(mux, customerSvc)
	server.Init()

	svr := &http.Server{
		Addr:    net.JoinHostPort(host, port),
		Handler: server,
	}
	return svr.ListenAndServe()
}
