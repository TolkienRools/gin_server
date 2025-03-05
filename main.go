package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
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
	debugMode := strings.ToLower(os.Getenv("DEBUG")) == "true"

	if debugMode {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// 3. Создаем логгер в зависимости от режима
	var logger *zap.Logger
	var err error

	if debugMode {
		// Режим разработки с дебаг-логированием
		config := zap.NewDevelopmentConfig()
		config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
		logger, err = config.Build()
	} else {
		// Production-режим с ограничением уровня INFO
		config := zap.NewProductionConfig()
		config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
		logger, err = config.Build()
	}

	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	// 4. Инициализация Gin
	r := gin.New()

	// 5. Подключаем middleware с правильным уровнем логирования
	r.Use(ginzap.Ginzap(logger, "2006-01-02 15:04:05", true))
	r.Use(ginzap.RecoveryWithZap(logger, true))

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
