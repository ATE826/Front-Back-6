package main

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"
)

var cacheData struct {
	Data      interface{}
	Timestamp time.Time
}

var cacheMutex sync.RWMutex

// Обработчик для данных с кэшированием
func dataHandler(w http.ResponseWriter, r *http.Request) {
	cacheMutex.RLock()
	if time.Since(cacheData.Timestamp) < time.Minute {
		// Возвращаем данные из кэша
		json.NewEncoder(w).Encode(cacheData.Data)
		cacheMutex.RUnlock()
		return
	}
	cacheMutex.RUnlock()

	// Генерация новых данных
	newData := map[string]string{"message": "This is fresh data!"}

	// Обновление кэша
	cacheMutex.Lock()
	cacheData.Data = newData
	cacheData.Timestamp = time.Now()
	cacheMutex.Unlock()

	// Отправка новых данных
	json.NewEncoder(w).Encode(newData)
}
