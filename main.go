package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type WeatherData struct {
	Temperature float64 `json:"temperature"`
	WindSpeed   float64 `json:"wind_speed"`
	Humidity    float64 `json:"humidity"`
}

func getWeatherData(city string) (WeatherData, error) {

	apiURL := fmt.Sprintf("http://api.weatherprovider.com/weather?city=%s", city)

	response, err := http.Get(apiURL)
	if err != nil {
		return WeatherData{}, err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return WeatherData{}, err
	}

	var weather WeatherData
	if err := json.Unmarshal(body, &weather); err != nil {
		return WeatherData{}, err
	}

	return weather, nil
}

func weatherHandler(w http.ResponseWriter, r *http.Request) {

	city := r.URL.Query().Get("city")

	if city == "" {
		http.Error(w, "Місто не вказане", http.StatusBadRequest)
		return
	}

	weather, err := getWeatherData(city)
	if err != nil {
		http.Error(w, "Не вдалося отримати погоду", http.StatusInternalServerError)
		return
	}

	jsonResponse, err := json.Marshal(weather)
	if err != nil {
		http.Error(w, "Помилка обробки даних про погоду", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.Write(jsonResponse)
}

func main() {

	http.HandleFunc("/weather", weatherHandler)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Помилка при запуску сервера:", err)
		return
	}
}
