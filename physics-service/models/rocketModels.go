package models

type Rocket struct {
	X            float64 `json:"x"`
	Y            float64 `json:"y"`
	Width        float64 `json:"width"`
	Height       float64 `json:"height"`
	VelocityY    float64 `json:"velocity_y"`
	Thrust       float64 `json:"thrust"`
	Mass         float64 `json:"mass"`
	Acceleration float64 `json:"acceleration"`
}
