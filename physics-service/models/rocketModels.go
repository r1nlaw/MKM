package models

type Rocket struct {
	X            float64 `json:"x"`
	Y            float64 `json:"y"`
	Width        float64 `json:"width"`
	Height       float64 `json:"height"`
	VelocityY    float64 `json:"velocity_y"`
	Thrust       int     `json:"thrust"`
	Mass         float64 `json:"mass"`
	Acceleration float64 `json:"acceleration"`
}

// Структура для запроса данных о ракете
type RocketDataRequest struct {
	X         float64 `json:"x"`
	Y         float64 `json:"y"`
	Thrust    int     `json:"thrust"`
	Mass      float64 `json:"mass"`
	VelocityY float64 `json:"velocity_y"`
}

type RocketDataResponse struct {
	Acceleration float64 `json:"acceleration"`
	VelocityY    float64 `json:"velocity_y"`
	NewY         float64 `json:"new_y"`
}
