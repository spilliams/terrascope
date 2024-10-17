package server

import (
	"embed"
	"net/http"

	"github.com/gin-gonic/gin"
)

//go:embed index.html public/assets/*
var static embed.FS

func Run() error {
	router := setupRouter()
	// Listen and Server in 0.0.0.0:8080
	return router.Run(":8080")
}

var db = make(map[string]string)

func setupRouter() *gin.Engine {
	// gin.SetMode(gin.ReleaseMode)

	// Disable Console Color
	// gin.DisableConsoleColor()
	router := gin.Default()
	router.SetTrustedProxies([]string{"127.0.0.1"})

	router.StaticFS("/", http.FS(static))

	return router
}
