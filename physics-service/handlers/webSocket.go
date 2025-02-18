package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"physics-service/models"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Разрешаем любые соединения
	},
}

var wsClients = make(map[*websocket.Conn]bool) // Список клиентов WebSocket

// Обработчик WebSocket соединений
func wsHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	response := map[string]string{
		"message": "WebSocket data",
	}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Unable to marshal response", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Ошибка при установке соединения WebSocket: ", err)
		return
	}
	defer conn.Close()

	wsClients[conn] = true

	for {
		var rocketState models.RocketState

		err := conn.ReadJSON(&rocketState)
		if err != nil {
			log.Println("Ошибка при чтении данных от клиента: ", err)
			delete(wsClients, conn)
			return
		}

		// Обработка запросов к математическому микросервису
		force, err := sendRequestToMathService("http://localhost:8085/math/force", rocketState)
		if err != nil {
			log.Println("Ошибка при запросе на математический микросервис:", err)
			return
		}

		vector, err := sendRequestToMathService("http://localhost:8085/math/vector", rocketState)
		if err != nil {
			log.Println("Ошибка при запросе на математический микросервис:", err)
			return
		}

		trajectory, err := sendRequestToMathService("http://localhost:8085/math/trajectory", rocketState)
		if err != nil {
			log.Println("Ошибка при запросе на математический микросервис:", err)
			return
		}

		integratedState, err := sendRequestToMathService("http://localhost:8085/math/integrate", rocketState)
		if err != nil {
			log.Println("Ошибка при запросе на математический микросервис:", err)
			return
		}

		// Отправка обновленных данных обратно клиенту через WebSocket
		response := map[string]interface{}{
			"force":           force,
			"vector":          vector,
			"trajectory":      trajectory,
			"integratedState": integratedState,
		}

		err = conn.WriteJSON(response)
		if err != nil {
			log.Println("Ошибка при отправке данных на клиент:", err)
			delete(wsClients, conn)
			return
		}
	}
}
