package api

type contextKey int

//TransactionKey is the context key for the database transaction for a request
const TransactionKey contextKey = 0
