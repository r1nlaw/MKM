package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io"
	"net/http"
	"os"
	"physics-service/models"
)

// Инициализация начальных данных ракеты
var rocketData = &models.Rocket{
	X:         100,
	Y:         8000,
	Width:     20,
	Height:    60,
	VelocityY: 0,
	FuelMass:  270000,
	Thrust:    0, // Начальная тяга
	Mass:      300000,
}

const (
	imageHeight  = 600  // Высота изображения в пикселях
	worldHeight  = 8000 // Реальная высота в метрах
	surfaceY     = 100  // Высота земли в метрах
	rocketWidth  = 20   // Ширина ракеты
	rocketHeight = 60   // Высота ракеты в пикселях (оставляем без изменений)
)

func drawRocket(y float64) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, 100, imageHeight))

	// Белый фон
	white := color.RGBA{255, 255, 255, 255}
	draw.Draw(img, img.Bounds(), &image.Uniform{white}, image.Point{}, draw.Src)

	// Коэффициент масштабирования
	scale := float64(imageHeight) / worldHeight

	// Преобразуем координаты из метров в пиксели
	rocketY := imageHeight - int(y*scale) // Инверсия оси Y
	groundY := imageHeight - int(surfaceY*scale)

	// Рисуем землю
	groundColor := color.RGBA{0, 255, 0, 255}
	groundRect := image.Rect(0, groundY, 100, imageHeight)
	draw.Draw(img, groundRect, &image.Uniform{groundColor}, image.Point{}, draw.Over)

	// Позиционируем ракету по центру экрана
	centerX := float64(img.Bounds().Dx()) / 2 // Получаем ширину изображения и делим на 2 для центра
	rocketColor := color.RGBA{255, 0, 0, 255}
	rocketRect := image.Rect(int(centerX)-rocketWidth/2, rocketY-rocketHeight, int(centerX)+rocketWidth/2, rocketY)
	draw.Draw(img, rocketRect, &image.Uniform{rocketColor}, image.Point{}, draw.Over)

	return img
}

// Обработчик для обновления позиции ракеты и отправки изображения
func RocketHandler(w http.ResponseWriter, r *http.Request) {
	// Вызываем функцию для вычислений
	_, err := getAccelerationFromMathService(rocketData)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error processing rocket data: %v", err), http.StatusInternalServerError)
		return
	}

	// Проверяем, достигла ли ракета земли
	if rocketData.Y <= surfaceY {
		fmt.Println("Rocket has landed. Shutting down server...")
		os.Exit(0)
	}

	// Печатаем данные ракеты для отладки
	fmt.Printf("Updated rocket data: Y = %.2f, VelocityY = %.2f, Thrust = %v\n", rocketData.Y, rocketData.VelocityY, rocketData.Thrust)

	// Генерируем изображение ракеты
	img := drawRocket(float64(rocketData.Y))

	// Буфер для изображения
	imgBuffer := new(bytes.Buffer)
	if err := png.Encode(imgBuffer, img); err != nil {
		http.Error(w, fmt.Sprintf("Error encoding image: %v", err), http.StatusInternalServerError)
		return
	}

	// Отправляем изображение клиенту
	w.Header().Set("Content-Type", "image/png")
	w.WriteHeader(http.StatusOK)
	w.Write(imgBuffer.Bytes())
}

// Отправка данных в математический микросервис для вычислений
func getAccelerationFromMathService(rocket *models.Rocket) (float64, error) {
	if rocket.Mass <= 0 {
		return 0, fmt.Errorf("invalid rocket mass: %.2f", rocket.Mass)
	}

	// Формируем структуру данных для отправки на математический микросервис
	requestData := models.RocketDataRequest{
		Y:         rocket.Y,
		Thrust:    rocket.Thrust, // Тяга передается сюда
		Mass:      rocket.Mass,
		FuelMass:  rocket.FuelMass,
		VelocityY: rocket.VelocityY,
	}

	// Кодируем данные в JSON
	requestBody, err := json.Marshal(requestData)
	if err != nil {
		return 0, fmt.Errorf("error encoding request data: %v", err)
	}

	// Отправляем POST-запрос на математический микросервис
	resp, err := http.Post("http://localhost:8085/math/integrate", "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return 0, fmt.Errorf("error sending request to math service: %v", err)
	}
	defer resp.Body.Close()

	// Читаем ответ от микросервиса
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("error reading response from math service: %v", err)
	}

	// Проверка, если статус код не 200
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("error from math service: %s", body)
	}

	// Декодируем ответ от математического микросервиса
	var response models.RocketDataResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return 0, fmt.Errorf("error decoding response from math service: %v", err)
	}

	// Обновляем данные ракеты с учетом ответа
	rocket.VelocityY = response.VelocityY
	rocket.Y = response.NewY

	// Логируем результат для отладки
	fmt.Printf("Acceleration: %.2f, New Y: %.2f, New VelocityY: %.2f, Thrust: %v\n", response.Acceleration, rocket.Y, rocket.VelocityY, rocket.Thrust)

	return response.Acceleration, nil
}

// Обработчик для обновления тяги ракеты
func updateRocketThrust(w http.ResponseWriter, r *http.Request) {
	var thrustData struct {
		Thrust int `json:"thrust"`
	}

	// Чтение данных с тела запроса
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	// Разбор JSON
	err = json.Unmarshal(body, &thrustData)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Обновление тяги ракеты
	rocketData.Thrust = thrustData.Thrust

	// Логируем изменение тяги для отладки
	fmt.Printf("Updated Thrust: %d\n", rocketData.Thrust)

	// Пересчет данных с новым значением тяги
	acceleration, err := getAccelerationFromMathService(rocketData)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error updating math service: %v", err), http.StatusInternalServerError)
		return
	}

	// Логируем обновленные данные ракеты
	fmt.Printf("Updated rocket data: Y = %.2f, VelocityY = %.2f, Thrust = %d\n", rocketData.Y, rocketData.VelocityY, rocketData.Thrust)

	// Отправляем обновленные данные клиенту
	response := map[string]float64{
		"new_thrust":   float64(rocketData.Thrust),
		"acceleration": acceleration,
		"velocity_y":   rocketData.VelocityY,
		"new_y":        rocketData.Y,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
