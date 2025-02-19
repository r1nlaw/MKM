package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"math-service/models"
	"net/http"
)

const (
	g  = 9.81  // Ускорение свободного падения
	dt = 0.016 // Интервал времени
	ve = 3000  // Скорость истечения газа (м/с)
)

var rocketState = models.Rocket{}

// Обработчик для расчета движения ракеты
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

	// Если топлива нет, сбрасываем тягу до 0 и не позволяем увеличивать
	if rocketState.FuelMass <= 0 {
		rocketState.Thrust = 0
	} else {
		rocketState.Thrust = inputRocket.Thrust
	}

	// Рассчитываем новое состояние
	_, rocketState = calculateRocketMovement(rocketState)

	// Формируем ответ
	response := map[string]float64{
		"acceleration": rocketState.Acceleration,
		"velocity_y":   rocketState.VelocityY,
		"new_y":        rocketState.Y,
		"fuel_mass":    rocketState.FuelMass,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Функция для расчета движения ракеты с учетом уравнения Мещерского
func calculateRocketMovement(rocket models.Rocket) (float64, models.Rocket) {
	if rocket.Mass == 0 {
		return math.NaN(), rocket
	}

	// Расчет изменения массы топлива
	massFlowRate := float64(rocket.Thrust) / ve // Расход топлива
	fuelConsumed := massFlowRate * dt

	if fuelConsumed > rocket.FuelMass {
		fuelConsumed = rocket.FuelMass // Не сжигаем больше, чем есть
	}
	rocket.FuelMass -= fuelConsumed
	rocket.Mass -= fuelConsumed // Общая масса уменьшается

	// Если топлива нет, сбрасываем тягу до 0
	if rocket.FuelMass <= 0 {
		rocket.Thrust = 0
	}

	// Рассчитываем ускорение по уравнению Мещерского
	acceleration := g - (float64(rocket.Thrust) / rocket.Mass)
	rocket.Acceleration = acceleration

	// Обновляем скорость и позицию
	rocket.VelocityY += acceleration * dt
	rocket.Y += rocket.VelocityY * dt

	// Если ракета приземлилась, останавливаем ее
	if rocket.Y < 0 {
		rocket.Y = 0
		rocket.VelocityY = 0
	}

	// Проверка сохранения энергии
	kineticEnergy := 0.5 * rocket.Mass * rocket.VelocityY * rocket.VelocityY
	potentialEnergy := rocket.Mass * g * rocket.Y
	// Добавляем только энергию топлива, если оно сжигается
	fuelEnergy := 0.0
	if fuelConsumed > 0 {
		fuelEnergy = (rocket.FuelMass + fuelConsumed) * ve * ve / 2
	}
	totalEnergy := kineticEnergy + potentialEnergy + fuelEnergy

	fmt.Printf("Acceleration: %.2f, New Y: %.2f, New VelocityY: %.2f, Fuel: %.2f, Energy: %.2f, Mass: %.2f\n",
		acceleration, rocket.Y, rocket.VelocityY, rocket.FuelMass, totalEnergy, rocket.Mass)

	return acceleration, rocket
}
