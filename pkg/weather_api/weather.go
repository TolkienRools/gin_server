package weatherapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type WeatherAPI struct {
	APIkey string
}

func (wh *WeatherAPI) GetWeatherData(lat string, lon string) interface{} {

	url := fmt.Sprintf("http://api.weatherapi.com/v1/current.json?key=%s&q=%s,%s", wh.APIkey, lat, lon)
	resp, err := http.Get(url)

	if err != nil {
		// Лучше использовать логирование из Gin
		panic(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		panic(err)
	}

	var json_result map[string]interface{}

	if err := json.Unmarshal(body, &json_result); err != nil {
		panic(err)
	}

	return json_result
}
