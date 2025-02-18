package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

func sendRequestToMathService(url string, data interface{}) (interface{}, error) {

	body, err := json.Marshal(data)
	if err != nil {
		log.Printf("Ошибка при маршаллинге данных: ", err)
		return nil, err
	}

	// Создаем новый HTTP-запрос
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		log.Printf("Ошибка при создании запроса: ", err)
		return nil, err
	}
	// Устанавливаем заголовки
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Создаем HTTP-клиент с таймаутом
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Отправляем запрос и получаем ответ
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Ошибка при отправке запроса: ", err)
		return nil, err
	}
	defer resp.Body.Close()

	// Проверяем статус
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Ошибка при получении ответа: статус: ", resp.StatusCode)
	}

	// Декодируем ответ
	var result interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		log.Printf("Ошибка при декодировании ответа: ", err)
		return nil, err
	}

	return result, nil
}
