package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

var cacheData []byte
var lastCacheTime time.Time

func dataHandler(w http.ResponseWriter, r *http.Request) {
	// Проверка на устаревший кэш
	if time.Since(lastCacheTime) > time.Minute {
		// Генерация новых данных (например, считывание из базы)
		data := map[string]string{
			"message": "Данные, обновлённые на сервере",
		}
		dataJSON, _ := json.Marshal(data)

		// Сохранение в кэш
		cacheData = dataJSON
		lastCacheTime = time.Now()
	} else {
		log.Println("Отправка данных из кэша")
	}

	// Отправка данных
	w.Header().Set("Content-Type", "application/json")
	w.Write(cacheData)
}
