package main

import (
	"fmt"
	"net/http"
	"physics-service/handlers"

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

	http.HandleFunc("/physics/force", handlers.ForceHandler)
	http.HandleFunc("/physics/vector", handlers.VectorHandler)
	http.HandleFunc("/physics/trajectory", handlers.TrajectoryHandler)
	http.HandleFunc("/physics/integrate", handlers.IntegrateHandler)
	http.HandleFunc("/physics/ws", handlers.WebSocketHandler)

	http.ListenAndServe(":8086", handlerWithCors)
}
