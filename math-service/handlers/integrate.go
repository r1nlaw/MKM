package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"math-service/models"
	"net/http"
)

const g = 9.81
const dt = 0.016 // Интервал времени

var rocketState = models.Rocket{}

// Обработчик для расчета ускорения
func integrateHandler(w http.ResponseWriter, r *http.Request) {
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
	if rocketState.Mass == 0 {
		rocketState = inputRocket
	}

	// Используем предыдущее состояние скорости и позиции
	rocketState.Thrust = inputRocket.Thrust

	// Рассчитываем новое состояние
	_, rocketState = calculateRocketMovement(rocketState)

	// Формируем ответ
	response := map[string]float64{
		"acceleration": rocketState.Acceleration,
		"velocity_y":   rocketState.VelocityY,
		"new_y":        rocketState.Y,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Функция для расчета ускорения, обновления скорости и позиции ракеты
func calculateRocketMovement(rocket models.Rocket) (float64, models.Rocket) {
	gravityForce := rocket.Mass * g        // Сила тяжести
	thrustForce := float64(rocket.Thrust)  // Тяга двигателя
	netForce := gravityForce - thrustForce // Суммарная сила (тяга - гравитация)

	if rocket.Mass == 0 {
		return math.NaN(), rocket
	}

	// Рассчитываем ускорение
	acceleration := netForce / rocket.Mass
	rocket.Acceleration = acceleration

	// Обновляем скорость ракеты с учетом ускорения
	rocket.VelocityY += acceleration * dt

	// Обновляем позицию ракеты
	rocket.Y += rocket.VelocityY * dt

	// Если ракета приземлилась, то она остается на поверхности
	if rocket.Y < 0 {
		rocket.Y = 0
		rocket.VelocityY = 0
	}

	// Печатаем данные для отладки
	fmt.Printf("Acceleration: %.2f, New Y: %.2f, New VelocityY: %.2f, Thrust: %v\n", acceleration, rocket.Y, rocket.VelocityY, rocket.Thrust)

	return acceleration, rocket
}
