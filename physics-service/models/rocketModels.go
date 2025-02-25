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
	TotalEnergy  float64 `json:"energy"`
	Drag         float64 `json:"drag"`
}

type RocketDataRequest struct {
	Y         float64 `json:"y"`
	Thrust    int     `json:"thrust"`
	Mass      float64 `json:"mass"`
	FuelMass  float64 `json:"fuel_mass"`
	VelocityY float64 `json:"velocity_y"`
}

type RocketDataResponse struct {
	Acceleration float64 `json:"acceleration"`
	VelocityY    float64 `json:"velocity_y"`
	NewY         float64 `json:"y"`
	TotalEnergy  float64 `json:"energy"`
	Drag         float64 `json:"drag"`
	Mass         float64 `json:"mass"`
	FuelMass     float64 `json:"fuel_mass"`
}
