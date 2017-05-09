package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/handlers"
	"github.com/korylprince/bisd-inventory-identifier-server/api"
	"github.com/korylprince/bisd-inventory-identifier-server/httpapi"
)

func main() {
	chromeSvc, err := api.NewChromebookService(config.OAuthJSONPath, config.OAuthImpersonateUser)
	if err != nil {
		log.Fatalln("Could not create ChromebookService:", err)
	}

	db, err := sql.Open(config.SQLDriver, config.SQLDSN)
	if err != nil {
		log.Fatalln("Could not open database:", err)
	}

	r := httpapi.NewRouter(os.Stdout, chromeSvc, db)

	chain := handlers.CompressHandler(handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Accept", "Content-Type", "Origin", "X-Session-Key"}),
	)(http.StripPrefix(config.Prefix, r)))

	log.Println("Listening on:", config.ListenAddr)
	log.Println(http.ListenAndServe(config.ListenAddr, chain))
}
