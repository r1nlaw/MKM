package main

import (
	"fmt"
	"math-service/handlers"
	"net/http"

	"github.com/rs/cors"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, world!")
}

func main() {

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8080"}, // Разрешаем доступ только с порта 8080
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
	})

	http.HandleFunc("/", handler)

	// Применяем middleware для CORS
	handlerWithCors := c.Handler(http.DefaultServeMux)

	http.HandleFunc("/math/integrate", handlers.IntegrateHandler)
	http.HandleFunc("/math/vector", handlers.CurrentVectorHandler)
	http.HandleFunc("/math/force", handlers.ForceHandler)
	http.HandleFunc("/math/trajectory", handlers.TrajectoryHandler)

	http.ListenAndServe(":8085", handlerWithCors)

}
