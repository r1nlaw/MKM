package models

type Rocket struct {
	Mass   float64 `json:"mass"`
	Fuel   float64 `json:"fuel"`
	Thrust float64 `json:"thrust"`
}

type ForceResult struct {
	ThrustY  float64 `json:"thrust_y"`    // Сила тяги для Y
	GravityY float64 `json:"gravity_y"`   // Сила тяжести для Y
	ResultFY float64 `json:"resultant_y"` // Результирующая сила для Y
}
