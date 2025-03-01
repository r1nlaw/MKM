package models

type Rocket struct {
	Y            float64 `json:"y"`
	Width        float64 `json:"width"`
	Height       float64 `json:"height"`
	VelocityY    float64 `json:"velocity_y"`
	Thrust       int     `json:"thrust"`
	Mass         float64 `json:"mass"`
	FuelMass     float64 `json:"fuel_mass"`
	Acceleration float64 `json:"acceleration"`
	Drag         float64 `json:"drag"`
	Energy       float64 `json:"energy"`
	InitialMass  float64 `json:"initial_mass"`
	Losses       float64 `json:"losses"`
}
