package handlers

import (
	"encoding/json"
	"log"
	"math"
	"net/http"
)

func currentVectorHandler(w http.ResponseWriter, r *http.Request) {
	var vector1, vector2 Vector
	var result Vector
	var err error

	// Декодируем запрос для получения двух векторов
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&struct {
		Vector1 Vector `json:"vector1"`
		Vector2 Vector `json:"vector2"`
	}{
		Vector1: vector1,
		Vector2: vector2,
	})
	if err != nil {
		http.Error(w, "Невозможно обработать запрос", http.StatusBadRequest)
		return
	}
	vector1.Add(&vector2)

	result = vector1
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		log.Println("Ошибка кодировки:", err)
		http.Error(w, "Ошибка отправки ответа", http.StatusInternalServerError)
	}
}

type Vector struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// Сложение двух векторов
func (v1 *Vector) Add(v2 *Vector) {
	v1.X += v2.X
	v1.Y += v2.Y
}

// Вычитание двух векторов
func (v1 *Vector) Subtract(v2 *Vector) {
	v1.X -= v2.X
	v1.Y -= v2.Y
}

// Умножение вектора на скаляр
func (v *Vector) MultiplyByScalar(scalar float64) {
	v.X *= scalar
	v.Y *= scalar
}

// Скалярное произведение двух векторов
func (v1 *Vector) Dot(v2 *Vector) float64 {
	return v1.X*v2.Y + v1.Y*v2.Y
}

// Длина (модуль) вектора
func (v *Vector) Length() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

// Нормализация вектора
func (v *Vector) Normalize() {
	length := v.Length()
	if length != 0 {
		v.X /= length
		v.Y /= length
	}
}
