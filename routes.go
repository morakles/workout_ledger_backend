package main

import (
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"
)

func registerRoutes() {
	http.HandleFunc("/hello", helloHandler)
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/dbtest", testDBselect)

	// Swagger UI: /swagger/index.html
	http.HandleFunc("/swagger/", httpSwagger.WrapHandler)
}
