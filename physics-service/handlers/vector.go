package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"physics-service/models"
)

func vectorHandler(w http.ResponseWriter, r *http.Request) {

	var state models.RocketState

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&state)
	if err != nil {
		http.Error(w, "Невозможно обработать запрос", http.StatusBadRequest)
		return
	}

	// Отправляем запрос в математический микросервис для вычсления вектора
	result, err := sendRequestToMathService("http://localhost:8085/math/vector", state)
	if err != nil {
		http.Error(w, "Ошибка при вычислении вектора", http.StatusInternalServerError)
		return
	}

	// Преобразуем в нужный тип
	vector, ok := result.(models.Vector)
	if !ok {
		http.Error(w, "Ошибка при получении данных о векторе", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(vector)
	if err != nil {
		log.Println("Ошибка кодировки: ", err)
		http.Error(w, "Ошибка отправки ответа", http.StatusInternalServerError)
	}

}
