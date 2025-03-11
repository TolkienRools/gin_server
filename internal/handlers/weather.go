package weather

import (
	"net/http"

	"github.com/gin-gonic/gin"

	weatherapi "github.com/TolkienRools/gin_server/pkg/weather_api"
)

func GetWeatherHandler(c *gin.Context) {
	latitude := c.Query("lat")
	longitude := c.Query("lon")

	weather_api := weatherapi.WeatherAPI{
		APIkey: "6ef704c680454e1eb7691220242208",
	}

	result := weather_api.GetWeatherData(latitude, longitude, 3)

	c.JSON(http.StatusOK, result)
}
