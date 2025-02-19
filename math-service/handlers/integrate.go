package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"math-service/models"
	"net/http"
)

const g = -9.81
const dt = 0.016 // Интервал времени

// Обработчик для расчета ускорения
func integrateHandler(w http.ResponseWriter, r *http.Request) {
	var rocket models.Rocket

	// Читаем тело запроса
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	fmt.Println("Received raw JSON:", string(body))

	// Декодируем JSON
	err = json.Unmarshal(body, &rocket)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	fmt.Printf("Decoded rocket data: X=%.2f, Y=%.2f, Thrust=%v, Mass=%.2f\n",
		rocket.X, rocket.Y, rocket.Thrust, rocket.Mass)

	if rocket.Mass <= 0 {
		http.Error(w, fmt.Sprintf("Invalid rocket data: mass = %f, thrust = %v", rocket.Mass, rocket.Thrust), http.StatusBadRequest)
		return
	}

	// Рассчитываем ускорение, обновляем скорость и позицию
	acceleration, updatedRocket := calculateRocketMovement(rocket)

	if math.IsNaN(acceleration) {
		http.Error(w, "Invalid acceleration data", http.StatusInternalServerError)
		return
	}

	// Формируем ответ с расчетами
	response := map[string]float64{
		"acceleration": updatedRocket.Acceleration,
		"velocity_y":   updatedRocket.VelocityY,
		"new_y":        updatedRocket.Y,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Функция для расчета ускорения, обновления скорости и позиции ракеты
func calculateRocketMovement(rocket models.Rocket) (float64, models.Rocket) {
	gravityForce := rocket.Mass * g                   // Сила тяжести
	netForce := float64(rocket.Thrust) - gravityForce // Суммарная сила (тяга минус сила тяжести)

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
