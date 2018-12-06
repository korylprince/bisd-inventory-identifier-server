package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/handlers"
	"github.com/korylprince/bisd-inventory-identifier-server/httpapi"
)

func main() {
	db, err := sql.Open(config.SQLDriver, config.SQLDSN)
	if err != nil {
		log.Fatalln("Could not open database:", err)
	}
	httpapi.Debug = config.Debug

	s := httpapi.NewServer(db, os.Stdout, config.Secret)

	chain := handlers.CombinedLoggingHandler(os.Stdout,
		handlers.CompressHandler(handlers.CORS(
			handlers.AllowedOrigins([]string{"*"}),
			handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS"}),
			handlers.AllowedHeaders([]string{"Accept", "Content-Type", "Origin", "Authorization"}),
		)(
			http.StripPrefix(config.Prefix, s.Router()),
		)))

	log.Println("Listening on:", config.ListenAddr)
	log.Println(http.ListenAndServe(config.ListenAddr, chain))
}
