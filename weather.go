package main

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

type weatherModule struct{}

func (weatherModule) Name() string { return "Weather" }
func (weatherModule) Collect() (string, error) {
	url := "https://wttr.in/Oryol?0?Q"

	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("Ошибка создания запроса: %v\n", err)
		return "Ошибка создания запроса", err
	}

	// Обязательно выставляем User-Agent как у curl, чтобы получить чистый текст
	req.Header.Set("User-Agent", "curl/7.79.1")

	// Выполняем запрос
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Ошибка выполнения запроса: %v\n", err)
		return "Ошибка выполнения запроса", err
	}
	defer resp.Body.Close()

	// Читаем ответ
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Ошибка чтения ответа: %v\n", err)
		return "Ошибка чтения ответа", err
	}

	return string(body), nil
}

func (weatherModule) Interval() time.Duration {
	return time.Hour // Обновляем раз в час
}
