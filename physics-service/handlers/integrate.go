package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"physics-service/models"
)

func integrateHandler(w http.ResponseWriter, r *http.Request) {

	var state models.RocketState

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&state)
	if err != nil {
		http.Error(w, "Невозможно обработать запрос", http.StatusBadRequest)
		return
	}

	// Отправляем запрос в математический микросервис для численного интегрирования
	result, err := sendRequestToMathService("http://localhost:8085/math/integrate", state)
	if err != nil {
		http.Error(w, "Ошибка при выполнении интегрирования", http.StatusInternalServerError)
		return
	}

	// Преобразуем результат в нужный тип
	rocketState, ok := result.(models.RocketState)
	if !ok {
		http.Error(w, "Ошибка при получении результата интегрирования", http.StatusInternalServerError)
		return
	}

	// Возвращаем результат
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(rocketState)
	if err != nil {
		log.Println("Ошибка кодировки: ", err)
		http.Error(w, "Ошибка отправки ответа", http.StatusInternalServerError)
	}
}
