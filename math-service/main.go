package main

import (
	"encoding/json"
	"math-service/handlers"
	"net/http"

	"github.com/rs/cors"
)

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{"message": "Поехали!"}
	json.NewEncoder(w).Encode(response)
}

func main() {

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8081"}, // Разрешаем доступ только с порта 8080
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
	})

	http.HandleFunc("/", handler)

	// Применяем middleware для CORS
	handlerWithCors := c.Handler(http.DefaultServeMux)

	integrationService := &handlers.IntegrattionService{}

	http.HandleFunc("/math/integrate", integrationService.Integrate)

	http.ListenAndServe(":8085", handlerWithCors)

}
