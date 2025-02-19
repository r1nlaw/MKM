package handlers

import (
	"fmt"
	"log"
	"net/http"
	"physics-service/models"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Handle WebSocket connection
func webSocketConnectionHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade WebSocket:", err)
		return
	}
	defer conn.Close()

	for {
		// Получаем данные ракеты от клиента
		var incomingRocket models.Rocket
		err := conn.ReadJSON(&incomingRocket)
		if err != nil {
			log.Println("Error reading rocket data:", err)
			break
		}

		// Обновляем глобальные данные ракеты
		rocketData.X = incomingRocket.X
		rocketData.Y = incomingRocket.Y
		rocketData.Thrust = incomingRocket.Thrust
		rocketData.Mass = incomingRocket.Mass

		fmt.Printf("Before calling math service: X=%.2f, Y=%.2f, Thrust=%v, Mass=%.2f\n",
			rocketData.X, rocketData.Y, rocketData.Thrust, rocketData.Mass)

		// Проверяем, что данные корректны
		if rocketData.Mass <= 0 {
			log.Println("Error: Invalid rocket data (Mass must be > 0), skipping request")
			continue
		}

		// Отправляем данные в математический микросервис для вычислений
		acceleration, err := getAccelerationFromMathService(rocketData)
		if err != nil {
			log.Println("Error fetching acceleration from math service:", err)
			continue
		}

		// Обновляем физические параметры ракеты
		rocketData.VelocityY += acceleration
		rocketData.Y += rocketData.VelocityY * dt

		// Проверяем, достигла ли ракета поверхности
		if rocketData.Y >= surfaceY {
			rocketData.Y = surfaceY
			rocketData.VelocityY = 0
			log.Println("Rocket has landed.")
		}

		// Отправляем обновленные данные обратно на фронтенд
		err = conn.WriteJSON(rocketData)
		if err != nil {
			log.Println("Error sending data to client:", err)
			break
		}
	}
}
