package handlers

import (
	"encoding/json"
	"log"
	"math-service/models"
	"net/http"
)

func integrateHandler(w http.ResponseWriter, r *http.Request) {
	var state models.RocketState
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&state)
	if err != nil {
		http.Error(w, "Недопустимы формат запроса", http.StatusBadRequest)
		return
	}
	// Параметры интеграции (шаг времени)
	const dt = 0.01 // 10 мс

	updateRocket(&state, dt)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(state)
	if err != nil {
		log.Println("Ошибка кодировки", err)
		http.Error(w, "Не удалось отправить ответ", http.StatusInternalServerError)
	}

}

func updateRocket(state *models.RocketState, dt float64) {
	g := 9.81

	// Вычисляем силу тяги
	var Fx, Fy float64

	if state.Fuel > 0 {
		Fx = 0
		Fy = state.Thrust               // Считаем, что тяга направлена только вверх
		state.Fuel -= state.Thrust * dt // Трата топлива
	}

	// Вычисляем ускорение методом Эйлера
	state.Ax = Fx / state.Mass
	state.Ay = (Fy / state.Mass) - g

	// Обновляем скорости
	state.Vx += state.Ax * dt
	state.Vy += state.Ay * dt

	// Обновляем позицию
	state.X += state.Vx * dt
	state.Y += state.Vy * dt

}
