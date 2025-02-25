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
	"time"

	"golang.org/x/exp/rand"
)

// Инициализация начальных данных ракеты
var rocketData = &models.Rocket{
	Y:         3000,
	Width:     20,
	Height:    60,
	VelocityY: 0,
	FuelMass:  270000,
	Thrust:    75000,
	Mass:      300000,
}

const (
	imageHeight  = 800  // Высота изображения в пикселях
	worldHeight  = 3000 // Реальная высота в метрах
	surfaceY     = 100  // Высота земли в метрах
	platform     = 115
	rocketWidth  = 20
	rocketHeight = 60
)

func drawRocket(y, velocityY float64) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, 1900, imageHeight))

	// Чёрный фон (космос)
	black := color.RGBA{0, 0, 0, 255}
	draw.Draw(img, img.Bounds(), &image.Uniform{black}, image.Point{}, draw.Src)

	// Коэффициент масштабирования
	scale := float64(imageHeight) / worldHeight
	// Добавляем звёзды случайными точками
	white := color.RGBA{255, 255, 255, 255}
	for i := 0; i < 50; i++ {
		x := rand.Intn(img.Bounds().Dx())
		y := rand.Intn(imageHeight - int(platform*scale)) // Только в воздухе
		img.Set(x, y, white)
	}
	// Позиционируем ракету по центру
	centerX := float64(img.Bounds().Dx()) / 2

	// Преобразуем координаты из метров в пиксели
	rocketY := imageHeight - int(y*scale) // Инверсия оси Y
	groundY := imageHeight - int(surfaceY*scale)

	// Рисуем землю
	groundColor := color.RGBA{0, 255, 0, 255}
	groundRect := image.Rect(0, groundY, 1900, imageHeight)
	draw.Draw(img, groundRect, &image.Uniform{groundColor}, image.Point{}, draw.Over)

	// Рисуем корпус ракеты
	rocketBody := color.RGBA{200, 200, 200, 255}
	bodyRect := image.Rect(int(centerX)-rocketWidth/2, rocketY-rocketHeight, int(centerX)+rocketWidth/2, rocketY)
	draw.Draw(img, bodyRect, &image.Uniform{rocketBody}, image.Point{}, draw.Over)

	// Рисуем нос ракеты (треугольник)
	noseColor := color.RGBA{200, 200, 200, 255}
	for i := 0; i < rocketWidth; i++ {
		for j := 0; j < i/2; j++ { // Создаём треугольную форму
			img.Set(int(centerX)-i/2+j, rocketY-rocketHeight-5-j, noseColor)
			img.Set(int(centerX)+i/2-j, rocketY-rocketHeight-5-j, noseColor)
		}
	}

	// Рисуем стабилизаторы
	stabilizerColor := color.RGBA{255, 0, 0, 255}
	for i := -3; i <= 3; i++ {
		img.Set(int(centerX)-rocketWidth/2-3, rocketY-10+i, stabilizerColor)
		img.Set(int(centerX)+rocketWidth/2+3, rocketY-10+i, stabilizerColor)
	}

	// Если ракета достигла земли, рисуем флаг успеха или аварии
	if y <= platform {
		flagX := int(centerX) + rocketWidth/2 + 15
		flagY := imageHeight - int(platform*scale) - 15

		var flagColor color.RGBA
		if velocityY < -5 {
			flagColor = color.RGBA{255, 0, 0, 255} // Красный флаг (авария)
		} else {
			flagColor = color.RGBA{0, 255, 0, 255} // Зелёный флаг (успешно)
		}

		// Палка для флага
		poleColor := color.RGBA{150, 150, 150, 255}
		poleRect := image.Rect(flagX-2, flagY-8, flagX, flagY+15) // Высота флагштока ~30 пикселей
		draw.Draw(img, poleRect, &image.Uniform{poleColor}, image.Point{}, draw.Over)

		// Рисуем флаг (прямоугольник)
		flagRect := image.Rect(flagX, flagY-8, flagX+10, flagY+5)
		draw.Draw(img, flagRect, &image.Uniform{flagColor}, image.Point{}, draw.Over)
	}

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

	// Генерируем изображение ракеты
	img := drawRocket(float64(rocketData.Y), float64(rocketData.VelocityY))

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

	// Проверяем, достигла ли ракета земли
	if rocketData.Y <= platform {
		if rocketData.VelocityY >= -5 {
			fmt.Println("Rocket has landed successfully!")
		} else {
			fmt.Println("Rocket crashed!")
		}

		// Даем клиенту время отобразить последний кадр перед закрытием
		go func() {
			time.Sleep(16 * time.Millisecond)
			os.Exit(0)
		}()
	}

}

// Отправка данных в математический микросервис для вычислений
func getAccelerationFromMathService(rocket *models.Rocket) (float64, error) {
	if rocket.Mass <= 0 {
		return 0, fmt.Errorf("invalid rocket mass: %.2f", rocket.Mass)
	}

	// Формируем структуру данных для отправки на математический микросервис
	requestData := models.RocketDataRequest{
		Y:         rocket.Y,
		Thrust:    rocket.Thrust,
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
	rocket.Acceleration = response.Acceleration
	rocket.VelocityY = response.VelocityY
	rocket.Y = response.NewY
	rocket.FuelMass = response.FuelMass
	rocket.Mass = response.Mass
	rocket.Drag = response.Drag
	rocket.TotalEnergy = response.TotalEnergy

	// Логируем результат
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

func updateDataHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Кодируем текущее состояние ракеты
	if err := json.NewEncoder(w).Encode(rocketData); err != nil {
		http.Error(w, fmt.Sprintf("Failed to encode rocket data: %v", err), http.StatusInternalServerError)
		return
	}
}
