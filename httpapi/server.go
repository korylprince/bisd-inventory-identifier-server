package httpapi

import (
	"database/sql"
	"io"
)

//Server represents shared resources
type Server struct {
	db     *sql.DB
	output io.Writer
	secret string
}

//NewServer returns a new server with the given resources
func NewServer(db *sql.DB, output io.Writer, secret string) *Server {
	return &Server{db: db, output: output, secret: secret}
}
