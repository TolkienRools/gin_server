package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type WeatherServer struct{}

const API_KEY = "6ef704c680454e1eb7691220242208"

// Добавить передачу параметров
func (wh *WeatherServer) getWeatherData(lat string, lon string) interface{} {

	url := fmt.Sprintf("http://api.weatherapi.com/v1/current.json?key=%s&q=%s,%s", API_KEY, lat, lon)
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

	fmt.Println(string(body))

	return json_result
}

func (wh *WeatherServer) getWeatherHandler(c *gin.Context) {
	latitude := c.Query("lat")
	longitude := c.Query("lon")

	result := wh.getWeatherData(latitude, longitude)

	c.JSON(http.StatusOK, result)
}

func (wh *WeatherServer) postUploadFileHandler(c *gin.Context) {
	form, _ := c.MultipartForm()
	files := form.File["upload[]"]

	for _, file := range files {
		save_path := fmt.Sprintf("./uploads/%s", file.Filename)
		c.SaveUploadedFile(file, save_path)
	}
	c.JSON(http.StatusOK, gin.H{"message": "files uploaded"})
}

func main() {
	gin.DisableConsoleColor()

	f, _ := os.Create("log.txt")
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
	r := gin.Default()

	r.Use(cors.Default())
	r.LoadHTMLGlob("templates/*")
	server := WeatherServer{}

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	r.GET("/upload", func(c *gin.Context) {
		c.HTML(http.StatusOK, "file_upload.html", nil)
	})

	{
		api := r.Group("/api")
		api.GET("location/", server.getWeatherHandler)
		api.POST("upload/", server.postUploadFileHandler)
	}
	r.Run()
}
