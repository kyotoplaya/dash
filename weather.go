package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	tea "charm.land/bubbletea/v2"
)

type weather string

func tickWeather() tea.Cmd {
	url := "https://wttr.in/Oryol?0?Q"

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("Ошибка создания запроса: %v\n", err)
		os.Exit(1)
	}

	// Обязательно выставляем User-Agent как у curl, чтобы получить чистый текст
	req.Header.Set("User-Agent", "curl/7.79.1")

	// Выполняем запрос
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Ошибка выполнения запроса: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	// Читаем ответ
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Ошибка чтения ответа: %v\n", err)
		os.Exit(1)
	}

	return tea.Tick(5*time.Second, func(t time.Time) tea.Msg {
		return weather(body)
	})
}
