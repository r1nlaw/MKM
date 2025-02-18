package models

type RocketState struct {
	X      float64 `json:"x"`      // Позиция x
	Y      float64 `json:"y"`      // Позиция y
	Vx     float64 `json:"vx"`     // Скорость x
	Vy     float64 `json:"vy"`     // Скорость y
	Ax     float64 `json:"ax"`     // Ускорение x
	Ay     float64 `json:"ay"`     // Ускорение y
	Fuel   float64 `json:"fuel"`   // Оставшееся топливо
	Mass   float64 `json:"mass"`   // Масса
	Thrust float64 `json:"thrust"` // Тяга
}
