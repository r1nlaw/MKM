package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"physics-service/models"
)

func trajectoryHandler(w http.ResponseWriter, r *http.Request) {

	var state models.RocketState

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&state)
	if err != nil {
		http.Error(w, "Невозможно обработать запрос", http.StatusBadRequest)
		return
	}

	// Отправляем запрос в математический микросервис для вычисления траектории
	result, err := sendRequestToMathService("http://localhost:8085/math/trajectory", state)
	if err != nil {
		http.Error(w, "Ошибка при вычислении траектории", http.StatusInternalServerError)
		return
	}

	// Преобразуем в нужный тип
	trajectory, ok := result.(models.Trajectory)
	if !ok {
		http.Error(w, "Ошибка при получении данных о траектории", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(trajectory)
	if err != nil {
		log.Println("Ошибка кодировки: ", err)
		http.Error(w, "Ошибка отправки ответа", http.StatusInternalServerError)
	}
}
