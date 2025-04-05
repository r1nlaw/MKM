package handlers

import (
	"encoding/json"
	"io"
	"math-service/models"
	"net/http"
)

type IntegrattionService struct {
	rocketState models.Rocket
}

// Обработчик для расчета движения ракеты
func (s *IntegrattionService) Integrate(w http.ResponseWriter, r *http.Request) {
	var inputRocket models.Rocket

	// Читаем тело запроса
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	// Декодируем JSON
	err = json.Unmarshal(body, &inputRocket)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Если это первый запрос — инициализируем состояние
	if s.rocketState.Mass == 0 {
		s.rocketState = inputRocket
	}
	if s.rocketState.InitialMass == 0 {
		s.rocketState.InitialMass = s.rocketState.Mass // Запоминаем стартовую массу
	}

	// Если топлива нет, сбрасываем тягу до 0
	if s.rocketState.FuelMass <= 0 {
		s.rocketState.Thrust = 0
	} else {
		s.rocketState.Thrust = inputRocket.Thrust
	}

	// Рассчитываем новое состояние
	_, s.rocketState = s.сalculateRocketMovement(s.rocketState)

	// Формируем ответ
	response := map[string]float64{
		"acceleration": s.rocketState.Acceleration,
		"velocity_y":   s.rocketState.VelocityY,
		"y":            s.rocketState.Y,
		"fuel_mass":    s.rocketState.FuelMass,
		"mass":         s.rocketState.Mass,
		"drag":         s.rocketState.Drag,
		"energy":       s.rocketState.Energy,
		"initial_mass": s.rocketState.InitialMass,
		"losses":       s.rocketState.Losses,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
