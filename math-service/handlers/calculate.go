package handlers

import (
	"fmt"
	"math"
	"math-service/models"
)

const (
	g       = 9.81  // Ускорение свободного падения (м/с²)
	dt      = 0.016 // Интервал времени (шаг симуляции)
	ve      = 6000  // Скорость истечения газа (м/с)
	rho0    = 1.225 // Плотность воздуха у поверхности (кг/м³)
	H_scale = 10000 // Высота масштаба атмосферы (10 км)
	v_set   = 350   // Установленная скорость (м/с)

)

// Функция для расчета движения ракеты с уравнением Мещерского
func (s *IntegrattionService) сalculateRocketMovement(rocket models.Rocket) (float64, models.Rocket) {
	if rocket.Mass == 0 {
		return math.NaN(), rocket
	}

	// Расход топлива
	fuelConsumed := s.calculateFuelConsumed(rocket.Thrust)

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
	rocket.Losses = s.calculateLossesVelocityTsiolkovsky(rocket)

	// Сила сопротивления воздуха
	rocket.Drag = s.calculateDragForce(rocket)

	// Ускорение по уравнению Мещерского: a = (T - mg + D) / m
	rocket.Acceleration = s.calculateAcceleration(rocket)

	// Обновляем скорость и позицию
	rocket.VelocityY += rocket.Acceleration * dt
	rocket.Y += rocket.VelocityY * dt

	// Если ракета приземлилась, останавливаем ее
	if rocket.Y < 0 {
		rocket.Y = 0
		rocket.VelocityY = 0
	}

	// Проверка сохранения энергии
	rocket.Energy = s.calculateTotalEnergy(rocket, fuelConsumed)

	fmt.Printf("Acceleration: %.2f, Y: %.2f, VelocityY: %.2f, Fuel: %.2f, Energy: %.2f, Mass: %.2f, Drag: %.2f\n",
		rocket.Acceleration, rocket.Y, rocket.VelocityY, rocket.FuelMass, rocket.Energy, rocket.Mass, rocket.Drag)

	return rocket.Acceleration, rocket
}

func (s *IntegrattionService) calculateFuelConsumed(thrust int) float64 {
	massFlowRate := float64(thrust) / ve

	return massFlowRate * dt
}

func (s *IntegrattionService) calculateLossesVelocityTsiolkovsky(rocket models.Rocket) float64 {
	if rocket.InitialMass > rocket.Mass {
		deltaV := ve * math.Log(rocket.InitialMass/rocket.Mass)
		actualDeltaV := math.Abs(rocket.VelocityY)

		return deltaV - actualDeltaV
	}
	return 0
}

func (s *IntegrattionService) calculateDragForce(rocket models.Rocket) float64 {
	return rocket.Mass * g * (rocket.VelocityY * rocket.VelocityY) / (v_set * v_set) * rho0
}

func (s *IntegrattionService) calculateAcceleration(rocket models.Rocket) float64 {
	return (float64(rocket.Thrust) - (rocket.Mass * g) + rocket.Drag) / rocket.Mass
}

func (s *IntegrattionService) calculateTotalEnergy(rocket models.Rocket, fuelConsumed float64) float64 {
	kineticEnergy := 0.5 * rocket.Mass * rocket.VelocityY * rocket.VelocityY
	potentialEnergy := rocket.Mass * g * rocket.Y
	fuelEnergy := 0.0
	if fuelConsumed > 0 {
		fuelEnergy = (rocket.FuelMass + fuelConsumed) * ve * ve / 2
	}
	return kineticEnergy + potentialEnergy + fuelEnergy
}
