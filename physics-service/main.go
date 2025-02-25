package main

import (
	"encoding/json"
	"net/http"
	"physics-service/handlers"

	"github.com/rs/cors"
)

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{"message": "Поехали!"}
	json.NewEncoder(w).Encode(response)
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

	http.HandleFunc("/physics/rocket-image", handlers.DrowRocketHandler)
	http.HandleFunc("/physics/update-thrust", handlers.UpdateRocketThrust)
	http.HandleFunc("/physics/update-data", handlers.UpdateDataHandler)

	http.ListenAndServe(":8086", handlerWithCors)
}
