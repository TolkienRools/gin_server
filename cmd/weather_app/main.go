package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/TolkienRools/gin_server/internal/config"
	weather "github.com/TolkienRools/gin_server/internal/handlers"
	"github.com/gin-contrib/cors"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type WeatherServer struct{}

// const API_KEY = "6ef704c680454e1eb7691220242208"

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
	cfg := config.MustLoad()

	fmt.Println(cfg)

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
	r.LoadHTMLGlob("web/templates/*")
	server := WeatherServer{}

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	r.GET("/upload", func(c *gin.Context) {
		c.HTML(http.StatusOK, "file_upload.html", nil)
	})

	{
		api := r.Group("/api")
		api.GET("location/", weather.GetWeatherHandler)
		api.POST("upload/", server.postUploadFileHandler)
	}
	r.Run()
}
