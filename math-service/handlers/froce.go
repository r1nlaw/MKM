package handlers

import (
	"encoding/json"
	"log"
	"math-service/models"
	"net/http"
)

func forceHandler(w http.ResponseWriter, r *http.Request) {
	var rocket models.Rocket

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&rocket)
	if err != nil {
		http.Error(w, "Невозможно обработать запрос", http.StatusBadRequest)
		return
	}

	// Рассчет силы
	force := calculateForce(&rocket)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(force)
	if err != nil {
		log.Println("Ошибка кодировки:", err)
		http.Error(w, "Ошибка отправки ответа", http.StatusInternalServerError)
	}

}

const g = 9.81

func calculateForce(rocket *models.Rocket) models.ForceResult {
	var force models.ForceResult

	// Сила тяги
	if rocket.Fuel > 0 {
		force.ThrustY = rocket.Thrust       // Сила тяги вверх
		rocket.Fuel -= rocket.Thrust * 0.01 // Трата топлива
	} else {
		force.ThrustY = 0 // Нет топлива, нет тяги
	}

	// Сила тяжести
	force.GravityY = -rocket.Mass * g

	// Результирующая сила
	force.ResultFY = force.ThrustY + force.GravityY

	return force
}
