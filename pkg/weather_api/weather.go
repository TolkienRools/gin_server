package weatherapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type WeatherAPI struct {
	APIkey string
}

func (wh *WeatherAPI) GetWeatherData(lat, lon string, timeout int) interface{} {

	url := fmt.Sprintf("http://api.weatherapi.com/v1/current.json?key=%s&q=%s,%s", wh.APIkey, lat, lon)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)

	defer cancel()

	resp, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	// resp, err := http.Get(url)

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
