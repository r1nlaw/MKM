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
	g       = 9.81  // Ускорение свободного падения (м/с²)
	dt      = 0.016 // Интервал времени (шаг симуляции)
	ve      = 6000  // Скорость истечения газа (м/с)
	rho0    = 1.225 // Плотность воздуха у поверхности (кг/м³)
	Cd      = 0.5   // Коэффициент лобового сопротивления
	A       = 12.0  // Площадь поперечного сечения ракеты (м²)
	H_scale = 8000  // Высота масштаба атмосферы (~8 км)
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

	// Если топлива нет, сбрасываем тягу до 0
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
		"y":            rocketState.Y,
		"fuel_mass":    rocketState.FuelMass,
		"mass":         rocketState.Mass,
		"drag":         rocketState.Drag,
		"energy":       rocketState.Energy,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Функция для расчета движения ракеты с уравнением Мещерского
func calculateRocketMovement(rocket models.Rocket) (float64, models.Rocket) {
	if rocket.Mass == 0 {
		return math.NaN(), rocket
	}

	// Расход топлива
	massFlowRate := float64(rocket.Thrust) / ve
	fuelConsumed := massFlowRate * dt

	// Ограничение, чтобы не сжигать больше топлива, чем есть
	if fuelConsumed > rocket.FuelMass {
		fuelConsumed = rocket.FuelMass
	}
	rocket.FuelMass -= fuelConsumed
	rocket.Mass -= fuelConsumed

	// Если топлива нет, сбрасываем тягу до 0
	if rocket.FuelMass <= 0 {
		rocket.Thrust = 0
	}

	// Плотность воздуха на текущей высоте (экспоненциальное убывание)
	rho := rho0 * math.Exp(-rocket.Y/H_scale)

	// Сила сопротивления воздуха
	dragForce := 0.5 * Cd * rho * A * rocket.VelocityY * rocket.VelocityY
	if rocket.VelocityY > 0 {
		dragForce = -dragForce // Сопротивление всегда противоположно движению
	}
	rocket.Drag = dragForce

	// Ускорение по уравнению Мещерского: a = (T - mg + D) / m
	acceleration := (float64(rocket.Thrust) - (rocket.Mass * g) + dragForce) / rocket.Mass
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
	fuelEnergy := 0.0
	if fuelConsumed > 0 {
		fuelEnergy = (rocket.FuelMass + fuelConsumed) * ve * ve / 2
	}
	totalEnergy := kineticEnergy + potentialEnergy + fuelEnergy
	rocket.Energy = totalEnergy

	fmt.Printf("Acceleration: %.2f, Y: %.2f, VelocityY: %.2f, Fuel: %.2f, Energy: %.2f, Mass: %.2f, Drag: %.2f\n",
		acceleration, rocket.Y, rocket.VelocityY, rocket.FuelMass, totalEnergy, rocket.Mass, dragForce)

	return acceleration, rocket
}
