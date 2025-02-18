package handlers

import (
	"encoding/json"
	"log"
	"math-service/models"
	"net/http"
)

func trajectoryHandler(w http.ResponseWriter, r *http.Request) {
	var rocketState models.RocketState

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&rocketState)
	if err != nil {
		http.Error(w, "Невозможно обработать запрос", http.StatusBadRequest)
		return
	}

	// Параметры расчета траектории
	const dt = 0.01
	const maxTime = 500.0

	trajectory := calculateTrajectory(&rocketState, dt, maxTime)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(trajectory)
	if err != nil {
		log.Println("Ошибка кодировки:", err)
		http.Error(w, "Ошибка отправки ответа", http.StatusInternalServerError)
	}

}

func calculateTrajectory(rocket *models.RocketState, dt, maxTime float64) models.Trajectory {
	var trajectory models.Trajectory

	var time float64

	// Пошаговое численное интегрированние траектории ракеты
	for time < maxTime {
		// Обновляем положение и скорость ракеты
		updateRocket(rocket, dt)

		// Добавляем текущее положение в траекторию
		trajectory.X = append(trajectory.X, rocket.X)
		trajectory.Y = append(trajectory.Y, rocket.Y)

		// Обновляем время
		time += dt

	}
	return trajectory

}
