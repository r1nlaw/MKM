package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"physics-service/models"
)

func forceHandler(w http.ResponseWriter, r *http.Request) {
	var rocket models.Rocket

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&rocket)
	if err != nil {
		http.Error(w, "Невозможно обработать запрос", http.StatusBadRequest)
		return
	}

	// Отправляем запрос в математический микросервис для вычисления силы
	result, err := sendRequestToMathService("http://localhost:8085/math/force", rocket)
	if err != nil {
		http.Error(w, "Ошибка при вычислении силы", http.StatusInternalServerError)
		return
	}

	// Преобразуем результат в нужный тип
	forceResult, ok := result.(models.ForceResult)
	if !ok {
		http.Error(w, "Ошибка при получении данных о силе", http.StatusInternalServerError)
		return
	}

	// Возвращаем результат
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(forceResult)
	if err != nil {
		log.Println("Ошибка кодировки:", err)
		http.Error(w, "Ошибка отправки ответа", http.StatusInternalServerError)
	}
}
